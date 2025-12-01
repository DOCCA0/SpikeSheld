package api

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"spikeshield/contracts"
	"spikeshield/db"
	"spikeshield/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PayoutService handles on-chain payout executions
type PayoutService struct {
	Client          *ethclient.Client
	ContractAddress common.Address
	Contract        *contracts.InsurancePool
	PrivateKey      *ecdsa.PrivateKey
	ChainID         *big.Int
}

// NewPayoutService creates a new payout service instance
func NewPayoutService(rpcURL, contractAddr, privateKeyHex string) (*PayoutService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	contractAddress := common.HexToAddress(contractAddr)

	// Create contract instance
	insuranceContract, err := contracts.NewInsurancePool(contractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract instance: %w", err)
	}

	// Verify oracle address matches
	oracleAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	currentOracle, err := insuranceContract.Oracle(&bind.CallOpts{})
	if err != nil {
		utils.LogError("Could not verify oracle address: %v", err)
	} else if currentOracle != oracleAddr {
		return nil, fmt.Errorf("private key does not match contract oracle. Expected: %s, Got: %s", currentOracle.Hex(), oracleAddr.Hex())
	}

	utils.LogInfo("✅ Connected to InsurancePool contract at %s", contractAddress.Hex())
	utils.LogInfo("✅ Oracle address: %s", oracleAddr.Hex())

	return &PayoutService{
		Client:          client,
		ContractAddress: contractAddress,
		Contract:        insuranceContract,
		PrivateKey:      privateKey,
		ChainID:         chainID,
	}, nil
}

// ExecutePayout triggers on-chain payout for a spike event
func (ps *PayoutService) ExecutePayout(spike *db.Spike) error {
	utils.LogInfo("Executing payout for spike ID %d", spike.ID)

	// Get all active policies
	policies, err := db.GetActivePolicies()
	if err != nil {
		return fmt.Errorf("failed to get active policies: %w", err)
	}

	if len(policies) == 0 {
		utils.LogInfo("No active policies found, skipping payout")
		return nil
	}

	utils.LogInfo("Found %d active policy/policies", len(policies))

	// Execute payout for each active policy
	for _, policy := range policies {
		if err := ps.executeForPolicy(policy, spike); err != nil {
			utils.LogError("Failed to execute payout for policy %d: %v", policy.ID, err)
			continue
		}
	}

	return nil
}

// executeForPolicy executes payout for a single policy
func (ps *PayoutService) executeForPolicy(policy *db.Policy, spike *db.Spike) error {
	// Create transaction auth
	auth, err := bind.NewKeyedTransactorWithChainID(ps.PrivateKey, ps.ChainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	// Get current nonce
	publicKeyECDSA, ok := ps.PrivateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("failed to get public key")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := ps.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// Set gas parameters
	auth.GasLimit = uint64(300000) // Increase if needed

	// Get suggested gas price
	gasPrice, err := ps.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}
	auth.GasPrice = gasPrice

	// Query on-chain policies to find correct active policy ID
	userAddr := common.HexToAddress(policy.UserAddress)
	onchainPolicies, err := ps.Contract.GetUserPolicies(&bind.CallOpts{}, userAddr)
	if err != nil {
		return fmt.Errorf("failed to get user policies: %w", err)
	}

	now := time.Now().Unix()
	var targetPolicyId int64 = -1
	for i := range onchainPolicies {
		p := onchainPolicies[i]
		if p.Active && !p.Claimed && p.ExpiryTime.Int64() >= now {
			targetPolicyId = int64(i)
			break // Use first active unclaimed policy
		}
	}

	if targetPolicyId == -1 {
		utils.LogInfo("No active unclaimed policy found on-chain for user %s (DB policy %d)", policy.UserAddress, policy.ID)
		return nil
	}

	// Check pool balance
	onchainPolicy, err := ps.Contract.GetPolicy(&bind.CallOpts{}, userAddr, big.NewInt(targetPolicyId))
	if err != nil {
		return fmt.Errorf("failed to get on-chain policy %d: %w", targetPolicyId, err)
	}
	poolBal, err := ps.Contract.GetPoolBalance(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get pool balance: %w", err)
	}

	// Convert coverage to wei (6 decimals)
	coverageWei := onchainPolicy.CoverageAmount
	covUSDT := new(big.Float).Quo(new(big.Float).SetInt(coverageWei), big.NewFloat(1e6))
	covFloat, _ := covUSDT.Float64()
	if poolBal.Cmp(coverageWei) < 0 {
		utils.LogError("Insufficient pool balance for user %s policy %d: %s < %s wei (%.2f USDT)", userAddr.Hex(), targetPolicyId, poolBal.String(), coverageWei.String(), covFloat)
		return nil
	}

	utils.LogInfo("🚀 Calling executePayout on-chain for user %s, on-chain policy ID %d", userAddr.Hex(), targetPolicyId)
	utils.LogInfo("   Gas Price: %s wei", gasPrice.String())
	utils.LogInfo("   Gas Limit: %d", auth.GasLimit)
	utils.LogInfo("   Pool balance: %s wei OK", poolBal.String())

	// *** REAL ON-CHAIN TRANSACTION ***
	tx, err := ps.Contract.ExecutePayout(
		auth,
		userAddr,
		big.NewInt(targetPolicyId),
		big.NewInt(int64(spike.ID)),
	)
	if err != nil {
		return fmt.Errorf("failed to execute payout transaction: %w", err)
	}

	utils.LogInfo("📤 Transaction sent: %s", tx.Hash().Hex())
	utils.LogInfo("⏳ Waiting for transaction to be mined...")

	// Wait for transaction to be mined
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, ps.Client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != 1 {
		utils.LogInfo("Payout tx reverted (status %d) for policy %d - likely race condition or already claimed", receipt.Status, policy.ID)
		return nil
	}

	utils.LogInfo("✅ Transaction mined in block %d", receipt.BlockNumber.Uint64())
	utils.LogInfo("   Gas Used: %d", receipt.GasUsed)

	// // Record payout in database (use DB policy ID)
	// payout := &db.Payout{
	// 	PolicyID:    policy.ID,
	// 	UserAddress: policy.UserAddress,
	// 	Amount:      policy.CoverageAmount,
	// 	SpikeID:     spike.ID,
	// 	TxHash:      tx.Hash().Hex(),
	// }

	// if err := db.InsertPayout(payout); err != nil {
	// 	return fmt.Errorf("failed to insert payout record: %w", err)
	// }

	// // Update policy status
	// if err := db.UpdatePolicyStatus(policy.ID, "claimed"); err != nil {
	// 	return fmt.Errorf("failed to update policy status: %w", err)
	// }

	utils.LogInfo("💰 Payout executed successfully for user %s policy %d: $%.2f (tx: %s)",
		userAddr.Hex(), targetPolicyId, covFloat, tx.Hash().Hex())

	return nil
}

// Close closes the client connection
func (ps *PayoutService) Close() {
	if ps.Client != nil {
		ps.Client.Close()
	}
}
