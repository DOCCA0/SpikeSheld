package eventlistener

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"spikeshield/db"
	"spikeshield/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EventListener monitors and syncs contract events to the database
type EventListener struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
	pollInterval    time.Duration
}

// PolicyPurchasedEvent represents the PolicyPurchased event from the contract
type PolicyPurchasedEvent struct {
	PolicyId   *big.Int
	Premium    *big.Int
	Coverage   *big.Int
	ExpiryTime *big.Int
}

// PayoutExecutedEvent represents the PayoutExecuted event from the contract
type PayoutExecutedEvent struct {
	PolicyId *big.Int
	Amount   *big.Int
}

// NewEventListener creates a new event listener instance
func NewEventListener(rpcURL, contractAddr string, pollInterval time.Duration) (*EventListener, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	// Parse contract ABI
	contractABI, err := abi.JSON(strings.NewReader(getInsurancePoolABI()))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	return &EventListener{
		client:          client,
		contractAddress: common.HexToAddress(contractAddr),
		contractABI:     contractABI,
		pollInterval:    pollInterval,
	}, nil
}

// Start begins listening for contract events
func (el *EventListener) Start(ctx context.Context) error {
	utils.LogInfo("🎧 Event listener started for contract: %s", el.contractAddress.Hex())

	ticker := time.NewTicker(el.pollInterval)
	defer ticker.Stop()

	// Initial sync
	if err := el.syncEvents(ctx); err != nil {
		utils.LogError("Initial event sync failed: %v", err)
	}

	// Poll for new events
	for {
		select {
		case <-ctx.Done():
			utils.LogInfo("Event listener stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := el.syncEvents(ctx); err != nil {
				utils.LogError("Event sync failed: %v", err)
			}
		}
	}
}

// Close closes the Ethereum client connection
func (el *EventListener) Close() {
	if el.client != nil {
		el.client.Close()
	}
}

// syncEvents fetches and processes new events since last sync
func (el *EventListener) syncEvents(ctx context.Context) error {
	// Get last synced block
	lastBlock, err := db.GetLastSyncedBlock(el.contractAddress)
	if err != nil {
		return fmt.Errorf("failed to get last synced block: %w", err)
	}

	// Get current block number
	currentBlock, err := el.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block: %w", err)
	}

	if lastBlock >= currentBlock {
		return nil // No new blocks
	}

	// If this is the first sync, start from recent blocks (last 1000 blocks or from lastBlock)
	if lastBlock == 0 {
		if currentBlock > 1000 {
			lastBlock = currentBlock - 1000
		}
	}

	utils.LogInfo("Syncing events from block %d to %d", lastBlock+1, currentBlock)

	// Fetch events in chunks to avoid RPC limits
	chunkSize := uint64(1000)
	for fromBlock := lastBlock + 1; fromBlock <= currentBlock; fromBlock += chunkSize {
		toBlock := fromBlock + chunkSize - 1
		if toBlock > currentBlock {
			toBlock = currentBlock
		}

		if err := el.fetchAndProcessEvents(ctx, fromBlock, toBlock); err != nil {
			return fmt.Errorf("failed to fetch events [%d-%d]: %w", fromBlock, toBlock, err)
		}
	}

	// Update last synced block
	if err := db.UpdateLastSyncedBlock(el.contractAddress, currentBlock); err != nil {
		return fmt.Errorf("failed to update sync state: %w", err)
	}

	return nil
}

// fetchAndProcessEvents fetches events from a block range and processes them
func (el *EventListener) fetchAndProcessEvents(ctx context.Context, fromBlock, toBlock uint64) error {
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: []common.Address{el.contractAddress},
	}

	logs, err := el.client.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	for _, vLog := range logs {
		if err := el.processLog(vLog); err != nil {
			utils.LogError("Failed to process log (tx: %s, index: %d): %v", vLog.TxHash.Hex(), vLog.Index, err)
			continue
		}
	}

	if len(logs) > 0 {
		utils.LogInfo("Processed %d events from blocks %d-%d", len(logs), fromBlock, toBlock)
	}

	return nil
}

// processLog processes a single log entry
func (el *EventListener) processLog(vLog types.Log) error {
	eventSignature := vLog.Topics[0].Hex()

	switch eventSignature {
	case el.contractABI.Events["PolicyPurchased"].ID.Hex():
		return el.handlePolicyPurchased(vLog)
	case el.contractABI.Events["PayoutExecuted"].ID.Hex():
		return el.handlePayoutExecuted(vLog)
	default:
		// Unknown event, skip
		return nil
	}
}

// handlePolicyPurchased processes PolicyPurchased events
func (el *EventListener) handlePolicyPurchased(vLog types.Log) error {
	var event PolicyPurchasedEvent

	// Parse indexed user address from topics[1]
	user := common.HexToAddress(vLog.Topics[1].Hex())

	// Parse non-indexed data
	err := el.contractABI.UnpackIntoInterface(&event, "PolicyPurchased", vLog.Data)
	if err != nil {
		return fmt.Errorf("failed to unpack PolicyPurchased event: %w", err)
	}

	utils.LogInfo("📝 PolicyPurchased: user=%s, policyId=%s, premium=%s, coverage=%s",
		user.Hex(), event.PolicyId.String(), event.Premium.String(), event.Coverage.String())

	// Store event in database
	return db.InsertPolicyFromEvent(user, event.Premium, event.Coverage, event.ExpiryTime, vLog)
}

// handlePayoutExecuted processes PayoutExecuted events
func (el *EventListener) handlePayoutExecuted(vLog types.Log) error {
	var event PayoutExecutedEvent

	// Parse indexed user address from topics[1]
	user := common.HexToAddress(vLog.Topics[1].Hex())

	// Parse non-indexed data
	err := el.contractABI.UnpackIntoInterface(&event, "PayoutExecuted", vLog.Data)
	if err != nil {
		return fmt.Errorf("failed to unpack PayoutExecuted event: %w", err)
	}

	utils.LogInfo("💰 PayoutExecuted: user=%s, policyId=%s, amount=%s",
		user.Hex(), event.PolicyId.String(), event.Amount.String())

	// Store event in database
	return db.InsertPayoutFromEvent(user, event.PolicyId, event.Amount, vLog)
}

// getInsurancePoolABI returns the ABI string for the InsurancePool contract
func getInsurancePoolABI() string {
	return `[
		{
			"anonymous": false,
			"inputs": [
				{"indexed": true, "internalType": "address", "name": "user", "type": "address"},
				{"indexed": false, "internalType": "uint256", "name": "policyId", "type": "uint256"},
				{"indexed": false, "internalType": "uint256", "name": "premium", "type": "uint256"},
				{"indexed": false, "internalType": "uint256", "name": "coverage", "type": "uint256"},
				{"indexed": false, "internalType": "uint256", "name": "expiryTime", "type": "uint256"}
			],
			"name": "PolicyPurchased",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{"indexed": true, "internalType": "address", "name": "user", "type": "address"},
				{"indexed": false, "internalType": "uint256", "name": "policyId", "type": "uint256"},
				{"indexed": false, "internalType": "uint256", "name": "amount", "type": "uint256"}
			],
			"name": "PayoutExecuted",
			"type": "event"
		}
	]`
}
