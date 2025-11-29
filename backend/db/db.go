package db

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"spikeshield/contracts"
	"spikeshield/utils"

	"github.com/bytedance/gopkg/util/logger"
	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var DB *sql.DB

// InsertNotifier is a channel for notifying when new prices are inserted
var InsertNotifier chan struct{}

// InitNotifier initializes the insert notification channel
func InitNotifier() {
	InsertNotifier = make(chan struct{}, 100)
}

// PriceData represents a price record
type PriceData struct {
	ID        int
	Timestamp time.Time
	Symbol    string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// Spike represents a detected spike event
type Spike struct {
	ID                int
	Timestamp         time.Time
	Symbol            string
	Open              float64
	High              float64
	Low               float64
	Close             float64
	BodyRatio         float64
	RangeClosePercent float64
	DetectedAt        time.Time
}

// Policy represents an insurance policy
type Policy struct {
	ID             int
	UserAddress    string
	Premium        float64
	CoverageAmount float64
	PurchaseTime   time.Time
	ExpiryTime     time.Time
	Status         string
	TxHash         string
}

// Payout represents a payout record
type Payout struct {
	ID          int
	PolicyID    int
	UserAddress string
	Amount      float64
	SpikeID     int
	TxHash      string
	ExecutedAt  time.Time
}

// Balance represents an ERC20 token balance cached in DB
type Balance struct {
	ID           int
	TokenAddress string
	UserAddress  string
	Balance      float64
	LastUpdated  time.Time
}

// Connect establishes database connection
func Connect(cfg *utils.Config) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	utils.LogInfo("Database connected successfully")
	return nil
}

// InsertPrice inserts a price record
func InsertPrice(p *PriceData) error {
	query := `INSERT INTO prices (timestamp, symbol, open, high, low, close, volume) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := DB.QueryRow(query, p.Timestamp, p.Symbol, p.Open, p.High, p.Low, p.Close, p.Volume).Scan(&p.ID)

	// Notify listeners that a new price was inserted
	if err == nil && InsertNotifier != nil {
		select {
		case InsertNotifier <- struct{}{}:
		default:
			// Channel full, skip notification
		}
	}

	return err
}

// InsertSpike inserts a spike detection record
func InsertSpike(s *Spike, priceID int) error {
	query := `INSERT INTO spikes (timestamp, symbol, price_id, body_ratio, range_close_percent) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return DB.QueryRow(query, s.Timestamp, s.Symbol, priceID, s.BodyRatio, s.RangeClosePercent).Scan(&s.ID)
}

