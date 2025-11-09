# 📋 SpikeShield Quick Reference

## Essential Commands

### Setup & Installation
```bash
# Initial setup
cp .env.example .env
./setup.sh              # Linux/Mac
setup.bat               # Windows

# Deploy contracts
cd contracts
npm install
npx hardhat compile
npx hardhat run scripts/deploy.js --network sepolia
```

### Running the System

**Option 1: Docker (Easiest)**
```bash
docker-compose up -d              # Start all services
docker-compose logs -f backend    # View backend logs
docker-compose down               # Stop all services
```

**Option 2: Manual**
```bash
# Terminal 1: Database
docker-compose up -d postgres

# Terminal 2: Backend (Replay Mode)
cd backend
go run main.go --mode replay --symbol BTCUSDT \
  --start "2021-05-19T00:00:00" \
  --end "2021-05-19T03:00:00"

# Terminal 2: Backend (Live Mode)
cd backend
go run main.go --mode live --symbol BTCUSDT

# Terminal 3: Frontend
cd frontend
npm start
```

## Important Addresses

Update these in `.env` after deployment:

```bash
# Contract addresses
INSURANCE_POOL_ADDRESS=0x...
USDT_ADDRESS=0x...

# RPC endpoints
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_KEY
BSC_TESTNET_RPC_URL=https://data-seed-prebsc-1-s1.binance.org:8545

# Private key (backend wallet)
PRIVATE_KEY=0x...
```

## Key Parameters

### Insurance Policy
- **Premium**: 10 USDT
- **Coverage**: 100 USDT  
- **Duration**: 24 hours
- **Spike Threshold**: 10% drop in 5 minutes

### Detection Settings (config.yaml)
```yaml
detector:
  threshold_percent: 10.0    # Trigger threshold
  window_minutes: 5          # Analysis window

chainlink:
  update_interval: 60        # Price update frequency (seconds)
```

## Database Queries

```bash
# Access database
psql -U postgres -d spikeshield

# Useful queries
SELECT * FROM prices ORDER BY timestamp DESC LIMIT 10;
SELECT * FROM spikes;
SELECT * FROM policies WHERE status = 'active';
SELECT * FROM payouts;

# Clear all data
TRUNCATE prices, spikes, policies, payouts RESTART IDENTITY CASCADE;
```

## Contract Interaction (via Frontend)

1. **Connect Wallet** → MetaMask extension
2. **Switch Network** → Sepolia or BSC Testnet
3. **Get Test Funds** → Faucet or "Mint Test USDT" button
4. **Buy Insurance** → "Buy Insurance (10 USDT)" button
5. **View Policies** → Scroll down to see policy cards

## Testing Spike Detection

### Perfect Demo Scenario
```bash
# This time range includes a 10% spike at 01:40
go run main.go --mode replay --symbol BTCUSDT \
  --start "2021-05-19T00:00:00" \
  --end "2021-05-19T03:00:00"

# Expected output:
# - Loads ~30 price records
# - Detects spike at 2021-05-19 01:40
# - Triggers payout automatically
```

## File Locations

### Smart Contracts
- `contracts/InsurancePool.sol` - Main insurance logic
- `contracts/MockUSDT.sol` - Test USDT token

### Backend
- `backend/main.go` - Entry point
- `backend/config.yaml` - Configuration
- `backend/detector/detect.go` - Spike detection algorithm
- `backend/api/payout.go` - Blockchain interaction
- `backend/datafeed/replay_feed.go` - CSV loader
- `backend/datafeed/live_feed.go` - Chainlink integration

### Frontend  
- `frontend/src/App.js` - Main UI
- `frontend/src/hooks/useContract.js` - Web3 logic
- `frontend/src/components/index.js` - UI components

### Data
- `data/btcusdt_2021-05-19.csv` - Historical price data

## Common Issues

| Issue | Solution |
|-------|----------|
| Contract deployment fails | Check testnet ETH/BNB balance |
| Backend can't connect to DB | Run `docker-compose up -d postgres` |
| Frontend shows wrong balance | Click "Refresh Data" button |
| No spike detected | Check time range includes 01:40 |
| Transaction reverts | Ensure USDT approved first |

## Network Information

### Sepolia Testnet
- **Chain ID**: 11155111
- **RPC**: https://sepolia.infura.io/v3/YOUR_KEY
- **Explorer**: https://sepolia.etherscan.io
- **Faucet**: https://sepoliafaucet.com

### BSC Testnet
- **Chain ID**: 97
- **RPC**: https://data-seed-prebsc-1-s1.binance.org:8545
- **Explorer**: https://testnet.bscscan.com
- **Faucet**: https://testnet.binance.org/faucet-smart

## Port Usage

- **Frontend**: http://localhost:3000
- **PostgreSQL**: localhost:5432
- **Backend**: No HTTP server (CLI only)

## Environment Variables

### Frontend (.env)
```bash
REACT_APP_INSURANCE_POOL_ADDRESS=0x...
REACT_APP_USDT_ADDRESS=0x...
```

### Backend (config.yaml or .env)
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=spikeshield
```

## Logs & Debugging

```bash
# Backend logs
go run main.go 2>&1 | tee backend.log

# Frontend logs (browser console)
# Press F12 in browser

# PostgreSQL logs
docker-compose logs postgres

# Contract events (Etherscan)
# Go to contract address → Events tab
```

## Demo Checklist

- [ ] Contracts deployed
- [ ] .env files updated
- [ ] Database running
- [ ] MetaMask installed & connected to testnet
- [ ] Test ETH/BNB in wallet
- [ ] Frontend running on :3000
- [ ] Backend ready to run
- [ ] CSV data in place
- [ ] Practice run completed

## Support Resources

- **Full Documentation**: README.md
- **Demo Guide**: DEMO_GUIDE.md  
- **Troubleshooting**: TROUBLESHOOTING.md
- **Project Summary**: PROJECT_COMPLETE.md

## Quick Test

Verify everything works:
```bash
# 1. Database
psql -U postgres -d spikeshield -c "SELECT 1"

# 2. Contracts (replace with your address)
curl -X POST https://sepolia.infura.io/v3/YOUR_KEY \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getCode","params":["0xYOUR_CONTRACT","latest"],"id":1}'

# 3. Frontend
curl http://localhost:3000

# 4. Backend
cd backend && go run main.go --help
```

---

**Pro Tip**: Keep this file open during your demo for quick reference! 🚀
