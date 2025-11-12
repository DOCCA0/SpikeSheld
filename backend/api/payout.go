package api

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"spikeshield/contracts"
	"spikeshield/db"
	"spikeshield/utils"
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
		utils.LogWarning("Could not verify oracle address: %v", err)
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

	utils.LogInfo("🚀 Calling executePayout on-chain for user %s, policy %d", policy.UserAddress, policy.ID)
	utils.LogInfo("   Gas Price: %s wei", gasPrice.String())
	utils.LogInfo("   Gas Limit: %d", auth.GasLimit)

	// *** REAL ON-CHAIN TRANSACTION ***
	tx, err := ps.Contract.ExecutePayout(
		auth,
		common.HexToAddress(policy.UserAddress),
		big.NewInt(int64(policy.ID)),
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
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	utils.LogInfo("✅ Transaction mined in block %d", receipt.BlockNumber.Uint64())
	utils.LogInfo("   Gas Used: %d", receipt.GasUsed)

	// Record payout in database
	payout := &db.Payout{
		PolicyID:    policy.ID,
		UserAddress: policy.UserAddress,
		Amount:      policy.CoverageAmount,
		SpikeID:     spike.ID,
		TxHash:      tx.Hash().Hex(),
	}

	if err := db.InsertPayout(payout); err != nil {
		return fmt.Errorf("failed to insert payout record: %w", err)
	}

	// Update policy status
	if err := db.UpdatePolicyStatus(policy.ID, "claimed"); err != nil {
		return fmt.Errorf("failed to update policy status: %w", err)
	}

	utils.LogInfo("💰 Payout executed successfully for user %s: $%.2f (tx: %s)",
		policy.UserAddress, policy.CoverageAmount, tx.Hash().Hex())

	return nil
}

// Close closes the client connection
func (ps *PayoutService) Close() {
	if ps.Client != nil {
		ps.Client.Close()
	}
}

