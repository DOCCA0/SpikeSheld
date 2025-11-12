import { useState, useEffect } from 'react';
import { ethers } from 'ethers';

// Contract ABI (simplified for demo)
const INSURANCE_POOL_ABI = [
  "function buyInsurance() external",
  "function getUserPoliciesCount(address user) external view returns (uint256)",
  "function getPolicy(address user, uint256 policyId) external view returns (tuple(address user, uint256 premium, uint256 coverageAmount, uint256 purchaseTime, uint256 expiryTime, bool active, bool claimed))",
  "function hasActivePolicy(address user) external view returns (bool)",
  "function premiumAmount() external view returns (uint256)",
  "function coverageAmount() external view returns (uint256)",
  "function getPoolBalance() external view returns (uint256)",
  "event PolicyPurchased(address indexed user, uint256 policyId, uint256 premium, uint256 coverage, uint256 expiryTime)",
  "event PayoutExecuted(address indexed user, uint256 policyId, uint256 amount)"
];

const USDT_ABI = [
  "function balanceOf(address owner) external view returns (uint256)",
  "function approve(address spender, uint256 amount) external returns (bool)",
  "function mint(address to, uint256 amount) external",
  "function decimals() external view returns (uint8)"
];

export const useContract = () => {
  const [account, setAccount] = useState(null);
  const [provider, setProvider] = useState(null);
  const [signer, setSigner] = useState(null);
  const [insuranceContract, setInsuranceContract] = useState(null);
  const [usdtContract, setUsdtContract] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Replace with your deployed contract addresses
  const INSURANCE_POOL_ADDRESS = process.env.REACT_APP_INSURANCE_POOL_ADDRESS || "0x0000000000000000000000000000000000000000";
  const USDT_ADDRESS = process.env.REACT_APP_USDT_ADDRESS || "0x0000000000000000000000000000000000000000";

  // Connect wallet
  const connectWallet = async () => {
    try {
      setLoading(true);
      setError(null);

      if (!window.ethereum) {
        throw new Error("Please install MetaMask!");
      }

      const accounts = await window.ethereum.request({
        method: 'eth_requestAccounts'
      });

      const provider = new ethers.BrowserProvider(window.ethereum);
      const signer = await provider.getSigner();
      
      setProvider(provider);
      setSigner(signer);
      setAccount(accounts[0]);

      // Initialize contracts
      const insurance = new ethers.Contract(INSURANCE_POOL_ADDRESS, INSURANCE_POOL_ABI, signer);
      const usdt = new ethers.Contract(USDT_ADDRESS, USDT_ABI, signer);
      
      setInsuranceContract(insurance);
      setUsdtContract(usdt);

      console.log("Wallet connected:", accounts[0]);
    } catch (err) {
      setError(err.message);
      console.error("Connection error:", err);
    } finally {
      setLoading(false);
    }
  };

  // Disconnect wallet
  const disconnectWallet = () => {
    setAccount(null);
    setProvider(null);
    setSigner(null);
    setInsuranceContract(null);
    setUsdtContract(null);
  };

  // Buy insurance
  const buyInsurance = async () => {
    try {
      setLoading(true);
      setError(null);

      if (!insuranceContract || !usdtContract) {
        throw new Error("Contracts not initialized");
      }

      // Get premium amount
      const premium = await insuranceContract.premiumAmount();
      console.log("Premium:", ethers.formatUnits(premium, 6), "USDT");

      // Approve USDT spending
      console.log("Approving USDT...");
      const approveTx = await usdtContract.approve(INSURANCE_POOL_ADDRESS, premium);
      await approveTx.wait();
      console.log("USDT approved");

      // Buy insurance
      console.log("Buying insurance...");
      const buyTx = await insuranceContract.buyInsurance();
      const receipt = await buyTx.wait();
      console.log("Insurance purchased:", receipt.hash);

      return receipt.hash;
    } catch (err) {
      setError(err.message);
      console.error("Buy insurance error:", err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // Mint test USDT
  const mintTestUSDT = async (amount) => {
    try {
      setLoading(true);
      setError(null);

      if (!usdtContract || !account) {
        throw new Error("Wallet not connected");
      }

      const tx = await usdtContract.mint(account, ethers.parseUnits(amount.toString(), 6));
      await tx.wait();
      console.log("Minted", amount, "test USDT");
    } catch (err) {
      setError(err.message);
      console.error("Mint error:", err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // Get user balance
  const getBalance = async () => {
    try {
      if (!usdtContract || !account) return 0;
      const balance = await usdtContract.balanceOf(account);
      return ethers.formatUnits(balance, 6);
    } catch (err) {
      console.error("Balance error:", err);
      return 0;
    }
  };

  // Get user policies
  const getUserPolicies = async () => {
    try {
      if (!insuranceContract || !account) return [];
      
      const count = await insuranceContract.getUserPoliciesCount(account);
      const policies = [];

      for (let i = 0; i < count; i++) {
        const policy = await insuranceContract.getPolicy(account, i);
        policies.push({
          id: i,
          premium: ethers.formatUnits(policy.premium, 6),
          coverage: ethers.formatUnits(policy.coverageAmount, 6),
          purchaseTime: new Date(Number(policy.purchaseTime) * 1000),
          expiryTime: new Date(Number(policy.expiryTime) * 1000),
          active: policy.active,
          claimed: policy.claimed
        });
      }

      return policies;
    } catch (err) {
      console.error("Get policies error:", err);
      return [];
    }
  };

  // Check if user has active policy
  const hasActivePolicy = async () => {
    try {
      if (!insuranceContract || !account) return false;
      return await insuranceContract.hasActivePolicy(account);
    } catch (err) {
      console.error("Check active policy error:", err);
      return false;
    }
  };

  // Listen for payout events
  const listenForPayouts = (callback) => {
    if (!insuranceContract || !account) return;

    const filter = insuranceContract.filters.PayoutExecuted(account);
    
    insuranceContract.on(filter, (user, policyId, amount, event) => {
      console.log("🎉 Payout received!", {
        user,
        policyId: policyId.toString(),
        amount: ethers.formatUnits(amount, 6),
        blockNumber: event.log.blockNumber
      });
      
      if (callback) {
        callback({
          user,
          policyId: Number(policyId),
          amount: ethers.formatUnits(amount, 6),
          blockNumber: event.log.blockNumber
        });
      }
    });

    // Return cleanup function
    return () => {
      insuranceContract.off(filter);
    };
  };

  // Listen to account changes
  useEffect(() => {
    if (window.ethereum) {
      window.ethereum.on('accountsChanged', (accounts) => {
        if (accounts.length === 0) {
          disconnectWallet();
        } else {
          setAccount(accounts[0]);
        }
      });

      window.ethereum.on('chainChanged', () => {
        window.location.reload();
      });
    }
  }, []);

  return {
    account,
    loading,
    error,
    connectWallet,
    disconnectWallet,
    buyInsurance,
    mintTestUSDT,
    getBalance,
    getUserPolicies,
    hasActivePolicy,
    listenForPayouts,
    isConnected: !!account
  };
};
