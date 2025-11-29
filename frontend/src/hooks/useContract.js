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
  const INSURANCE_POOL_ADDRESS = process.env.REACT_APP_INSURANCE_POOL_ADDRESS ;
  const USDT_ADDRESS = process.env.REACT_APP_USDT_ADDRESS;

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
      const checksummedAccount = ethers.getAddress(accounts[0]);
      setProvider(provider);
      setSigner(signer);
      setAccount(checksummedAccount);

      // Initialize contracts
      const insuranceAddress = ethers.getAddress(INSURANCE_POOL_ADDRESS);
      const usdtAddress = ethers.getAddress(USDT_ADDRESS);
      const insurance = new ethers.Contract(insuranceAddress, INSURANCE_POOL_ABI, signer);
      const usdt = new ethers.Contract(usdtAddress, USDT_ABI, signer);
      
      setInsuranceContract(insurance);
      setUsdtContract(usdt);

      console.log("Wallet connected:", checksummedAccount);
      // Notify backend to upsert this user's balances and policies
      try {
        const apiBase = process.env.REACT_APP_API_URL || '';
        const url = apiBase ? `${apiBase.replace(/\/$/, '')}/api/wallet/link` : '/api/wallet/link';
        console.log('address', checksummedAccount);
        fetch(url, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ address: checksummedAccount, token: usdtAddress }),
        }).then((res) => {
          if (!res.ok) console.error('Wallet link API responded with', res.status);
          else console.log('Backend wallet-link upsert triggered');
        }).catch((err) => {
          console.error('Failed to notify backend on wallet link', err);
        });
      } catch (err) {
        console.error('Wallet link notify error:', err);
      }
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
      const approveTx = await usdtContract.approve(insuranceAddress, premium);
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
  // Note: Reading user policies, balances and event subscriptions should be done server-side.

  // Listen to account changes
  useEffect(() => {
    if (window.ethereum) {
      window.ethereum.on('accountsChanged', (accounts) => {
        if (accounts.length === 0) {
          disconnectWallet();
        } else {
          setAccount(ethers.getAddress(accounts[0]));
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
    // getBalance and listenForPayouts removed: use backend APIs for reads/push updates
    isConnected: !!account
  };
};
