# SpikeShield 🛡️

Decentralized wick insurance protocol that protects users against cryptocurrency price wicks (small body, large range candles).

**Updated: December 2025**

## 🌟 Features

- **Dual Mode Operation**
  - 📊 **Replay Mode**: Monitor DB inserts for testing (feed data via API)
  - ⚡ **Live Mode**: Real-time Chainlink price monitoring
- **Smart Contracts** (Upgradeable)
  - Purchase insurance with MockUSDT
  - Automatic payout on wick detection
  - Transparent proxy pattern for upgrades
- **Backend API Server**
  - Price data ingestion (API POST or Chainlink)
  - Wick detection (body ratio ≤30%, range/close ≥10%)
  - Automated on-chain payouts
  - Event & token transfer polling
  - PostgreSQL with balances/policies sync
  - REST API for frontend (localhost:8080)
- **Frontend DApp**
  - Web3 integration (ethers v6)
  - Buy insurance, mint test USDT
  - Real-time policies, balances, charts
  - System stats, recent wicks/payouts

## 📁 Project Structure

```
SpikeShield/
├── contracts/                    # Solidity (Hardhat + Upgrades)
│   ├── contracts/
│   │   ├── InsurancePool.sol     # Upgradeable insurance logic
│   │   └── MockUSDT.sol          # Test ERC20
│   ├── scripts/
│   │   ├── deploy_init.js        # Initial deploy
│   │   └── deploy_upgrade.js     # Proxy upgrades
│   ├── test/                     # Jest tests
│   ├── hardhat.config.js
│   ├── package.json
│   └── USAGE.md
│
├── backend/                      # Go API service
│   ├── .env.example              # PRIVATE_KEY env
│   ├── main.go                   # Entry point (--mode replay/live)
│   ├── config.yaml               # Config (DB, RPC, detector...)
│   ├── api/                      # HTTP API (prices, stats...)
│   ├── contracts/                # ABI bindings
│   ├── datafeed/                 # Live/replay feeds
│   ├── db/                       # PostgreSQL ORM
│   │   └── schema.sql
│   ├── detector/                 # Wick detection
│   ├── eventlistener/            # Event polling
│   └── utils/
│
├── frontend/                     # React 18 DApp
│   ├── src/
│   │   ├── App.js                # Main UI
│   │   ├── components/
│   │   │   ├── PriceChart.js
│   │   │   └── PayoutNotification.js
│   │   ├── hooks/
│   │   │   └── useContract.js
│   │   └── services/
│   │       └── api.js            # Backend API client
│   ├── package.json              # ethers v6, lightweight-charts
│   └── Dockerfile
│
├── data/                         # Test data
│   └── btcusdt_wick_test.csv
│
├── docker-compose.yml            # Postgres + Backend + Frontend
├── setup.sh / setup.bat          # Quick setup
└── .env.example                  # Docker env vars
```

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.24+
- PostgreSQL 15+ (or Docker)
- MetaMask (Sepolia testnet)
- Testnet ETH

### 1. Clone & Setup
```bash
git clone <repo> && cd SpikeShield
cp .env.example .env
./setup.sh  # or setup.bat on Windows
```

### 2. Deploy Contracts (Sepolia)
```bash
cd contracts
npm install
npx hardhat compile

# Deploy initial proxy
npx hardhat run scripts/deploy_init.js --network sepolia

# Copy addresses to backend/config.yaml & frontend/.env
# rpc.contract_address, rpc.usdt_address
```

### 3. Backend Setup
```bash
cd backend
cp .env.example .env
go mod tidy

# Edit config.yaml:
# - database creds
# - rpc.url, contract_address, usdt_address
# - rpc.private_key: "${PRIVATE_KEY}" (loads from .env)
# - chainlink.btc_usd_feed (Sepolia: 0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43)
```

**Database Init** (if not Docker):
```bash
psql -U postgres -d spikeshield -f db/schema.sql
```

### 4. Frontend Setup
```bash
cd frontend
npm install
# Add to .env:
# REACT_APP_INSURANCE_POOL_ADDRESS=0x...
# REACT_APP_USDT_ADDRESS=0x...
```

