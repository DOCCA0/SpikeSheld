const { ethers, upgrades } = require("hardhat");

async function main() {
  const [deployer] = await ethers.getSigners();
  
  console.log("Deploying contracts with the account:", deployer.address);
  console.log("Account balance:", (await ethers.provider.getBalance(deployer.address)).toString());

  // Deploy MockUSDT (not upgradeable)
  const MockUSDT = await ethers.getContractFactory("MockUSDT");
  const mockUSDT = await MockUSDT.deploy();
  await mockUSDT.waitForDeployment();
  const usdtAddress = await mockUSDT.getAddress();
  console.log("MockUSDT deployed to:", usdtAddress);
  
  // Deploy InsurancePool as upgradeable proxy
  const InsurancePool = await ethers.getContractFactory("InsurancePool");
  console.log("Deploying InsurancePool proxy...");
  const insurancePool = await upgrades.deployProxy(InsurancePool, [usdtAddress], {
    initializer: 'initialize',
    kind: 'transparent'
  });
  await insurancePool.waitForDeployment();
  const proxyAddress = await insurancePool.getAddress();
  console.log("InsurancePool proxy deployed to:", proxyAddress);

  // Get implementation address
  const implementationAddress = await upgrades.erc1967.getImplementationAddress(proxyAddress);
  console.log("InsurancePool implementation deployed to:", implementationAddress);

  // Get admin address
  const adminAddress = await upgrades.erc1967.getAdminAddress(proxyAddress);
  console.log("ProxyAdmin deployed to:", adminAddress);

  // Fund the pool (for demo)
  const fundAmount = ethers.parseUnits("10000", 6); // 10,000 USDT
  await mockUSDT.approve(proxyAddress, fundAmount);
  await insurancePool.fundPool(fundAmount, {
    gasLimit: 300000,
    gasPrice: ethers.parseUnits("50", "gwei")
  });
  console.log("MockUSDT:", usdtAddress);
  console.log("InsurancePool Proxy:", proxyAddress);
  console.log("InsurancePool Implementation:", implementationAddress);
  console.log("ProxyAdmin:", adminAddress);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
