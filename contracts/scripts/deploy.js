const hre = require("hardhat");

async function main() {
  console.log("🚀 Deploying SpikeShield contracts...");

  // Deploy MockUSDT
  console.log("\n1️⃣ Deploying MockUSDT...");
  const MockUSDT = await hre.ethers.getContractFactory("MockUSDT");
  const usdt = await MockUSDT.deploy();
  await usdt.waitForDeployment();
  const usdtAddress = await usdt.getAddress();
  console.log("✅ MockUSDT deployed to:", usdtAddress);

  // Deploy InsurancePool
  console.log("\n2️⃣ Deploying InsurancePool...");
  const InsurancePool = await hre.ethers.getContractFactory("InsurancePool");
  const insurance = await InsurancePool.deploy(usdtAddress);
  await insurance.waitForDeployment();
  const insuranceAddress = await insurance.getAddress();
  console.log("✅ InsurancePool deployed to:", insuranceAddress);

  // Fund the pool with initial USDT
  console.log("\n3️⃣ Funding InsurancePool...");
  const fundAmount = hre.ethers.parseUnits("10000", 6); // 10,000 USDT
  await usdt.approve(insuranceAddress, fundAmount);
  await insurance.fundPool(fundAmount);
  console.log("✅ Pool funded with 10,000 USDT");

  // Display summary
  console.log("\n" + "=".repeat(60));
  console.log("📋 DEPLOYMENT SUMMARY");
  console.log("=".repeat(60));
  console.log("MockUSDT Address:        ", usdtAddress);
  console.log("InsurancePool Address:   ", insuranceAddress);
  console.log("Network:                 ", hre.network.name);
  console.log("=".repeat(60));
  
  console.log("\n📝 Update your .env file with these addresses:");
  console.log(`INSURANCE_POOL_ADDRESS=${insuranceAddress}`);
  console.log(`USDT_ADDRESS=${usdtAddress}`);
  console.log(`REACT_APP_INSURANCE_POOL_ADDRESS=${insuranceAddress}`);
  console.log(`REACT_APP_USDT_ADDRESS=${usdtAddress}`);

  console.log("\n✅ Deployment complete!");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
