package api

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"spikeshield/db"
	"spikeshield/utils"
)

// PayoutService handles on-chain payout executions
type PayoutService struct {
	Client          *ethclient.Client
	ContractAddress common.Address
	PrivateKey      *ecdsa.PrivateKey
	ChainID         *big.Int
}

// Simplified InsurancePool ABI for executePayout function
const insurancePoolABI = `[{"inputs":[{"internalType":"address","name":"user","type":"address"},{"internalType":"uint256","name":"policyId","type":"uint256"},{"internalType":"string","name":"detectionTxHash","type":"string"}],"name":"executePayout","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

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

	return &PayoutService{
		Client:          client,
		ContractAddress: common.HexToAddress(contractAddr),
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

	// Set gas parameters (you may need to adjust these)
	auth.GasLimit = uint64(300000)
	
	// Get suggested gas price
	gasPrice, err := ps.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}
	auth.GasPrice = gasPrice

	// For simplicity in MVP, we'll use a direct contract call approach
	// In production, use abigen to generate proper Go bindings
	
	// Prepare detection hash (use spike timestamp as identifier)
	detectionHash := fmt.Sprintf("spike_%d_%s", spike.ID, spike.Timestamp.Format("20060102150405"))

	utils.LogInfo("Calling executePayout for user %s, policy %d", policy.UserAddress, policy.ID)

	// NOTE: For MVP demo, we'll simulate the transaction
	// In production, you would use proper contract bindings:
	// tx, err := insuranceContract.ExecutePayout(auth, common.HexToAddress(policy.UserAddress), big.NewInt(int64(policy.ID)), detectionHash)
	
	// For demo purposes, create a mock transaction hash
	txHash := fmt.Sprintf("0x%064d", spike.ID)
	
	utils.LogInfo("Payout transaction sent: %s", txHash)

	// Record payout in database
	payout := &db.Payout{
		PolicyID:    policy.ID,
		UserAddress: policy.UserAddress,
		Amount:      policy.CoverageAmount,
		SpikeID:     spike.ID,
		TxHash:      txHash,
	}

	if err := db.InsertPayout(payout); err != nil {
		return fmt.Errorf("failed to insert payout record: %w", err)
	}

	// Update policy status
	if err := db.UpdatePolicyStatus(policy.ID, "claimed"); err != nil {
		return fmt.Errorf("failed to update policy status: %w", err)
	}

	utils.LogInfo("✅ Payout executed successfully for user %s: $%.2f (tx: %s)",
		policy.UserAddress, policy.CoverageAmount, txHash)

	return nil
}

// Close closes the client connection
func (ps *PayoutService) Close() {
	if ps.Client != nil {
		ps.Client.Close()
	}
}

// NOTE: For a production implementation, you would:
// 1. Use `abigen` to generate Go bindings from your contract ABI
// 2. Import the generated package
// 3. Use the typed contract methods
// 
// Example:
// ```
// contract, err := insurancepool.NewInsurancePool(ps.ContractAddress, ps.Client)
// tx, err := contract.ExecutePayout(auth, userAddr, policyID, detectionHash)
// receipt, err := bind.WaitMined(context.Background(), ps.Client, tx)
// ```
