# Backend Setup Guide - Real Contract Integration

This guide explains how to set up the backend to make **real on-chain transactions** for insurance payouts.

## Prerequisites

1. **Go 1.19+** installed
2. **PostgreSQL** database running
3. **Deployed contracts** on a testnet (Sepolia recommended)
4. **Oracle wallet** with some testnet ETH for gas

## Step 1: Deploy Contracts

First, deploy the contracts to get the addresses:

```bash
cd ../contracts
npm install
npx hardhat run scripts/deploy.js --network sepolia
```

**Save the output:**
- MockUSDT address: `0x...`
- InsurancePool address: `0x...`
- Deployer/Oracle address: `0x...`

## Step 2: Configure Backend

Edit `backend/config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: spikeshield

rpc:
  url: https://sepolia.infura.io/v3/YOUR_INFURA_KEY
  contract_address: "0xYourInsurancePoolAddress"
  private_key: "your_private_key_without_0x_prefix"

detector:
  threshold_percent: 10.0
  window_minutes: 5

chainlink:
  btc_usd_feed: "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43"
  update_interval: 60

mode: replay  # or 'live'
```

### Important Configuration Notes

#### RPC URL
- Get free API key from [Infura](https://infura.io/) or [Alchemy](https://www.alchemy.com/)
- Sepolia testnet recommended for testing

#### Contract Address
- Use the **InsurancePool** contract address from deployment
- NOT the USDT address

#### Private Key
- **MUST be the oracle address** set in the contract
- Remove the `0x` prefix
- Example: If your key is `0xabc123...`, use `abc123...`
- ⚠️ **NEVER commit this to git!** Use environment variables in production

#### Verifying Oracle Address

The backend will verify on startup that your private key matches the contract's oracle:

```bash
✅ Connected to InsurancePool contract at 0x123...
✅ Oracle address: 0x742d35Cc...
```

If you see an error:
```
❌ private key does not match contract oracle
```

Then either:
1. Use the correct private key (the one that deployed the contract)
2. Or call `setOracle(newAddress)` on the contract to change it

## Step 3: Fund the Oracle Wallet

The oracle needs ETH for gas fees:

```bash
# Get testnet ETH from Sepolia faucet
# https://sepoliafaucet.com/
# https://www.alchemy.com/faucets/ethereum-sepolia

# Send 0.1 - 0.5 ETH to your oracle address
```

## Step 4: Fund the Insurance Pool

The pool needs USDT to pay out claims:

```bash
# Using Hardhat console
cd ../contracts
npx hardhat console --network sepolia

# Mint USDT to your address
const USDT = await ethers.getContractAt("MockUSDT", "0xYourUSDTAddress");
await USDT.mint("0xYourAddress", ethers.parseUnits("10000", 6));

# Approve and fund the pool
await USDT.approve("0xInsurancePoolAddress", ethers.parseUnits("5000", 6));
const Pool = await ethers.getContractAt("InsurancePool", "0xPoolAddress");
await Pool.fundPool(ethers.parseUnits("5000", 6));

# Verify pool balance
await Pool.getPoolBalance(); // Should show 5000 USDT (in 6 decimals)
```

## Step 5: Install Go Dependencies

```bash
cd backend
go mod download
go mod tidy
```

The key dependency is:
```go
github.com/ethereum/go-ethereum v1.13.0 // For contract interaction
```

## Step 6: Run the Backend

```bash
# Set up database
psql -U postgres -f db/schema.sql

# Run the service
go run main.go
```

### Expected Output

```
2024/01/15 10:00:00 🚀 Starting SpikeShield Backend Oracle
2024/01/15 10:00:01 ✅ Database connected
2024/01/15 10:00:02 ✅ Connected to InsurancePool contract at 0x123...
2024/01/15 10:00:02 ✅ Oracle address: 0x742d35Cc...
2024/01/15 10:00:03 📊 Starting detector in replay mode...
2024/01/15 10:00:05 🔍 Analyzing price data...
2024/01/15 10:00:10 🚨 SPIKE DETECTED! Price dropped 15.23%
2024/01/15 10:00:10 Found 2 active policy/policies
2024/01/15 10:00:11 🚀 Calling executePayout on-chain for user 0xabc..., policy 0
2024/01/15 10:00:11    Gas Price: 2500000000 wei
2024/01/15 10:00:11    Gas Limit: 300000
2024/01/15 10:00:12 📤 Transaction sent: 0x456def...
2024/01/15 10:00:12 ⏳ Waiting for transaction to be mined...
2024/01/15 10:00:27 ✅ Transaction mined in block 12345678
2024/01/15 10:00:27    Gas Used: 187432
2024/01/15 10:00:27 💰 Payout executed successfully for user 0xabc...: $100.00
```

## Step 7: Test the Integration

### 1. Create a Test Policy (from frontend or Hardhat)

```javascript
// Using Hardhat console
const USDT = await ethers.getContractAt("MockUSDT", "0xUSDTAddress");
const Pool = await ethers.getContractAt("InsurancePool", "0xPoolAddress");

// Mint and approve USDT
await USDT.mint(await signer.getAddress(), ethers.parseUnits("100", 6));
await USDT.approve(Pool.target, ethers.parseUnits("10", 6));

// Buy insurance
await Pool.buyInsurance();

// Verify policy
const count = await Pool.getUserPoliciesCount(await signer.getAddress());
console.log("Total policies:", count.toString());
```

### 2. Trigger Backend Detection

Backend will automatically detect spikes from the CSV data in replay mode, or from live Chainlink feeds in live mode.

### 3. Verify Payout

Check on Etherscan (Sepolia):
- Navigate to `https://sepolia.etherscan.io/address/0xYourPoolAddress`
- Look for `executePayout` transactions
- Verify USDT transfers to policy holders

## Troubleshooting

### Error: "Only oracle can execute payout"

**Cause:** Private key doesn't match the contract's oracle address

**Solution:**
```bash
# Check oracle address in contract
cast call 0xPoolAddress "oracle()(address)" --rpc-url $RPC_URL

# Check address from your private key
cast wallet address --private-key $PRIVATE_KEY

# If different, update oracle in contract (as owner)
cast send 0xPoolAddress "setOracle(address)" 0xNewOracleAddress --private-key $OWNER_KEY
```

### Error: "Insufficient pool balance"

**Cause:** Pool doesn't have enough USDT

**Solution:**
```bash
# Fund the pool (see Step 4)
```

### Error: "Transaction underpriced"

**Cause:** Gas price too low

**Solution:** Backend auto-suggests gas price, but you can increase the multiplier:
```go
// In payout.go
gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120)) // 20% higher
gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))
```

### Error: "Failed to wait for transaction"

**Cause:** Network congestion or transaction reverted

**Solution:**
- Check transaction on Etherscan
- Increase gas limit if out-of-gas
- Check if policy is valid (not expired, not already claimed)

## Production Considerations

### Security

1. **Never commit private keys** to version control
2. Use environment variables:
   ```bash
   export ORACLE_PRIVATE_KEY="your_key"
   export RPC_URL="your_rpc_url"
   ```
3. Use secrets management (AWS Secrets Manager, HashiCorp Vault)
4. Restrict oracle wallet permissions (separate from owner wallet)

### Monitoring

1. Monitor oracle wallet ETH balance
2. Monitor pool USDT balance
3. Set up alerts for failed transactions
4. Log all payout attempts

### Gas Optimization

1. Batch payouts if possible
2. Monitor gas prices and delay non-urgent transactions
3. Use EIP-1559 gas pricing for better estimates

### Scaling

1. Use archive node for historical queries
2. Implement retry logic with exponential backoff
3. Queue system for high transaction volume
4. Multiple oracle wallets for load distribution

## Testing Checklist

- [ ] Contract deployed successfully
- [ ] Oracle address matches private key
- [ ] Pool funded with sufficient USDT
- [ ] Oracle wallet has ETH for gas
- [ ] Backend connects to RPC
- [ ] Test policy created
- [ ] Spike detection triggers
- [ ] Payout transaction sent
- [ ] Transaction mined successfully
- [ ] User receives USDT
- [ ] Database updated correctly
- [ ] Frontend shows payout event

## Next Steps

1. Set up monitoring and alerts
2. Configure proper logging
3. Implement error recovery
4. Add payout queue system
5. Set up continuous deployment
6. Add more comprehensive tests

## Support

For issues, check:
- Backend logs: `logs/spikeshield.log`
- Database: Check `payouts` table
- Etherscan: Verify transactions on-chain