### 5. Docker (Recommended - One Command)
```bash
docker compose up -d  # Starts DB+Backend(replay)+Frontend:3000
docker compose logs -f backend  # Watch logs
docker compose down
```
*Note: For Docker backend payouts, set PRIVATE_KEY=... in root .env. Frontend: REACT_APP_* in frontend/.env*

### 6. Manual Run
**T1 - DB** (if not Docker): `psql ... schema.sql`

**T2 - Backend**:
```bash
cd backend
# Replay (monitor DB inserts)
go run main.go --mode replay --symbol BTCUSDT --api-port 8080

# Live (Chainlink)
go run main.go --mode live --symbol BTCUSDT
```

**T3 - Frontend**: `cd frontend && npm start` (localhost:3000)

## 🎬 Demo Flow

1. **Start Services**: `docker compose up -d`
2. **Frontend** (localhost:3000):
   - Connect MetaMask (Sepolia)
   - Mint 100 Test USDT
   - Buy Insurance (10 USDT → 100 coverage, 24h)
3. **Simulate Replay**:
   - POST CSV prices to http://localhost:8080/prices (dev endpoint)
   - Or manual DB insert wick candle
4. **Observe**:
   - Backend detects wick (body≤30%, range≥10%)
   - Auto triggers payout
   - Frontend shows "claimed", balance +100 USDT
5. **Verify**: Etherscan tx hash

## ⚙️ Configuration

### backend/config.yaml
```yaml
database:
  host: localhost  # postgres for Docker
  port: 5432
  user: postgres
  password: postgres
  dbname: spikeshield

rpc:
  url: https://ethereum-sepolia-rpc.publicnode.com
  contract_address: "0x..."
  private_key: "${PRIVATE_KEY}"  # Payout signer
  usdt_address: "0x..."
  usdt_decimals: 6

detector:
  threshold_percent: 0.1     # (high-low)/close >= 10%
  body_ratio_max: 0.3        # |open-close|/(high-low) <= 30%

chainlink:
  btc_usd_feed: "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43"  # Sepolia BTC/USD
  update_interval: 60

eventlistener:
  enabled: true
  poll_interval: 1  # seconds

mode: replay
```

### Insurance Parameters
- Premium: 10 USDT
- Coverage: 100 USDT
- Duration: 24 hours
- Wick: Body ratio ≤30% + Range ratio ≥10%

## 🧪 Testing

```bash
# Contracts
cd contracts && npx hardhat test

# Backend replay (feed CSV manually)
go run main.go --mode replay

# Live mode
go run main.go --mode live
```

## 📊 Database Schema

| Table       | Description                  |
|-------------|------------------------------|
| `prices`    | OHLCV candles                |
| `spikes`    | Detected wicks               |
| `policies`  | User policies                |
| `payouts`   | Executed payouts             |
| `balances`  | ERC20 balance cache          |
| `sync_state`| Event sync tracking          |

## 🔧 Troubleshooting

- **Backend DB error**: Check config.yaml creds, Postgres running
- **No prices**: POST to /api/prices or enable live mode
- **Wallet not linked**: Frontend auto-links on connect
- **Payout fails**: Check private_key, contract USDT balance, oracle role
- **Events missing**: Enable eventlistener, check poll_interval

## 🌐 Networks
- Sepolia (primary)
- BSC Testnet
- Local Hardhat

## 📚 Tech Stack
- **Contracts**: Solidity ^0.8, OpenZeppelin Upgrades, Hardhat
- **Backend**: Go 1.24, Gin, go-ethereum, PostgreSQL
- **Frontend**: React 18, ethers 6, Lightweight Charts
- **Oracle**: Chainlink Price Feeds
- **DevOps**: Docker Compose

## 🎯 Enhancements
- [x] Event polling & balance sync
- [x] API server + charts
- [ ] Multi-asset
- [ ] Dynamic pricing
- [ ] Mainnet
- [ ] Audits

## 📄 License
MIT

## 👥 Team
Hackathon prototype → Production-ready MVP

---
**⚠️ Testnet only. Not audited. Demo purposes.**
