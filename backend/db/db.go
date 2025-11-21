package db

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"spikeshield/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	query := `SELECT id, timestamp, symbol, open, high, low, close, volume 
			  FROM prices WHERE symbol = $1 ORDER BY timestamp`

	rows, err := DB.Query(query, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	query := `SELECT id, timestamp, symbol, open, high, low, close, volume 
			  FROM prices WHERE symbol = $1 ORDER BY timestamp DESC LIMIT $2`

	rows, err := DB.Query(query, symbol, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// GetRecentPayouts retrieves recent payout records
func GetRecentPayouts(limit int) ([]*Payout, error) {
	query := `SELECT id, policy_id, user_address, amount, spike_id, COALESCE(tx_hash, ''), executed_at 
			  FROM payouts ORDER BY executed_at DESC LIMIT $1`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// GetActivePolicyIDForUser retrieves the most recent active policy ID for a user
func GetActivePolicyIDForUser(userAddr common.Address) (int, error) {
	var dbPolicyId int
	query := `
		SELECT id FROM policies
		WHERE user_address = $1 AND status = 'active'
		ORDER BY id DESC LIMIT 1
	`
	err := DB.QueryRow(query, userAddr.Hex()).Scan(&dbPolicyId)
	if err != nil {
		return 0, err
	}
	return dbPolicyId, nil
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
