const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("InsurancePool", function () {
  let insurancePool;
  let mockUSDT;
  let owner;
  let user1;
  let user2;

  const PREMIUM_AMOUNT = ethers.parseUnits("10", 6); // 10 USDT
  const COVERAGE_AMOUNT = ethers.parseUnits("100", 6); // 100 USDT
  const INITIAL_MINT = ethers.parseUnits("1000", 6); // 1000 USDT per user

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    // Deploy MockUSDT
    const MockUSDT = await ethers.getContractFactory("MockUSDT");
    mockUSDT = await MockUSDT.deploy();
    await mockUSDT.waitForDeployment();

    // Deploy InsurancePool
    const InsurancePool = await ethers.getContractFactory("InsurancePool");
    insurancePool = await InsurancePool.deploy(await mockUSDT.getAddress());
    await insurancePool.waitForDeployment();

    // Mint USDT to users for testing
    await mockUSDT.mint(user1.address, INITIAL_MINT);
    await mockUSDT.mint(user2.address, INITIAL_MINT);

    // Fund the insurance pool
    const poolFundAmount = ethers.parseUnits("10000", 6);
    await mockUSDT.approve(await insurancePool.getAddress(), poolFundAmount);
    await insurancePool.fundPool(poolFundAmount);
  });

  describe("Deployment", function () {
    it("Should set correct initial values", async function () {
      expect(await insurancePool.usdt()).to.equal(await mockUSDT.getAddress());
      expect(await insurancePool.owner()).to.equal(owner.address);
      expect(await insurancePool.oracle()).to.equal(owner.address);
    });
  });

  describe("Buy Insurance", function () {
    it("Should allow user to purchase insurance", async function () {
      await mockUSDT.connect(user1).approve(await insurancePool.getAddress(), PREMIUM_AMOUNT);
      await insurancePool.connect(user1).buyInsurance();

      expect(await insurancePool.getUserPoliciesCount(user1.address)).to.equal(1);
      expect(await insurancePool.hasActivePolicy(user1.address)).to.be.true;
    });

    it("Should transfer premium to pool", async function () {
      const initialBalance = await insurancePool.getPoolBalance();

      await mockUSDT.connect(user1).approve(await insurancePool.getAddress(), PREMIUM_AMOUNT);
      await insurancePool.connect(user1).buyInsurance();

      expect(await insurancePool.getPoolBalance()).to.equal(initialBalance + PREMIUM_AMOUNT);
    });

    it("Should fail without USDT approval", async function () {
      await expect(
        insurancePool.connect(user1).buyInsurance()
      ).to.be.reverted;
    });
  });

  describe("Execute Payout", function () {
    beforeEach(async function () {
      await mockUSDT.connect(user1).approve(await insurancePool.getAddress(), PREMIUM_AMOUNT);
      await insurancePool.connect(user1).buyInsurance();
    });

    it("Should execute payout for valid policy", async function () {
      const initialBalance = await mockUSDT.balanceOf(user1.address);

      await insurancePool.connect(owner).executePayout(user1.address, 0);

      expect(await mockUSDT.balanceOf(user1.address)).to.equal(initialBalance + COVERAGE_AMOUNT);
    });

    it("Should mark policy as claimed", async function () {
      await insurancePool.connect(owner).executePayout(user1.address, 0);

      const policy = await insurancePool.getPolicy(user1.address, 0);
      expect(policy.claimed).to.be.true;
      expect(policy.active).to.be.false;
    });

    it("Should fail when non-oracle tries payout", async function () {
      await expect(
        insurancePool.connect(user2).executePayout(user1.address, 0)
      ).to.be.revertedWith("Only oracle can execute payout");
    });

    it("Should fail when policy expired", async function () {
      await time.increase(24 * 60 * 60 + 1); // Fast forward 24 hours + 1 second

      await expect(
        insurancePool.connect(owner).executePayout(user1.address, 0)
      ).to.be.revertedWith("Policy expired");
    });

    it("Should fail for already claimed policy", async function () {
      await insurancePool.connect(owner).executePayout(user1.address, 0);

      await expect(
        insurancePool.connect(owner).executePayout(user1.address, 0)
      ).to.be.revertedWith("Policy not active");
    });
  });

  describe("Admin Functions", function () {
    it("Should allow owner to update oracle", async function () {
      await insurancePool.connect(owner).setOracle(user2.address);
      expect(await insurancePool.oracle()).to.equal(user2.address);
    });

    it("Should allow owner to update policy parameters", async function () {
      const newPremium = ethers.parseUnits("20", 6);
      const newCoverage = ethers.parseUnits("200", 6);

      await insurancePool.connect(owner).setPolicyParams(newPremium, newCoverage, 48 * 60 * 60);

      expect(await insurancePool.premiumAmount()).to.equal(newPremium);
      expect(await insurancePool.coverageAmount()).to.equal(newCoverage);
    });
  });
});
