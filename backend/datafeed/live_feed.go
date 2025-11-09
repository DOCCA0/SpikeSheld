package datafeed

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"spikeshield/db"
	"spikeshield/utils"
)

// LiveFeed fetches real-time price data from Chainlink Oracle
type LiveFeed struct {
	Client       *ethclient.Client
	FeedAddress  common.Address
	Symbol       string
	PollInterval time.Duration
}

// Simplified AggregatorV3Interface ABI for latestRoundData
const aggregatorABI = `[{"inputs":[],"name":"latestRoundData","outputs":[{"internalType":"uint80","name":"roundId","type":"uint80"},{"internalType":"int256","name":"answer","type":"int256"},{"internalType":"uint256","name":"startedAt","type":"uint256"},{"internalType":"uint256","name":"updatedAt","type":"uint256"},{"internalType":"uint80","name":"answeredInRound","type":"uint80"}],"stateMutability":"view","type":"function"}]`

// NewLiveFeed creates a new live feed instance
func NewLiveFeed(rpcURL, feedAddress, symbol string, pollInterval int) (*LiveFeed, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	return &LiveFeed{
		Client:       client,
		FeedAddress:  common.HexToAddress(feedAddress),
		Symbol:       symbol,
		PollInterval: time.Duration(pollInterval) * time.Second,
	}, nil
}

// Start begins polling price data from Chainlink
func (lf *LiveFeed) Start(ctx context.Context) error {
	utils.LogInfo("Starting live feed for %s with %v interval", lf.Symbol, lf.PollInterval)

	ticker := time.NewTicker(lf.PollInterval)
	defer ticker.Stop()

	// Fetch immediately on start
	if err := lf.fetchAndStore(); err != nil {
		utils.LogError("Failed to fetch initial price: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			utils.LogInfo("Live feed stopped")
			return nil
		case <-ticker.C:
			if err := lf.fetchAndStore(); err != nil {
				utils.LogError("Failed to fetch price: %v", err)
			}
		}
	}
}

// fetchAndStore fetches latest price from Chainlink and stores in DB
func (lf *LiveFeed) fetchAndStore() error {
	// Create bound contract instance
	boundContract := bind.NewBoundContract(lf.FeedAddress, parseABI(), lf.Client, lf.Client, lf.Client)

	var result []interface{}
	err := boundContract.Call(&bind.CallOpts{}, &result, "latestRoundData")
	if err != nil {
		return fmt.Errorf("failed to call latestRoundData: %w", err)
	}

	// Parse result: roundId, answer, startedAt, updatedAt, answeredInRound
	if len(result) < 5 {
		return fmt.Errorf("unexpected result length")
	}

	answer := result[1].(*big.Int)
	updatedAt := result[3].(*big.Int)

	// Chainlink returns price with 8 decimals for BTC/USD
	price := new(big.Float).SetInt(answer)
	price.Quo(price, big.NewFloat(1e8))
	priceFloat, _ := price.Float64()

	timestamp := time.Unix(updatedAt.Int64(), 0)

	priceData := &db.PriceData{
		Timestamp: timestamp,
		Symbol:    lf.Symbol,
		Close:     priceFloat,
		Open:      priceFloat,
		High:      priceFloat,
		Low:       priceFloat,
		Volume:    0,
	}

	if err := db.InsertPrice(priceData); err != nil {
		return fmt.Errorf("failed to insert price: %w", err)
	}

	utils.LogInfo("Fetched price for %s: $%.2f at %s", lf.Symbol, priceFloat, timestamp.Format(time.RFC3339))
	return nil
}

// parseABI is a helper to parse the aggregator ABI
func parseABI() *bind.MetaData {
	// For simplicity in MVP, we'll use a simplified approach
	// In production, you would use abigen to generate Go bindings
	return nil
}

// Close closes the client connection
func (lf *LiveFeed) Close() {
	if lf.Client != nil {
		lf.Client.Close()
	}
}
