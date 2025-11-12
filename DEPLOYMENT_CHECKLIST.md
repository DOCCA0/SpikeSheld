# Deployment Checklist ✅

Use this checklist to deploy and test the real contract integration.

## Pre-Deployment

### Contracts
- [ ] Contracts compiled successfully
- [ ] Tests passing locally
- [ ] Deployment script ready
- [ ] Testnet (Sepolia) selected

### Wallets
- [ ] Deployer wallet created
- [ ] Oracle wallet created (can be same as deployer initially)
- [ ] Both wallets funded with testnet ETH
  - [ ] Deployer: 0.1+ ETH
  - [ ] Oracle: 0.1+ ETH

### API Keys
- [ ] Infura/Alchemy account created
- [ ] Sepolia API key obtained
- [ ] API key tested (curl check)

## Deployment Phase

### 1. Deploy Contracts
```bash
cd contracts
npx hardhat run scripts/deploy.js --network sepolia
```

- [ ] MockUSDT deployed
- [ ] InsurancePool deployed
- [ ] Save addresses to `.env` and `config.yaml`
- [ ] Verify contracts on Etherscan (optional)

**Addresses to save:**
```
USDT Address: 0x_______________
Pool Address: 0x_______________
Oracle Address: 0x_______________
```

### 2. Fund Pool
```bash
npx hardhat console --network sepolia
```

```javascript
const usdt = await ethers.getContractAt("MockUSDT", "USDT_ADDR");
const pool = await ethers.getContractAt("InsurancePool", "POOL_ADDR");

// Mint 10,000 USDT
await usdt.mint(await signer.getAddress(), ethers.parseUnits("10000", 6));

// Approve pool
await usdt.approve(pool.target, ethers.parseUnits("5000", 6));

// Fund pool
await pool.fundPool(ethers.parseUnits("5000", 6));

// Verify
console.log("Pool balance:", await pool.getPoolBalance());
```

- [ ] USDT minted
- [ ] Pool funded with 5000+ USDT
- [ ] Balance verified

### 3. Configure Backend
```bash
cd backend
nano config.yaml
```

Update:
```yaml
rpc:
  url: "https://sepolia.infura.io/v3/YOUR_KEY"
  contract_address: "0xPoolAddress"
  private_key: "oracle_key_without_0x"
```

- [ ] RPC URL updated
- [ ] Contract address updated
- [ ] Private key updated
- [ ] Private key matches oracle address

### 4. Configure Frontend
```bash
cd frontend
nano .env
```

Update:
```bash
REACT_APP_INSURANCE_POOL_ADDRESS=0xPoolAddress
REACT_APP_USDT_ADDRESS=0xUSDTAddress
REACT_APP_RPC_URL=https://sepolia.infura.io/v3/YOUR_KEY
```

- [ ] Pool address updated
- [ ] USDT address updated
- [ ] RPC URL updated

### 5. Database Setup
```bash
cd backend
psql -U postgres -f db/schema.sql
```

- [ ] Database created
- [ ] Tables created
- [ ] Connection tested

## Testing Phase

### Backend Tests

#### Start Backend
```bash
cd backend
go run main.go
```

Expected output:
```
✅ Connected to InsurancePool contract at 0x...
✅ Oracle address: 0x...
📊 Starting detector...
```

- [ ] Backend starts without errors
- [ ] Connects to RPC
- [ ] Oracle address verified
- [ ] Detector running

#### Test Policy Creation
```javascript
// In Hardhat console
const usdt = await ethers.getContractAt("MockUSDT", "USDT_ADDR");
const pool = await ethers.getContractAt("InsurancePool", "POOL_ADDR");
const [signer] = await ethers.getSigners();

// Create test policy
await usdt.mint(signer.address, ethers.parseUnits("100", 6));
await usdt.approve(pool.target, ethers.parseUnits("10", 6));
const tx = await pool.buyInsurance();
await tx.wait();

// Verify
const count = await pool.getUserPoliciesCount(signer.address);
console.log("Policies:", count.toString());
```

- [ ] Policy created successfully
- [ ] Transaction mined
- [ ] Policy count incremented

#### Test Spike Detection
- [ ] Backend detects spike from CSV
- [ ] Queries active policies
- [ ] Prepares payout transaction

