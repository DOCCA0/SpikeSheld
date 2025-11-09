# 🎬 SpikeShield Hackathon Demo Guide

Complete step-by-step guide for presenting SpikeShield at your hackathon.

## 📋 Pre-Demo Checklist

### Before the Event

- [ ] Deploy contracts to testnet (Sepolia or BSC Testnet)
- [ ] Update all `.env` files with contract addresses
- [ ] Test the entire flow end-to-end
- [ ] Prepare MetaMask with testnet ETH/BNB
- [ ] Start services (database, backend, frontend)
- [ ] Open browser tabs for demo
- [ ] Prepare slides (optional)

### Required Tabs

1. Frontend: `http://localhost:3000`
2. Terminal: Backend logs
3. Etherscan/BscScan: To show transactions
4. MetaMask: Extension popup

## 🎯 Demo Script (5-10 minutes)

### Introduction (1 minute)

**Opening:**
> "Hi! I'm presenting SpikeShield - a decentralized insurance protocol that protects cryptocurrency traders against sudden price crashes, also known as 'spikes'."

**Problem Statement:**
> "On May 19, 2021, Bitcoin crashed from $45,000 to $30,000 in hours, liquidating millions of traders. What if you could buy insurance against such events?"

**Solution:**
> "SpikeShield allows users to purchase spike insurance. If the price drops more than 10% within 5 minutes, the protocol automatically pays out your coverage."

### Architecture Overview (1 minute)

**Show diagram or describe:**
> "The system has three components:
> 1. **Smart Contracts** on Ethereum testnet - handle insurance purchases and payouts
> 2. **Backend Service** in Go - monitors prices and detects spikes
> 3. **Frontend DApp** in React - user interface for buying insurance
>
> It supports two modes:
> - **Replay Mode**: Test with historical data from the May 2021 crash
> - **Live Mode**: Real-time monitoring using Chainlink price oracle"

### Live Demo (6-7 minutes)

#### Step 1: Connect Wallet (30 seconds)

**Action:**
- Open frontend
- Click "Connect Wallet"
- Approve MetaMask connection

**Say:**
> "First, I'll connect my wallet to the DApp. I'm using MetaMask on Sepolia testnet."

#### Step 2: Get Test Funds (30 seconds)

**Action:**
- Click "Mint 100 Test USDT"
- Show balance update to 100 USDT

**Say:**
> "For this demo, I'll mint some test USDT tokens. In production, users would use real USDT."

#### Step 3: Review Insurance Terms (30 seconds)

**Action:**
- Point to the insurance info cards

**Say:**
> "The insurance costs 10 USDT as premium and provides 100 USDT coverage for 24 hours. 
> If Bitcoin drops 10% or more within a 5-minute window, the payout triggers automatically."

#### Step 4: Purchase Insurance (1 minute)

**Action:**
- Click "Buy Insurance (10 USDT)"
- Approve USDT spending in MetaMask
- Approve purchase transaction
- Wait for confirmation
- Show policy card appearing

**Say:**
> "Now I'll purchase the insurance. This requires two transactions:
> 1. First, approve the contract to spend my USDT
> 2. Second, execute the purchase
>
> And there we go - my policy is now active!"

#### Step 5: Run Replay Demo (2 minutes)

**Action:**
- Switch to terminal
- Run command:
  ```bash
  go run main.go --mode replay --symbol BTCUSDT --start "2021-05-19T00:00:00" --end "2021-05-19T03:00:00"
  ```
- Show logs loading price data
- Point out spike detection message
- Highlight payout execution log

**Say:**
> "Now let me demonstrate the detection system. I'll run the backend in replay mode 
> using real historical data from the May 19, 2021 Bitcoin crash.
>
> Watch the logs... The system is:
> 1. Loading price data from CSV
> 2. Analyzing for price drops
> 3. Detecting a spike at 01:40 UTC - price dropped from $45,000 to $40,500
> 4. That's an 10% drop - trigger condition met!
> 5. Automatically executing the payout transaction"

#### Step 6: Verify Payout (1 minute)

**Action:**
- Switch back to frontend
- Click "Refresh Data"
- Show policy status changed to "Claimed"
- Show balance increased from 90 to 190 USDT (90 remaining + 100 payout)
- Show payout transaction hash

**Say:**
> "Back to the frontend, I'll refresh the data. Look at that!
> - Policy status is now 'Claimed'
> - My balance increased from 90 to 190 USDT
> - The 100 USDT coverage was automatically paid out
> - We can verify the transaction on Etherscan using this hash"

#### Step 7: Show Live Mode (Optional, 1 minute)

**Action:**
- In terminal, run:
  ```bash
  go run main.go --mode live --symbol BTCUSDT
  ```
- Show logs fetching from Chainlink

**Say:**
> "The system also works in live mode, fetching real-time prices from Chainlink oracles.
> In production, this would continuously monitor the market and trigger payouts instantly
> when spikes occur."

