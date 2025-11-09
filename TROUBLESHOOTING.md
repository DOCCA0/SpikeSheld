# 🔧 SpikeShield Troubleshooting Guide

## Common Issues and Solutions

### 1. Contract Deployment Issues

#### Error: "Invalid API Key" or "Network error"
**Problem:** Hardhat can't connect to RPC endpoint
**Solution:**
```bash
# Check your .env file
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_ACTUAL_KEY

# Test connection
curl https://sepolia.infura.io/v3/YOUR_KEY -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

#### Error: "Insufficient funds for gas"
**Problem:** Deployer account has no testnet ETH
**Solution:**
- Get Sepolia ETH from faucet: https://sepoliafaucet.com/
- Or BSC testnet BNB: https://testnet.binance.org/faucet-smart

#### Error: "Contract verification failed"
**Problem:** Etherscan verification issues
**Solution:**
```bash
# Flatten contract first
npx hardhat flatten contracts/InsurancePool.sol > InsurancePool_flat.sol
# Manually verify on Etherscan
```

---

### 2. Backend Issues

#### Error: "Failed to connect to database"
**Problem:** PostgreSQL not running or wrong credentials
**Solution:**
```bash
# Check if PostgreSQL is running
docker-compose ps

# Start PostgreSQL
docker-compose up -d postgres

# Check logs
docker-compose logs postgres

# Test connection manually
psql -h localhost -U postgres -d spikeshield
```

#### Error: "Failed to open CSV file"
**Problem:** CSV path is incorrect
**Solution:**
```bash
# Check file exists
ls -la data/btcusdt_2021-05-19.csv

# Use absolute path in config or run from correct directory
cd backend
go run main.go --mode replay
```

#### Error: "Failed to parse timestamp"
**Problem:** CSV date format doesn't match parser
**Solution:**
- Check CSV format matches: `2021-05-19 00:00:00`
- Update `parseTimestamp()` function in `replay_feed.go` if needed

#### Error: "Failed to get gas price"
**Problem:** RPC connection issues
**Solution:**
```bash
# Test RPC
curl -X POST https://sepolia.infura.io/v3/YOUR_KEY \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}'

# Use alternative RPC
# Sepolia: https://rpc.sepolia.org
# BSC Testnet: https://data-seed-prebsc-1-s1.binance.org:8545
```

#### Error: "No spikes detected"
**Problem:** Threshold too high or time window wrong
**Solution:**
```yaml
# Adjust in config.yaml
detector:
  threshold_percent: 5.0  # Lower threshold
  window_minutes: 10      # Wider window
```

---

### 3. Frontend Issues

#### Error: "Please install MetaMask"
**Problem:** MetaMask not installed
**Solution:**
- Install from: https://metamask.io/
- Refresh page after installation

#### Error: "Wrong network"
**Problem:** MetaMask on wrong network
**Solution:**
1. Click MetaMask extension
2. Click network dropdown (top)
3. Select "Sepolia" or add custom network:
   - Network Name: BSC Testnet
   - RPC URL: https://data-seed-prebsc-1-s1.binance.org:8545
   - Chain ID: 97
   - Currency: BNB
   - Block Explorer: https://testnet.bscscan.com

#### Error: "Transaction failed"
**Problem:** Various reasons
**Solution:**
```javascript
// Check contract addresses in .env
REACT_APP_INSURANCE_POOL_ADDRESS=0x... // Must match deployed
REACT_APP_USDT_ADDRESS=0x...           // Must match deployed

// Check MetaMask:
// 1. Sufficient gas (ETH/BNB)
// 2. Connected account
// 3. Correct network
```

#### Error: "Module not found"
**Problem:** Dependencies not installed
**Solution:**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

#### Error: Frontend shows 0 balance after minting
**Problem:** Need to refresh
**Solution:**
- Click "Refresh Data" button
- Or wait 5 seconds and auto-refresh should trigger

---

### 4. Docker Issues

#### Error: "Port already in use"
**Problem:** Services running on same ports
**Solution:**
```bash
# Find process using port
lsof -i :5432  # PostgreSQL
lsof -i :3000  # Frontend

# Kill process
kill -9 <PID>

# Or change ports in docker-compose.yml
ports:
  - "5433:5432"  # Use different host port
```

#### Error: "Cannot connect to Docker daemon"
**Problem:** Docker not running
**Solution:**
```bash
# Linux
sudo systemctl start docker

# Windows/Mac
# Start Docker Desktop application
```

#### Error: "Container keeps restarting"
**Problem:** Check container logs
**Solution:**
```bash
# View logs
docker-compose logs -f backend
docker-compose logs -f postgres

# Common fixes:
# 1. Database not ready - add healthcheck delay
# 2. Config file missing - mount volume correctly
# 3. Dependencies missing - rebuild image
docker-compose build --no-cache backend
```

---

### 5. Smart Contract Interaction Issues

#### Error: "Execution reverted"
**Problem:** Contract function requirements not met
**Solution:**
```solidity
// Common causes:

// 1. Insufficient USDT balance
// → Mint more test USDT

// 2. Not approved USDT spending
// → Call approve() first

// 3. Policy already active
// → Wait for expiry or use different account

// 4. Pool insufficient balance
// → Fund pool: insurance.fundPool(amount)

// Debug in MetaMask:
// Click failed transaction → View on Etherscan → "State" tab
```

#### Error: "Gas estimation failed"
**Problem:** Transaction would fail
**Solution:**
```javascript
// 1. Check function parameters
// 2. Verify contract state
// 3. Try with more gas
const tx = await contract.buyInsurance({
  gasLimit: 300000 // Manually set
});
```

---

### 6. Database Issues

#### Error: "Relation does not exist"
**Problem:** Tables not created
**Solution:**
```bash
# Run schema manually
psql -U postgres -d spikeshield -f backend/db/schema.sql

