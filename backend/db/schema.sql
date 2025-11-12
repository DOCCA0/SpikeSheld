-- Table: prices - stores price data from CSV or Oracle
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    open DECIMAL(20, 8),
    high DECIMAL(20, 8),
    low DECIMAL(20, 8),
    close DECIMAL(20, 8) NOT NULL,
    volume DECIMAL(20, 8),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for prices table
CREATE INDEX IF NOT EXISTS idx_symbol_timestamp ON prices(symbol, timestamp);

-- Table: spikes - records detected price spikes
CREATE TABLE IF NOT EXISTS spikes (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    price_before DECIMAL(20, 8) NOT NULL,
    price_after DECIMAL(20, 8) NOT NULL,
    drop_percent DECIMAL(5, 2) NOT NULL,
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: policies - stores user insurance policies
CREATE TABLE IF NOT EXISTS policies (
    id SERIAL PRIMARY KEY,
    user_address VARCHAR(42) NOT NULL,
    premium DECIMAL(20, 8) NOT NULL,
    coverage_amount DECIMAL(20, 8) NOT NULL,
    purchase_time TIMESTAMP NOT NULL,
    expiry_time TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'active', -- active, expired, claimed
    tx_hash VARCHAR(66),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: payouts - logs payout executions
CREATE TABLE IF NOT EXISTS payouts (
    id SERIAL PRIMARY KEY,
    policy_id INTEGER REFERENCES policies(id),
    user_address VARCHAR(42) NOT NULL,
    amount DECIMAL(20, 8) NOT NULL,
    spike_id INTEGER REFERENCES spikes(id),
    tx_hash VARCHAR(66),
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_policies_user ON policies(user_address);
CREATE INDEX IF NOT EXISTS idx_policies_status ON policies(status);
CREATE INDEX IF NOT EXISTS idx_payouts_user ON payouts(user_address);
CREATE INDEX IF NOT EXISTS idx_spikes_symbol ON spikes(symbol);
