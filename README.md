# SpikeShield 🛡️

Decentralized spike insurance protocol that protects users against sudden cryptocurrency price drops.

## 🌟 Features

- **Dual Mode Operation**
  - 📊 **Replay Mode**: Test with historical data (May 19, 2021 BTC crash)
  - ⚡ **Live Mode**: Real-time monitoring using Chainlink Oracle
  
- **Smart Contracts**
  - Purchase insurance with mock USDT
  - Automatic payout execution when spike detected
  - Configurable coverage parameters

- **Backend System**
  - Price data ingestion (CSV or Chainlink)
  - Spike detection algorithm
  - Automated on-chain payout triggering
  - PostgreSQL database for data persistence

- **Frontend DApp**
  - Web3 wallet integration (MetaMask)
  - Purchase insurance interface
  - View policy status and history
  - Real-time updates

## 📁 Project Structure

```
spikeshield/
├── contracts/              # Solidity smart contracts
│   ├── InsurancePool.sol   # Main insurance contract
│   ├── MockUSDT.sol        # Test USDT token
│   ├── hardhat.config.js   # Hardhat configuration
│   └── package.json
│
├── backend/                # Go backend service
│   ├── main.go            # Entry point
│   ├── config.yaml        # Configuration file
│   ├── db/                # Database layer
│   ├── datafeed/          # Price data feeds
│   ├── detector/          # Spike detection
│   ├── api/               # Blockchain interaction
│   └── utils/             # Helper functions
│
├── frontend/              # React frontend
│   ├── src/
│   │   ├── App.js        # Main component
│   │   ├── hooks/        # Web3 hooks
│   │   └── components/   # UI components
│   └── package.json
│
├── data/                  # Historical price data
│   └── btcusdt_2021-05-19.csv
│
└── docker-compose.yml     # One-command deployment
```

## 🚀 Quick Start

### Prerequisites

- Node.js 18+
- Go 1.21+
- PostgreSQL 15+
- MetaMask wallet
- Testnet ETH/BNB

### 1. Clone and Setup

```bash
git clone <repository>
cd SpikeShield
cp .env.example .env
```

### 2. Deploy Smart Contracts

```bash
cd contracts
npm install
npx hardhat compile

# Deploy to testnet (Sepolia or BSC Testnet)
npx hardhat run scripts/deploy.js --network sepolia

# Update .env with deployed contract addresses
```

### 3. Setup Backend

```bash
cd backend
go mod download

# Update config.yaml with your settings
# - Database connection
# - RPC URL
# - Contract addresses
# - Detection parameters

# Initialize database
psql -U postgres -d spikeshield -f db/schema.sql
```

### 4. Setup Frontend

```bash
cd frontend
npm install

# Update .env with contract addresses
echo "REACT_APP_INSURANCE_POOL_ADDRESS=<your_address>" >> .env
echo "REACT_APP_USDT_ADDRESS=<your_address>" >> .env
```

### 5. Run with Docker (Easiest)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### 6. Run Manually

**Terminal 1 - Database:**
```bash
# Start PostgreSQL (if not using Docker)
psql -U postgres
CREATE DATABASE spikeshield;
\c spikeshield
\i backend/db/schema.sql
```

**Terminal 2 - Backend:**
```bash
cd backend

# Replay mode (test with historical data)
go run main.go --mode replay --symbol BTCUSDT --start "2021-05-19T00:00:00" --end "2021-05-19T03:00:00"

# Live mode (real-time monitoring)
go run main.go --mode live --symbol BTCUSDT
```

**Terminal 3 - Frontend:**
```bash
cd frontend
npm start
# Opens http://localhost:3000
```

## 🎬 Demo Flow (Hackathon Presentation)

### Scenario: Protect Against May 19, 2021 BTC Flash Crash

1. **Connect Wallet**
   - Open frontend at http://localhost:3000
   - Click "Connect Wallet"
   - Connect MetaMask to testnet

2. **Get Test Funds**
   - Click "Mint 100 Test USDT"
   - Verify balance shows 100 USDT

3. **Purchase Insurance**
   - Click "Buy Insurance (10 USDT)"
   - Approve transaction in MetaMask
   - Wait for confirmation
   - See policy appear in "Your Insurance Policies"

4. **Run Backend Replay**
   ```bash
   go run main.go --mode replay --symbol BTCUSDT --start "2021-05-19T00:00:00" --end "2021-05-19T03:00:00"
   ```

5. **Observe Spike Detection**
   - Backend loads historical data
   - Detects 10%+ price drop at 01:40 (45000 → 40500)
   - Automatically triggers smart contract payout

6. **Verify Payout**
   - Refresh frontend
   - Policy status changes to "Claimed"
   - Balance increases by 100 USDT (coverage amount)
   - Transaction hash displayed

## 📊 Configuration

### Backend (config.yaml)

```yaml
detector:
  threshold_percent: 10.0   # Trigger at 10% drop
  window_minutes: 5         # Check last 5 minutes

chainlink:
  btc_usd_feed: "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43"  # Sepolia
  update_interval: 60       # Fetch every 60 seconds
```

### Smart Contract Parameters

- Premium: 10 USDT
- Coverage: 100 USDT
- Duration: 24 hours
- Spike Threshold: 10% drop within 5 minutes

## 🧪 Testing

### Test Replay Mode
```bash
# Test with May 19, 2021 crash data
go run main.go --mode replay --symbol BTCUSDT --start "2021-05-19T00:00:00" --end "2021-05-19T03:00:00"
```

### Test Live Mode
```bash
# Connect to Chainlink price feed
go run main.go --mode live --symbol BTCUSDT
```

### Contract Testing
```bash
cd contracts
npx hardhat test
```

## 📝 Database Schema

| Table | Description |
|-------|-------------|
| `prices` | Price data from CSV or Oracle |
| `spikes` | Detected spike events |
| `policies` | User insurance policies |
| `payouts` | Executed payout records |

## 🔧 Troubleshooting

### Backend won't start
- Check PostgreSQL is running
- Verify config.yaml database settings
- Run `go mod tidy` to fix dependencies

### Frontend can't connect to wallet
- Install MetaMask extension
- Switch to correct testnet
- Check contract addresses in .env

### Payout not executing
- Verify backend has oracle role in contract
- Check RPC URL is correct
- Ensure contract has sufficient USDT balance

## 🌐 Supported Networks

- Ethereum Sepolia Testnet
- BSC Testnet
- Any EVM-compatible testnet

## 📚 Key Technologies

- **Blockchain**: Solidity, Hardhat, OpenZeppelin
- **Backend**: Go, PostgreSQL
- **Frontend**: React, ethers.js
- **Oracle**: Chainlink Price Feeds
- **DevOps**: Docker, Docker Compose

## 🎯 Future Enhancements (Post-Hackathon)

- [ ] Support multiple assets (ETH, SOL, etc.)
- [ ] Dynamic premium calculation based on volatility
- [ ] NFT-based policy representation
- [ ] Liquidity pool mechanism
- [ ] Advanced charting and analytics
- [ ] Mobile app
- [ ] Mainnet deployment

## 📄 License

MIT

## 👥 Team

Built for hackathon demo purposes.

## 🙏 Acknowledgments

- Chainlink for price oracle infrastructure
- OpenZeppelin for secure contract libraries
- May 19, 2021 for providing excellent test data 😅

---

**⚠️ Disclaimer**: This is a hackathon MVP for demonstration purposes only. Not audited. Do not use with real funds.