### Technical Highlights (1 minute)

**Key Points:**
> "Technical highlights of SpikeShield:
> 
> - **Smart Contracts**: Built with Solidity and OpenZeppelin for security
> - **Oracle Integration**: Uses Chainlink for decentralized price feeds
> - **Backend**: Go service for high-performance price analysis
> - **Database**: PostgreSQL stores all price data and policy records
> - **Frontend**: React with ethers.js for seamless Web3 integration
> - **Deployment**: Docker Compose for one-command setup"

### Closing (30 seconds)

**Future Plans:**
> "Future enhancements include:
> - Support for multiple cryptocurrencies
> - Dynamic pricing based on volatility
> - Liquidity pool mechanism for sustainable payouts
> - Mobile application"

**Call to Action:**
> "Thank you! I'd be happy to answer any questions about the architecture, 
> the spike detection algorithm, or the smart contracts."

## 🎤 Q&A Preparation

### Likely Questions & Answers

**Q: How do you prevent false positives in spike detection?**
A: We use a 5-minute rolling window and require 10% threshold. We also compare against the highest price in the window, not just the previous candle. In production, we'd add more sophisticated filters like volume analysis.

**Q: What prevents the insurance pool from running out of funds?**
A: In this MVP, the pool is manually funded. In production, we'd implement:
- Risk-based premium pricing
- Reinsurance mechanisms
- Liquidity provider incentives
- Maximum coverage caps

**Q: Why not use existing DeFi insurance protocols?**
A: Existing protocols focus on smart contract risks. SpikeShield specifically targets price volatility insurance, which is a different market. It's more like options/derivatives but simplified for retail users.

**Q: How accurate is the Chainlink oracle?**
A: Chainlink aggregates data from multiple sources and updates every ~30 seconds on mainnet. For production, we'd use multiple oracles and implement outlier detection.

**Q: What's the business model?**
A: Multiple revenue streams:
- Premium collection (10% goes to protocol)
- Liquidity provider fees
- Staking rewards for governance token holders
- Enterprise API access

**Q: Gas costs?**
A: On Ethereum L1, buying insurance costs ~$5-20 in gas. We'd deploy on L2s (Arbitrum, Optimism) where costs are <$0.50. The payout is executed by our backend, so users don't pay gas for claims.

**Q: Could this be gamed by users?**
A: Potential attacks and mitigations:
- **Flash crash manipulation**: We use time-weighted average prices
- **Oracle manipulation**: Multiple oracle sources + outlier detection  
- **Front-running**: Policies have a 1-minute activation delay (not in MVP)
- **Sybil attacks**: KYC for large policies (future feature)

## 🐛 Common Demo Issues

### Issue: MetaMask won't connect
**Solution:** 
- Check you're on correct network
- Hard refresh page (Ctrl+Shift+R)
- Reconnect MetaMask to site

### Issue: Transaction failing
**Solution:**
- Ensure sufficient testnet ETH for gas
- Check contract addresses in .env match deployment
- Verify contract has USDT balance

### Issue: Backend can't connect to database
**Solution:**
- Run `docker-compose up -d postgres`
- Check config.yaml database settings
- Verify PostgreSQL is running on port 5432

### Issue: No spike detected in replay
**Solution:**
- Check time range includes 01:40 on May 19
- Verify CSV data loaded (check logs)
- Ensure threshold is 10% or lower

## 📊 Demo Data

### May 19, 2021 Timeline

- **00:00 - 01:30**: Normal trading, price ~$43,000-45,000
- **01:40**: SPIKE! Price drops to $40,500 (10% drop)
- **02:00+**: Continued volatility, stays around $39,000-40,000

### Sample Commands

```bash
# Quick replay (shows spike)
go run main.go --mode replay --symbol BTCUSDT --start "2021-05-19T00:00:00" --end "2021-05-19T03:00:00"

# Live mode (use for live price monitoring)
go run main.go --mode live --symbol BTCUSDT

# Check database records
psql -U postgres -d spikeshield -c "SELECT * FROM spikes;"
psql -U postgres -d spikeshield -c "SELECT * FROM payouts;"
```

## 🎥 Recording Tips

If recording video demo:
- Use 1080p resolution
- Zoom browser to 110% for visibility
- Use dark mode for terminal
- Record audio clearly
- Practice 2-3 times before final take
- Keep it under 10 minutes
- Add captions/annotations for key moments

## 🏆 Judging Criteria Alignment

**Innovation:**
- Novel approach to crypto insurance
- Dual-mode architecture
- Automated payout system

**Technical Implementation:**
- Full-stack: Smart contracts, backend, frontend
- Multiple technologies integrated
- Docker deployment ready

**User Experience:**
- Simple one-click insurance purchase
- Clear visual feedback
- Automatic payout (no claims process)

**Market Potential:**
- $500B+ crypto derivatives market
- Growing demand for retail protection
- Clear business model

Good luck with your demo! 🚀