// GetLatestPrice retrieves the most recent price for a symbol
func GetLatestPrice(symbol string) (*PriceData, error) {
	query := `SELECT id, timestamp, symbol, open, high, low, close, volume 
			  FROM prices WHERE symbol = $1 ORDER BY timestamp DESC LIMIT 1`

	p := &PriceData{}
	err := DB.QueryRow(query, symbol).Scan(&p.ID, &p.Timestamp, &p.Symbol, &p.Open, &p.High, &p.Low, &p.Close, &p.Volume)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetAllPrices retrieves prices within a time range
func GetAllPrices(symbol string) ([]*PriceData, error) {
	// wrapper to unified GetPrices implementation (no limit => all)
	return GetPrices(symbol, 0)
}

// GetActivePolicies retrieves all active policies
func GetActivePolicies() ([]*Policy, error) {
	query := `SELECT id, user_address, premium, coverage_amount, purchase_time, expiry_time, status, COALESCE(tx_hash, '') 
			  FROM policies WHERE status = 'active' AND expiry_time > NOW()`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []*Policy
	for rows.Next() {
		p := &Policy{}
		if err := rows.Scan(&p.ID, &p.UserAddress, &p.Premium, &p.CoverageAmount, &p.PurchaseTime, &p.ExpiryTime, &p.Status, &p.TxHash); err != nil {
			return nil, err
		}
		policies = append(policies, p)
	}
	return policies, nil
}

// GetPoliciesForUser retrieves policies for a specific user address
func GetPoliciesForUser(userAddr string) ([]*Policy, error) {
	query := `SELECT id, user_address, premium, coverage_amount, purchase_time, expiry_time, status, COALESCE(tx_hash, '')
			  FROM policies WHERE user_address = $1 ORDER BY id DESC`

	rows, err := DB.Query(query, userAddr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []*Policy
	for rows.Next() {
		p := &Policy{}
		if err := rows.Scan(&p.ID, &p.UserAddress, &p.Premium, &p.CoverageAmount, &p.PurchaseTime, &p.ExpiryTime, &p.Status, &p.TxHash); err != nil {
			return nil, err
		}
		policies = append(policies, p)
	}
	return policies, nil
}

// GetBalanceForUser returns the cached balance for a token and user
func GetBalanceForUser(tokenAddr string, userAddr string) (*Balance, error) {
	query := `SELECT id, token_address, user_address, balance, last_updated FROM balances WHERE token_address = $1 AND user_address = $2 LIMIT 1`

	b := &Balance{}
	err := DB.QueryRow(query, tokenAddr, userAddr).Scan(&b.ID, &b.TokenAddress, &b.UserAddress, &b.Balance, &b.LastUpdated)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// UpsertBalance inserts or updates the balance for a token/user
func UpsertBalance(tokenAddr string, userAddr string, balance float64) error {
	query := `INSERT INTO balances (token_address, user_address, balance, last_updated)
			  VALUES ($1, $2, $3, NOW())
			  ON CONFLICT (token_address, user_address)
			  DO UPDATE SET balance = EXCLUDED.balance, last_updated = NOW()`
	_, err := DB.Exec(query, tokenAddr, userAddr, balance)
	return err
}

// UpdateBalanceDelta adds delta to existing balance (creates row if not exists)
func UpdateBalanceDelta(tokenAddr string, userAddr string, delta float64) error {
	// Use SQL upsert to increment existing balance or insert new
	query := `INSERT INTO balances (token_address, user_address, balance, last_updated)
			  VALUES ($1, $2, $3, NOW())
			  ON CONFLICT (token_address, user_address)
			  DO UPDATE SET balance = balances.balance + $3, last_updated = NOW()`
	_, err := DB.Exec(query, tokenAddr, userAddr, delta)
	return err
}

// InsertPayout records a payout execution
func InsertPayout(p *Payout) error {
	query := `INSERT INTO payouts (policy_id, user_address, amount, spike_id, tx_hash) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return DB.QueryRow(query, p.PolicyID, p.UserAddress, p.Amount, p.SpikeID, p.TxHash).Scan(&p.ID)
}

// UpdatePolicyStatus updates policy status
func UpdatePolicyStatus(policyID int, status string) error {
	query := `UPDATE policies SET status = $1 WHERE id = $2`
	_, err := DB.Exec(query, status, policyID)
	return err
}

// Close closes database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// GetRecentSpikes retrieves recent spike detection events
func GetRecentSpikes(limit int) ([]*Spike, error) {
	query := `SELECT s.id, s.timestamp, s.symbol, p.open, p.high, p.low, p.close, 
			         s.body_ratio, s.range_close_percent, s.detected_at 
			  FROM spikes s 
			  JOIN prices p ON s.price_id = p.id 
			  ORDER BY s.detected_at DESC LIMIT $1`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spikes []*Spike
	for rows.Next() {
		s := &Spike{}
		if err := rows.Scan(&s.ID, &s.Timestamp, &s.Symbol, &s.Open, &s.High, &s.Low, &s.Close,
			&s.BodyRatio, &s.RangeClosePercent, &s.DetectedAt); err != nil {
			return nil, err
		}
		spikes = append(spikes, s)
	}
	return spikes, nil
}

// GetRecentPrices retrieves recent price data
func GetRecentPrices(symbol string, limit int) ([]*PriceData, error) {
	// wrapper to unified GetPrices implementation
	return GetPrices(symbol, limit)
}

// helper: scan rows into []*PriceData
func scanPriceRows(rows *sql.Rows) ([]*PriceData, error) {
	var prices []*PriceData
	for rows.Next() {
		p := &PriceData{}
		if err := rows.Scan(&p.ID, &p.Timestamp, &p.Symbol, &p.Open, &p.High, &p.Low, &p.Close, &p.Volume); err != nil {
			return nil, err
		}
		prices = append(prices, p)
	}
	return prices, nil
}

// GetPrices returns prices for a symbol. If limit <= 0 all rows are returned (ordered asc), otherwise returns latest `limit` rows.
func GetPrices(symbol string, limit int) ([]*PriceData, error) {
	if limit > 0 {
		query := `SELECT id, timestamp, symbol, open, high, low, close, volume 
				  FROM prices WHERE symbol = $1 ORDER BY timestamp DESC LIMIT $2`
		rows, err := DB.Query(query, symbol, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanPriceRows(rows)
	}

	query := `SELECT id, timestamp, symbol, open, high, low, close, volume 
			  FROM prices WHERE symbol = $1 ORDER BY timestamp`
	rows, err := DB.Query(query, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPriceRows(rows)
}

// helper: scan rows into []*Payout
func scanPayoutRows(rows *sql.Rows) ([]*Payout, error) {
	var payouts []*Payout
	for rows.Next() {
		p := &Payout{}
		if err := rows.Scan(&p.ID, &p.PolicyID, &p.UserAddress, &p.Amount, &p.SpikeID, &p.TxHash, &p.ExecutedAt); err != nil {
			return nil, err
		}
		payouts = append(payouts, p)
	}
	return payouts, nil
}

// GetRecentPayouts retrieves recent payout records. If userAddress is non-empty it filters by that user.
func GetRecentPayouts(limit int, userAddress string) ([]*Payout, error) {
	if userAddress == "" {
		query := `SELECT id, policy_id, user_address, amount, spike_id, COALESCE(tx_hash, ''), executed_at 
				  FROM payouts ORDER BY executed_at DESC LIMIT $1`
		rows, err := DB.Query(query, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanPayoutRows(rows)
	}

	query := `SELECT id, policy_id, user_address, amount, spike_id, COALESCE(tx_hash, ''), executed_at
			  FROM payouts WHERE user_address = $1 ORDER BY executed_at DESC LIMIT $2`
	rows, err := DB.Query(query, userAddress, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPayoutRows(rows)
}

// SystemStats represents system statistics
type SystemStats struct {
	TotalSpikes    int `json:"total_spikes"`
	TotalPayouts   int `json:"total_payouts"`
	TotalPolicies  int `json:"total_policies"`
	ActivePolicies int `json:"active_policies"`
	TotalPrices    int `json:"total_prices"`
}

// GetSystemStats retrieves system statistics
func GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{}

	// Count spikes
	DB.QueryRow("SELECT COUNT(*) FROM spikes").Scan(&stats.TotalSpikes)

	// Count payouts
	DB.QueryRow("SELECT COUNT(*) FROM payouts").Scan(&stats.TotalPayouts)

	// Count total policies
	DB.QueryRow("SELECT COUNT(*) FROM policies").Scan(&stats.TotalPolicies)

	// Count active policies
	DB.QueryRow("SELECT COUNT(*) FROM policies WHERE status = 'active' AND expiry_time > NOW()").Scan(&stats.ActivePolicies)

	// Count price records
	DB.QueryRow("SELECT COUNT(*) FROM prices").Scan(&stats.TotalPrices)

	return stats, nil
}

// DeleteAllPrices deletes all price records
func DeleteAllPrices() error {
	query := `DELETE FROM prices`
	_, err := DB.Exec(query)
	return err
}

// DeleteAllSpikes deletes all spike records
func DeleteAllSpikes() error {
	query := `DELETE FROM spikes`
	_, err := DB.Exec(query)
	return err
}

// InsertPolicyFromEvent inserts a policy purchase event
func InsertPolicyFromEvent(userAddr common.Address, premium, coverage *big.Int, expiryTime *big.Int, txHash types.Log) error {
	// Convert wei to USDT (6 decimals)
	premiumFloat := new(big.Float).Quo(new(big.Float).SetInt(premium), big.NewFloat(1e6))
	coverageFloat := new(big.Float).Quo(new(big.Float).SetInt(coverage), big.NewFloat(1e6))
	expiry := time.Unix(expiryTime.Int64(), 0)

	premiumValue, _ := premiumFloat.Float64()
	coverageValue, _ := coverageFloat.Float64()

	query := `
		INSERT INTO policies
		(user_address, premium, coverage_amount, purchase_time, expiry_time, status, tx_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tx_hash) DO NOTHING
	`

	_, err := DB.Exec(query,
		userAddr.Hex(),
		premiumValue,
		coverageValue,
		time.Now(),
		expiry,
		"active",
		txHash.TxHash.Hex(),
	)

	return err
}

// InsertPayoutFromEvent inserts a payout execution event
func InsertPayoutFromEvent(userAddr common.Address, policyID *big.Int, amount *big.Int, txHash types.Log) error {
	// Convert wei to USDT (6 decimals)
	amountFloat := new(big.Float).Quo(new(big.Float).SetInt(amount), big.NewFloat(1e6))
	amountValue, _ := amountFloat.Float64()

	// Find the most recent active policy for this user
	var dbPolicyId int
	policyQuery := `
		SELECT id FROM policies
		WHERE user_address = $1 AND status = 'active'
		ORDER BY id DESC LIMIT 1
	`
	err := DB.QueryRow(policyQuery, userAddr.Hex()).Scan(&dbPolicyId)
	if err != nil {
		utils.LogError("Active policy not found for user %s, storing payout without policy link", userAddr.Hex())
		dbPolicyId = 0
	} else {
		// Update policy status to 'claimed'
		updateQuery := `UPDATE policies SET status = 'claimed' WHERE id = $1`
		_, err = DB.Exec(updateQuery, dbPolicyId)
		if err != nil {
			utils.LogError("Failed to update policy status: %v", err)
		}
	}

	// Insert into payouts table, ignore if tx_hash already exists to prevent duplicates
	payoutQuery := `
		INSERT INTO payouts
		(policy_id, user_address, amount, tx_hash, executed_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tx_hash) DO NOTHING
	`

	var policyIdPtr interface{}
	if dbPolicyId > 0 {
		policyIdPtr = dbPolicyId
	} else {
		policyIdPtr = nil
	}

	_, err = DB.Exec(payoutQuery,
		policyIdPtr,
		userAddr.Hex(),
		amountValue,
		txHash.TxHash.Hex(),
		time.Now(),
	)

	return err
}

// GetLastSyncedBlock retrieves the last synced block number for a contract
func GetLastSyncedBlock(contractAddr common.Address) (uint64, error) {
	var lastBlock uint64
	query := `SELECT last_synced_block FROM sync_state WHERE contract_address = $1`

	err := DB.QueryRow(query, contractAddr.Hex()).Scan(&lastBlock)
	if err != nil {
		// If no record exists, return 0 (will start from recent blocks)
		return 0, nil
	}

	return lastBlock, nil
}

// UpdateLastSyncedBlock updates the last synced block number for a contract
func UpdateLastSyncedBlock(contractAddr common.Address, blockNumber uint64) error {
	query := `
		INSERT INTO sync_state (contract_address, last_synced_block, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (contract_address)
		DO UPDATE SET last_synced_block = $2, updated_at = NOW()
	`

	_, err := DB.Exec(query, contractAddr.Hex(), blockNumber)
	return err
}

// GetAllBalances retrieves all rows from balances table
func GetAllBalances() ([]*Balance, error) {
	query := `SELECT id, token_address, user_address, balance, last_updated FROM balances`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []*Balance
	for rows.Next() {
		b := &Balance{}
		if err := rows.Scan(&b.ID, &b.TokenAddress, &b.UserAddress, &b.Balance, &b.LastUpdated); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}

// GetDistinctPolicyUsers returns a list of distinct user addresses found in policies table
func GetDistinctPolicyUsers() ([]string, error) {
	query := `SELECT DISTINCT user_address FROM policies`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// UpsertPolicy inserts or updates a policy record based on user_address + purchase_time
func UpsertPolicy(userAddr string, premium float64, coverage float64, purchaseTime time.Time, expiryTime time.Time, status string) error {
	var policyID sql.NullInt64
	err := DB.QueryRow(`SELECT id FROM policies WHERE user_address = $1 AND purchase_time = $2`, userAddr, purchaseTime).Scan(&policyID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if policyID.Valid {
		// Update existing
		_, err = DB.Exec(`UPDATE policies SET premium = $1, coverage_amount = $2, expiry_time = $3, status = $4 WHERE id = $5`,
			premium, coverage, expiryTime, status, policyID.Int64)
		return err
	}

	// Insert new
	_, err = DB.Exec(`INSERT INTO policies (user_address, premium, coverage_amount, purchase_time, expiry_time, status, tx_hash)
		VALUES ($1, $2, $3, $4, $5, $6, NULL)`, userAddr, premium, coverage, purchaseTime, expiryTime, status)
	return err
}

// FullSync updates balances and policies for all rows found in the DB.
// This was previously in a separate syncer package; moved here per request.
func FullSync(cfg *utils.Config) error {
	if cfg == nil || cfg.RPC.URL == "" {
		utils.LogInfo("RPC not configured; skipping full sync")
		return nil
	}

	client, err := ethclient.Dial(cfg.RPC.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	// Update balances for every row in balances table
	balances, err := GetAllBalances()
	if err != nil {
		utils.LogError("Failed to query balances for full sync: %v", err)
	} else {
		utils.LogInfo("FullSync: updating %d balance rows", len(balances))
		for _, b := range balances {
			if err := updateBalanceRow(client, cfg, b.TokenAddress, b.UserAddress); err != nil {
				utils.LogError("Failed to update balance for %s %s: %v", b.UserAddress, b.TokenAddress, err)
			}
		}
	}

	// Update policies for every distinct user in policies table
	users, err := GetDistinctPolicyUsers()
	if err != nil {
		utils.LogError("Failed to query policy users for full sync: %v", err)
	} else {
		utils.LogInfo("FullSync: updating policies for %d users", len(users))
		for _, u := range users {
			if err := syncPoliciesForUser(client, cfg, u); err != nil {
				utils.LogError("Failed to sync policies for user %s: %v", u, err)
			}
		}
	}

	return nil
}

// UpsertForUser updates balances and policies for a single user (called when frontend links wallet)
func UpsertForUser(cfg *utils.Config, userAddr string) error {
	if cfg == nil || cfg.RPC.URL == "" {
		utils.LogInfo("RPC not configured; skipping user upsert")
		return nil
	}

	client, err := ethclient.Dial(cfg.RPC.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	if err := updateBalanceRow(client, cfg, cfg.RPC.UsdtAddress, userAddr); err != nil {
		utils.LogError("Failed to update balance for %s %s: %v", userAddr, cfg.RPC.UsdtAddress, err)
		return err
	}

	// Update policies for this user
	if err := syncPoliciesForUser(client, cfg, userAddr); err != nil {
		utils.LogError("Failed to sync policies for user %s: %v", userAddr, err)
		return err
	}

	return nil
}

// updateBalanceRow reads on-chain ERC20 balance and upserts into DB
func updateBalanceRow(client *ethclient.Client, cfg *utils.Config, tokenAddr string, userAddr string) error {
	logger.Info("updateBalanceRow for %s %s", userAddr, tokenAddr)
	const erc20ABI = `[{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"}]`

	parsed, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return err
	}

	contract := bind.NewBoundContract(common.HexToAddress(tokenAddr), parsed, client, client, client)
	var out []interface{}
	if err := contract.Call(&bind.CallOpts{Context: context.Background()}, &out, "balanceOf", common.HexToAddress(userAddr)); err != nil {
		return err
	}

	balanceRaw := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	decimals := cfg.RPC.UsdtDecimals
	if decimals <= 0 {
		decimals = 6
	}

	denom := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	balf := new(big.Float).SetInt(balanceRaw)
	human := new(big.Float).Quo(balf, denom)
	humanVal, _ := human.Float64()

	if err := UpsertBalance(tokenAddr, userAddr, humanVal); err != nil {
		return err
	}
	return nil
}

// UpsertBalanceFromChain reads on-chain ERC20 balance for the given token and user using the provided
// ethclient.Client and configuration, then upserts the normalized balance into the DB.
// This exported wrapper lets other packages refresh a single user's balance for a token
// by querying the chain instead of relying on transfer deltas.
func UpsertBalanceFromChain(client *ethclient.Client, cfg *utils.Config, tokenAddr string, userAddr string) error {
	return updateBalanceRow(client, cfg, tokenAddr, userAddr)
}

// syncPoliciesForUser reads user's policies on-chain and upserts into DB
func syncPoliciesForUser(client *ethclient.Client, cfg *utils.Config, userAddr string) error {
	poolAddr := common.HexToAddress(cfg.RPC.ContractAddress)
	pool, err := contracts.NewInsurancePool(poolAddr, client)
	if err != nil {
		return err
	}

	// Get all user policies at once using typed binding for safety
	policies, err := pool.GetUserPolicies(&bind.CallOpts{Context: context.Background()}, common.HexToAddress(userAddr))
	if err != nil {
		return err
	}
	if len(policies) == 0 {
		return nil
	}
	decimals := cfg.RPC.UsdtDecimals

	for _, policy := range policies {
		userAddrOnChain := policy.User
		premiumRaw := policy.Premium
		coverageRaw := policy.CoverageAmount
		purchaseRaw := policy.PurchaseTime
		expiryRaw := policy.ExpiryTime
		active := policy.Active
		claimed := policy.Claimed

		premiumFloat := new(big.Float).Quo(new(big.Float).SetInt(premiumRaw), new(big.Float).SetFloat64(math.Pow10(int(decimals))))
		premiumVal, _ := premiumFloat.Float64()
		coverageFloat := new(big.Float).Quo(new(big.Float).SetInt(coverageRaw), new(big.Float).SetFloat64(math.Pow10(int(decimals))))
		coverageVal, _ := coverageFloat.Float64()

		purchaseTime := time.Unix(purchaseRaw.Int64(), 0)
		expiryTime := time.Unix(expiryRaw.Int64(), 0)

		status := "inactive"
		if active {
			status = "active"
		}
		if claimed {
			status = "claimed"
		}

		if err := UpsertPolicy(userAddrOnChain.Hex(), premiumVal, coverageVal, purchaseTime, expiryTime, status); err != nil {
			utils.LogError("UpsertPolicy failed for user %s: %v", userAddrOnChain.Hex(), err)
		}
	}

	return nil
}