#### Test Payout Execution
Watch backend logs for:
```
🚨 SPIKE DETECTED!
Found 1 active policy/policies
🚀 Calling executePayout on-chain...
📤 Transaction sent: 0x...
⏳ Waiting for transaction to be mined...
✅ Transaction mined in block XXX
💰 Payout executed successfully: $100.00
```

- [ ] Payout transaction created
- [ ] Transaction submitted
- [ ] Transaction mined (within 1-2 minutes)
- [ ] USDT transferred to user
- [ ] Database updated

#### Verify on Etherscan
Visit: `https://sepolia.etherscan.io/address/POOL_ADDRESS`

- [ ] `executePayout` transaction visible
- [ ] Transaction status: Success
- [ ] USDT transfer event present
- [ ] Gas used: ~180,000-200,000

### Frontend Tests

#### Start Frontend
```bash
cd frontend
npm install
npm start
```

- [ ] Frontend starts on http://localhost:3000
- [ ] No console errors

#### Test Wallet Connection
- [ ] MetaMask connects
- [ ] Correct network (Sepolia)
- [ ] Account address displayed

#### Test Insurance Purchase
- [ ] Mint test USDT button works
- [ ] Balance updates
- [ ] Approve transaction successful
- [ ] Buy insurance transaction successful
- [ ] Policy appears in list

#### Test Notification System
- [ ] `PayoutNotification` component loaded
- [ ] Event listener active
- [ ] When payout occurs:
  - [ ] Notification appears
  - [ ] Shows correct amount
  - [ ] Shows policy ID
  - [ ] Shows transaction hash
  - [ ] Browser notification (if permitted)

## Integration Tests

### End-to-End Flow
1. [ ] User connects wallet
2. [ ] User mints test USDT
3. [ ] User buys insurance
4. [ ] Backend detects spike
5. [ ] Backend executes payout
6. [ ] User receives USDT
7. [ ] Frontend shows notification
8. [ ] Policy marked as claimed

### Verify All Components
- [ ] Smart contract working
- [ ] Backend detecting and paying
- [ ] Database in sync
- [ ] Frontend displaying correctly
- [ ] Events propagating
- [ ] Notifications appearing

## Post-Deployment

### Monitoring Setup
- [ ] Backend logs being written
- [ ] Error tracking configured
- [ ] Gas usage monitored
- [ ] Wallet balances monitored

### Documentation
- [ ] Deployment addresses recorded
- [ ] Oracle address documented
- [ ] API keys secured
- [ ] Setup guide updated

### Security
- [ ] Private keys not in git
- [ ] .env files in .gitignore
- [ ] Only testnet funds used
- [ ] Oracle wallet access restricted

## Troubleshooting

### If Backend Won't Start
1. [ ] Check RPC URL is valid
2. [ ] Check private key format (no 0x prefix)
3. [ ] Verify Go dependencies installed
4. [ ] Check database connection

### If Transactions Fail
1. [ ] Verify oracle has enough ETH for gas
2. [ ] Check pool has enough USDT
3. [ ] Verify oracle address matches contract
4. [ ] Check policy is active and not expired

### If Frontend Won't Connect
1. [ ] MetaMask on Sepolia network
2. [ ] Contract addresses correct in .env
3. [ ] RPC endpoint accessible
4. [ ] Contracts deployed to Sepolia

## Success Criteria

### ✅ Deployment Successful When:
- [ ] All contracts deployed
- [ ] Backend starts and connects
- [ ] Frontend loads and connects
- [ ] Can buy insurance
- [ ] Backend detects spikes
- [ ] Payouts execute on-chain
- [ ] Frontend receives events
- [ ] All data persisted correctly

## Next Steps After Successful Deployment

1. [ ] Run for 24-48 hours to monitor stability
2. [ ] Test with multiple policies
3. [ ] Test edge cases (expired policies, insufficient balance)
4. [ ] Document any issues encountered
5. [ ] Optimize gas usage if needed
6. [ ] Prepare for mainnet (if applicable)

## Emergency Procedures

### If Something Goes Wrong
1. Stop the backend service
2. Check error logs
3. Verify contract state on Etherscan
4. Check database state
5. Review recent transactions
6. Restore from backup if needed

### Contact Info
- [ ] Document support contacts
- [ ] Escalation procedures
- [ ] Emergency shutdown process

---

**Status**: [ ] Not Started | [ ] In Progress | [ ] Complete  
**Deployed By**: _____________  
**Date**: _____________  
**Network**: Sepolia Testnet  
**Notes**: _____________