# Or in psql:
\i /path/to/backend/db/schema.sql
```

#### Error: "Duplicate key violation"
**Problem:** Trying to insert duplicate data
**Solution:**
```sql
-- Clear tables for fresh start
TRUNCATE prices, spikes, policies, payouts RESTART IDENTITY CASCADE;
```

#### Error: "Too many connections"
**Problem:** Connection pool exhausted
**Solution:**
```go
// In db.go, add connection pool settings
DB.SetMaxOpenConns(25)
DB.SetMaxIdleConns(5)
DB.SetConnMaxLifetime(5 * time.Minute)
```

---

### 7. Go Build Issues

#### Error: "Cannot find package"
**Problem:** Dependencies not downloaded
**Solution:**
```bash
cd backend
go mod download
go mod tidy
go mod verify
```

#### Error: "Module declares its path as X but was required as Y"
**Problem:** Module path mismatch
**Solution:**
```bash
# Fix go.mod
module spikeshield

# Then
go mod tidy
```

---

### 8. Replay Mode Issues

#### No data loaded
**Problem:** CSV not found or time range wrong
**Solution:**
```bash
# Check CSV exists
cat data/btcusdt_2021-05-19.csv | head

# Use correct time range (data is from 2021-05-19)
--start "2021-05-19T00:00:00" --end "2021-05-19T03:00:00"
```

#### Spike not detected
**Problem:** Spike outside time window
**Solution:**
```bash
# The spike occurs around 01:40
# Ensure your time range includes it:
--start "2021-05-19T00:00:00" --end "2021-05-19T02:00:00"

# Or check threshold
# config.yaml: threshold_percent: 10.0
```

---

### 9. Live Mode Issues

#### Error: "Failed to call latestRoundData"
**Problem:** Chainlink integration issue
**Solution:**
```go
// For MVP, the live mode is simplified
// To fix properly, generate contract bindings:

// 1. Get Chainlink ABI
// 2. Generate bindings
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
abigen --abi aggregator.json --pkg datafeed --type Aggregator --out aggregator.go

// 3. Use generated code in live_feed.go
```

#### Price not updating
**Problem:** Polling interval too long
**Solution:**
```yaml
# config.yaml
chainlink:
  update_interval: 30  # Update every 30 seconds
```

---

## Performance Optimization

### Slow database queries
```sql
-- Add indexes
CREATE INDEX idx_prices_symbol_timestamp ON prices(symbol, timestamp);
CREATE INDEX idx_spikes_timestamp ON spikes(timestamp);

-- Analyze query performance
EXPLAIN ANALYZE SELECT * FROM prices WHERE symbol = 'BTCUSDT' AND timestamp > NOW() - INTERVAL '1 hour';
```

### High memory usage (Go)
```go
// Batch process large datasets
const batchSize = 1000
for i := 0; i < len(prices); i += batchSize {
    batch := prices[i:min(i+batchSize, len(prices))]
    processBatch(batch)
}
```

---

## Debug Mode

### Enable verbose logging

**Backend:**
```go
// In utils/helpers.go, add debug flag
var Debug = os.Getenv("DEBUG") == "true"

func LogDebug(format string, args ...interface{}) {
    if Debug {
        log.Printf("[DEBUG] "+format, args...)
    }
}

// Run with:
DEBUG=true go run main.go
```

**Frontend:**
```javascript
// In useContract.js
console.log("Transaction:", tx);
console.log("Receipt:", receipt);
console.log("Contract address:", INSURANCE_POOL_ADDRESS);
```

---

## Health Checks

### System Health Script
```bash
#!/bin/bash
# check_health.sh

echo "🔍 SpikeShield Health Check"
echo "=========================="

# Database
echo -n "PostgreSQL: "
pg_isready -h localhost -U postgres && echo "✅" || echo "❌"

# Contracts (check on Etherscan)
echo -n "Contracts deployed: "
# Add your check here

# Frontend
echo -n "Frontend: "
curl -s http://localhost:3000 > /dev/null && echo "✅" || echo "❌"

echo "=========================="
```

---

## Getting Help

### Where to Find Answers

1. **Contract Issues**: Hardhat docs - https://hardhat.org/
2. **Web3 Issues**: ethers.js docs - https://docs.ethers.org/
3. **Go Issues**: Go docs - https://go.dev/doc/
4. **React Issues**: React docs - https://react.dev/

### Debug Information to Collect

When asking for help, provide:
```bash
# System info
go version
node --version
docker --version

# Error logs
docker-compose logs backend > backend.log
docker-compose logs postgres > postgres.log

# Contract addresses
cat .env | grep ADDRESS

# Transaction hash (if failed)
# Network (Sepolia/BSC Testnet)
# Account address
```

---

## Emergency Reset

### Complete Fresh Start
```bash
# Stop everything
docker-compose down -v

# Clean database
docker volume rm spikeshield_postgres_data

# Clean Go cache
go clean -cache -modcache -i -r

# Clean npm
cd frontend && rm -rf node_modules package-lock.json
cd ../contracts && rm -rf node_modules package-lock.json

# Rebuild
./setup.sh

# Redeploy contracts
cd contracts && npx hardhat run scripts/deploy.js --network sepolia

# Update .env with new addresses
# Restart everything
docker-compose up -d
```

---

Remember: This is an MVP for demonstration. Some issues are expected in development. The important thing is that the core demo flow works! 🚀
