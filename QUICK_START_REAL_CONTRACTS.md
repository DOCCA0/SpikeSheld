# Quick Start - Real Contract Integration

## TL;DR

Get the backend making real blockchain transactions in 5 minutes.

## Prerequisites
- Contracts deployed to Sepolia
- Oracle wallet with 0.1+ ETH
- Pool funded with 1000+ USDT

## Step 1: Update Config (30 seconds)

```bash
cd backend
nano config.yaml
```

Update these 3 lines:
```yaml
rpc:
  url: "https://sepolia.infura.io/v3/YOUR_INFURA_KEY"
  contract_address: "0xYourInsurancePoolAddress"
  private_key: "your_private_key_without_0x_prefix"
```

## Step 2: Run Backend (10 seconds)

```bash
go run main.go
```

Look for:
```
✅ Connected to InsurancePool contract at 0x...
✅ Oracle address: 0x...
```

## Step 3: Test (1 minute)

### Create a policy (from Hardhat console):
```javascript
npx hardhat console --network sepolia

const usdt = await ethers.getContractAt("MockUSDT", "0xUSDTAddr");
const pool = await ethers.getContractAt("InsurancePool", "0xPoolAddr");

await usdt.mint(await signer.getAddress(), ethers.parseUnits("100", 6));
await usdt.approve(pool.target, ethers.parseUnits("10", 6));
await pool.buyInsurance();
```

### Watch backend detect and payout:
```
🚨 SPIKE DETECTED!
🚀 Calling executePayout on-chain...
📤 Transaction sent: 0x123...
✅ Transaction mined in block 456
💰 Payout executed successfully: $100.00
```

## Done! 🎉

Your backend is now making real blockchain transactions.

## Troubleshooting

### "Only oracle can execute payout"
→ Private key doesn't match contract oracle

**Fix:**
```bash
# Check what oracle address the contract expects
cast call 0xPoolAddr "oracle()(address)" --rpc-url $RPC_URL
# Make sure your private key controls that address
```

### "Insufficient pool balance"
→ Pool needs more USDT

**Fix:**
```javascript
const usdt = await ethers.getContractAt("MockUSDT", "0xUSDTAddr");
const pool = await ethers.getContractAt("InsurancePool", "0xPoolAddr");
await usdt.mint(await signer.getAddress(), ethers.parseUnits("5000", 6));
await usdt.approve(pool.target, ethers.parseUnits("5000", 6));
await pool.fundPool(ethers.parseUnits("5000", 6));
```

### "Failed to connect to RPC"
→ Check your Infura/Alchemy API key

## Full Documentation

See `SETUP_REAL_CONTRACTS.md` for complete details.

## Frontend Integration

Update `.env`:
```bash
REACT_APP_INSURANCE_POOL_ADDRESS=0xYourPoolAddress
REACT_APP_USDT_ADDRESS=0xYourUSDTAddress
```

Add to `App.js`:
```javascript
import PayoutNotification from './components/PayoutNotification';

function App() {
  return (
    <>
      <PayoutNotification />
      {/* rest of your app */}
    </>
  );
}
```

Users will see real-time notifications when they receive payouts! 🎉
