# Contract Usage Guide

## Calling On-Chain Methods

### Using ethers.js
```javascript
const { ethers } = require("hardhat");

// Get contract instance
const InsurancePool = await ethers.getContractFactory("InsurancePool");
const insurance = await InsurancePool.attach("0xb2C5582d7a782A6bA79b2751E9cE3c0dA3F4BeA5");

// Get usdt instance
const MUSDT = await ethers.getContractFactory("MockUSDT");
const musdt = await MUSDT.attach("0xa56E945954B795c00b49D514cB53ecdC3932bC55");

const aaa = await insurance.xxx
const bbb = await musdt.yyy

// With signer
const [owner, user1] = await ethers.getSigners();
const poolAsUser = pool.connect(user1);
await poolAsUser.deposit(amount);
```

## Common Hardhat Commands

```bash
# Install dependencies
npm install

# Compile contracts
npx hardhat compile

# Run tests
npx hardhat test
npx hardhat test test/InsurancePool.test.js  # Single file


# Deploy
npx hardhat run scripts/deploy.js --network localhost
npx hardhat run scripts/deploy.js --network sepolia
npx hardhat run scripts/deploy.js --network mainnet

# Start local node
npx hardhat node

# Console
npx hardhat console --network localhost
npx hardhat console --network sepolia

# Clean
npx hardhat clean

# Gas report
REPORT_GAS=true npx hardhat test

# Coverage
npx hardhat coverage
```
