import React, { useState, useEffect } from 'react';
import './App.css';
import { useContract } from './hooks/useContract';
import { apiService } from './services/api';
import PriceChart from './components/PriceChart';

function App() {
  const {
    account,
    loading,
    error,
    connectWallet,
    disconnectWallet,
    buyInsurance,
    mintTestUSDT,
    isConnected
  } = useContract();

  const [balance, setBalance] = useState('0');
  const [policies, setPolicies] = useState([]);
  const [successMsg, setSuccessMsg] = useState('');
  const [refreshing, setRefreshing] = useState(false);
  
  // Backend API data
  const [spikes, setSpikes] = useState([]);
  const [prices, setPrices] = useState([]);
  const [payouts, setPayouts] = useState([]);
  const [stats, setStats] = useState(null);
  const [apiStatus, setApiStatus] = useState('checking');

  // Check API health
  useEffect(() => {
    const checkAPI = async () => {
      try {
        await apiService.healthCheck();
        setApiStatus('online');
        loadBackendData();
      } catch (err) {
        setApiStatus('offline');
        console.error('API health check failed:', err);
      }
    };
    checkAPI();
    const interval = setInterval(checkAPI, 5000); // Check every 5s
    return () => clearInterval(interval);
  }, []);

  // Load backend data
  const loadBackendData = async () => {
    try {
      const [spikesData, pricesData, payoutsData, statsData] = await Promise.all([
        apiService.getSpikes(10),
        apiService.getPrices('BTCUSDT', 100),
        apiService.getPayouts(10),
        apiService.getStats()
      ]);
      setSpikes(spikesData.spikes || []);
      setPrices(pricesData.prices || []);
      setPayouts(payoutsData.payouts || []);
      setStats(statsData.stats || null);
    } catch (err) {
      console.error('Failed to load backend data:', err);
    }
  };

  // Refresh data
  const refreshData = async () => {
    if (!isConnected) return;
    
    setRefreshing(true);
    try {
      // Read balance and policies from backend API (no direct chain reads in frontend)
      const [balResp, polsResp] = await Promise.all([
        apiService.getUserBalance(account),
        apiService.getUserPolicies(account)
      ]);

      setBalance(balResp?.balance ?? '0');

      const rawPolicies = polsResp?.policies || [];
      const mapped = rawPolicies.map((p, idx) => ({
        id: p.ID ?? p.id ?? idx,
        premium: (p.Premium ?? p.premium ?? 0).toString(),
        coverage: p.CoverageAmount ?? p.coverage ?? p.coverage_amount ?? 0,
        purchaseTime: new Date(p.PurchaseTime ?? p.purchase_time ?? p.purchaseTime),
        expiryTime: new Date(p.ExpiryTime ?? p.expiry_time ?? p.expiryTime),
        active: (p.Status ?? p.status) === 'active',
        claimed: (p.Status ?? p.status) === 'claimed'
      }));

      setPolicies(mapped);

      const activeFlag = mapped.some(pol => pol.active && pol.expiryTime > new Date());
      setHasActive(activeFlag);
      
      // Also refresh backend data
      if (apiStatus === 'online') {
        await loadBackendData();
      }
    } catch (err) {
      console.error("Refresh error:", err);
    } finally {
      setRefreshing(false);
    }
  };

  // Auto refresh when connected + link wallet
  useEffect(() => {
    if (isConnected && account) {
      const initWalletData = async () => {
        try {
          // Link wallet to trigger backend sync
          await apiService.linkWallet(account, process.env.REACT_APP_USDT_ADDRESS);
          console.log('Wallet linked successfully');
          // Then refresh data
          await refreshData();
        } catch (err) {
          console.error('Wallet initialization failed:', err);
        }
      };
      initWalletData();
    }
  }, [isConnected, account]);

  // Handle buy insurance
  const handleBuyInsurance = async () => {
    try {
      setSuccessMsg('');
      const txHash = await buyInsurance();
      setSuccessMsg(`Insurance purchased successfully! Tx: ${txHash ? `${txHash.slice(0, 10)}...` : 'pending'}`);
      setTimeout(() => refreshData(), 2000);
    } catch (err) {
      console.error("Purchase failed:", err);
    }
  };

  // Handle mint USDT
  const handleMintUSDT = async () => {
    try {
      setSuccessMsg('');
      await mintTestUSDT(100);
      setSuccessMsg('Minted 100 test USDT successfully!');
      setTimeout(() => refreshData(), 2000);
    } catch (err) {
      console.error("Mint failed:", err);
    }
  };

  // Format date (UTC: YYYY-MM-DD HH:mm:ss)
  const formatDate = (date) => {
    if (!date) return 'N/A';
    const d = new Date(date);
    const year = d.getUTCFullYear();
    const month = String(d.getUTCMonth() + 1).padStart(2, '0');
    const day = String(d.getUTCDate()).padStart(2, '0');
    const hour = String(d.getUTCHours()).padStart(2, '0');
    const minute = String(d.getUTCMinutes()).padStart(2, '0');
    const second = String(d.getUTCSeconds()).padStart(2, '0');
    return `${year}-${month}-${day} ${hour}:${minute}:${second}`;
  };

  // Get policy status
  const getPolicyStatus = (policy) => {
    if (policy.claimed) return 'claimed';
    if (policy.expiryTime < new Date()) return 'expired';
    if (policy.active) return 'active';
    return 'inactive';
  };

  return (
    <div className="App">
      <div className="container">
        {/* Header */}
        <div className="header">
          <h1>🛡️ SpikeShield</h1>
          <p>Decentralized Spike Insurance Protocol</p>
        </div>

        {/* Wallet Section */}
        <div className="wallet-section">
          {!isConnected ? (
            <div>
              <h2>Connect Your Wallet</h2>
              <p>Connect your wallet to purchase spike insurance</p>
              <button 
                className="button" 
                onClick={connectWallet}
                disabled={loading}
              >
                {loading ? <span className="loading"></span> : 'Connect Wallet'}
              </button>
            </div>
          ) : (
            <div>
              <h2>Wallet Connected</h2>
              <div className="account-info">
                <div className="account-address">
                  {account ? `${account.slice(0, 6)}...${account.slice(-4)}` : 'Loading...'}
                </div>
                <div className="balance-display">
                  {parseFloat(balance || '0').toFixed(2)} USDT
                </div>
              </div>
              <button 
                className="button button-secondary" 
                onClick={handleMintUSDT}
                disabled={loading}
              >
                {loading ? <span className="loading"></span> : 'Mint 100 Test USDT'}
              </button>
              <button 
                className="button button-danger" 
                onClick={disconnectWallet}
              >
                Disconnect
              </button>
              <button 
                className="button" 
                onClick={refreshData}
                disabled={refreshing}
              >
                {refreshing ? <span className="loading"></span> : 'Refresh Data'}
              </button>
            </div>
          )}
        </div>

        {/* Error/Success Messages */}
        {error && (
          <div className="error-message">
            ❌ {error}
          </div>
        )}
        {successMsg && (
          <div className="success-message">
            ✅ {successMsg}
          </div>
        )}

        {/* Insurance Section */}
        {isConnected && (
          <div className="insurance-section">
            <h2>📋 Insurance Coverage</h2>
            <div className="insurance-info">
              <div className="info-card">
                <h3>Premium</h3>
                <div className="value">10 USDT</div>
              </div>
              <div className="info-card">
                <h3>Coverage Amount</h3>
                <div className="value">100 USDT</div>
              </div>
              <div className="info-card">
                <h3>Duration</h3>
                <div className="value">24 Hours</div>
              </div>
              <div className="info-card">
                <h3>
                  Body Ratio 
                  <span className="tooltip">
                    ❓
                    <span className="tooltiptext">
                      Maximum ratio of candle body size to total range.<br/>
                      Formula: |open-close| / (high-low)<br/>
                      Smaller value = longer wick
                    </span>
                  </span>
                </h3>
                <div className="value">≤ 30%</div>
              </div>
              <div className="info-card">
                <h3>
                  Range Ratio 
                  <span className="tooltip">
                    ❓
                    <span className="tooltiptext">
                      Minimum ratio of price range to close price.<br/>
                      Formula: (high-low) / close<br/>
                      Larger value = bigger price swing
                    </span>
                  </span>
                </h3>
                <div className="value">≥ 10%</div>
              </div>
            </div>

            <div style={{ margin: '30px 0' }}>
              <h3>How It Works</h3>
              <ul className="feature-list">
                <li>💰 Pay 10 USDT premium to get 100 USDT coverage</li>
                <li>⏱️ Protection valid for 24 hours</li>
                <li>📊 <strong>Body Ratio</strong>: abs(open-close)/(high-low) ≤ 0.3 (small body = long wick)</li>
                <li>📈 <strong>Range Ratio</strong>: (high-low)/close ≥ 0.1 (large price range)</li>
                <li>🚨 When both conditions met, automatic payout triggered</li>
                <li>⚡ Backend monitors price in real-time or replay mode</li>
                <li>🤖 Smart contract executes payout automatically</li>
              </ul>
            </div>

            <button 
              className="button" 
              onClick={handleBuyInsurance}
              disabled={loading || parseFloat(balance) < 10}
              style={{ fontSize: '1.3em', padding: '20px 40px' }}
            >
              {loading ? (
                <span className="loading"></span>
              ) : (
                '🛡️ Buy Insurance (10 USDT)'
              )}
            </button>
            
            {parseFloat(balance) < 10 && (
              <p style={{ color: '#dc3545', marginTop: '10px' }}>
                ⚠️ Insufficient balance. Mint test USDT first.
              </p>
            )}
          </div>
        )}

        {/* Policies Section */}
        {isConnected && policies.length > 0 && (
          <div className="policies-section">
            <h2>📜 Your Insurance Policies</h2>
            <div className="policies-list">
              {policies.map((policy) => {
                const status = getPolicyStatus(policy);
                return (
                  <div key={policy.id} className={`policy-card ${status}`}>
                    <h4>
                      Policy #{policy.id}
                      <span 
                        className={`status-badge status-${status}`}
                        style={{ marginLeft: '10px' }}
                      >
                        {status}
                      </span>
                    </h4>
                    <div className="policy-detail">
                      <span className="policy-label">Premium:</span>
                      <span className="policy-value">{policy.premium} USDT</span>
                    </div>
                    <div className="policy-detail">
                      <span className="policy-label">Coverage:</span>
                      <span className="policy-value">{policy.coverage} USDT</span>
                    </div>
                    <div className="policy-detail">
                      <span className="policy-label">Purchased:</span>
                      <span className="policy-value">{formatDate(policy.purchaseTime)}</span>
                    </div>
                    <div className="policy-detail">
                      <span className="policy-label">Expires:</span>
                      <span className="policy-value">{formatDate(policy.expiryTime)}</span>
                    </div>
                    {policy.claimed && (
                      <div style={{ 
                        marginTop: '15px', 
                        padding: '10px', 
                        background: '#667eea', 
                        color: 'white', 
                        borderRadius: '5px',
                        textAlign: 'center',
                        fontWeight: 'bold'
                      }}>
                        🎉 Payout Executed!
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {/* Footer Info */}
        {!isConnected && (
          <div style={{ 
            background: 'rgba(255, 255, 255, 0.95)', 
            borderRadius: '20px', 
            padding: '30px',
            marginTop: '20px',
            textAlign: 'left'
          }}>
            <h2>🚀 About SpikeShield</h2>
            <p>
              SpikeShield is a decentralized insurance protocol that protects you against 
              sudden price drops (wicks) in cryptocurrency markets.
            </p>
            <h3>Features:</h3>
            <ul className="feature-list">
              <li>🔄 <strong>Replay Mode:</strong> Test with historical data</li>
              <li>⚡ <strong>Live Mode:</strong> Real-time monitoring using Chainlink Oracle</li>
              <li>🤖 <strong>Automatic Payout:</strong> Smart contract executes payout when wick detected</li>
              <li>🧪 <strong>Testnet Ready:</strong> Deploy on Sepolia or BSC Testnet</li>
            </ul>
          </div>
        )}

        {/* Backend Status and Stats */}
        <div style={{ 
          background: 'rgba(255, 255, 255, 0.95)', 
          borderRadius: '20px', 
          padding: '30px',
          marginTop: '20px'
        }}>
          <h2>
            📊 System Status 
            <span style={{ 
              marginLeft: '15px', 
              fontSize: '0.8em',
              padding: '5px 15px',
              borderRadius: '20px',
              background: apiStatus === 'online' ? '#48bb78' : '#f56565',
              color: 'white'
            }}>
              {apiStatus === 'online' ? '🟢 Online' : '🔴 Offline'}
            </span>
          </h2>
          
          {stats && (
            <div className="stats-grid" style={{ 
              display: 'grid', 
              gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))',
              gap: '15px',
              marginTop: '20px'
            }}>
              <div className="stat-card" style={{
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                color: 'white',
                padding: '20px',
                borderRadius: '15px'
              }}>
                <div style={{ fontSize: '2em', fontWeight: 'bold' }}>{stats.total_spikes}</div>
                <div>Total Wicks Detected</div>
              </div>
              <div className="stat-card" style={{
                background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
                color: 'white',
                padding: '20px',
                borderRadius: '15px'
              }}>
                <div style={{ fontSize: '2em', fontWeight: 'bold' }}>{stats.total_payouts}</div>
                <div>Total Payouts</div>
              </div>
              <div className="stat-card" style={{
                background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
                color: 'white',
                padding: '20px',
                borderRadius: '15px'
              }}>
                <div style={{ fontSize: '2em', fontWeight: 'bold' }}>{stats.active_policies}</div>
                <div>Active Policies</div>
              </div>
              <div className="stat-card" style={{
                background: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
                color: 'white',
                padding: '20px',
                borderRadius: '15px'
              }}>
                <div style={{ fontSize: '2em', fontWeight: 'bold' }}>{stats.total_prices}</div>
                <div>Price Records</div>
              </div>
            </div>
          )}
        </div>

        {/* Recent Wicks */}
        {apiStatus === 'online' && spikes.length > 0 && (
          <div style={{ 
            background: 'rgba(255, 255, 255, 0.95)', 
            borderRadius: '20px', 
            padding: '30px',
            marginTop: '20px'
          }}>
            <h2>⚡ Recent Wick Detections</h2>
            <div style={{ overflowX: 'auto' }}>
              <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: '15px' }}>
                <thead>
                  <tr style={{ borderBottom: '2px solid #e2e8f0' }}>
                    <th style={{ padding: '12px', textAlign: 'left' }}>Time</th>
                    <th style={{ padding: '12px', textAlign: 'left' }}>Symbol</th>
                    <th style={{ padding: '12px', textAlign: 'right' }}>
                      Body Ratio 
                      <span className="tooltip" style={{ marginLeft: '5px' }}>
                        ❓
                        <span className="tooltiptext">|open-close| / (high-low)</span>
                      </span>
                    </th>
                    <th style={{ padding: '12px', textAlign: 'right' }}>
                      Range Ratio 
                      <span className="tooltip" style={{ marginLeft: '5px' }}>
                        ❓
                        <span className="tooltiptext">(high-low) / close</span>
                      </span>
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {spikes.slice(0, 5).map((spike, idx) => {
                    return (
                      <tr key={idx} style={{ borderBottom: '1px solid #e2e8f0' }}>
                        <td style={{ padding: '12px' }}>
                          {formatDate(new Date(spike.Timestamp))}
                        </td>
                        <td style={{ padding: '12px' }}>{spike.Symbol}</td>
                        <td style={{ padding: '12px', textAlign: 'right', color: '#667eea', fontWeight: 'bold' }}>
                          {(spike.BodyRatio * 100)?.toFixed(2)}%
                        </td>
                        <td style={{ padding: '12px', textAlign: 'right', color: '#f56565', fontWeight: 'bold' }}>
                          {(spike.RangeClosePercent * 100)?.toFixed(2)}%
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* Recent Prices - K-line Chart */}
        {apiStatus === 'online' && prices.length > 0 && (
          <div style={{ 
            background: 'rgba(255, 255, 255, 0.95)', 
            borderRadius: '20px', 
            padding: '30px',
            marginTop: '20px'
          }}>
            <h2>📈 Recent Price Data</h2>
            <PriceChart prices={prices} />
          </div>
        )}

        {/* Payout History */}
        {apiStatus === 'online' && payouts.length > 0 && (
          <div style={{ 
            background: 'rgba(255, 255, 255, 0.95)', 
            borderRadius: '20px', 
            padding: '30px',
            marginTop: '20px'
          }}>
            <h2>💰 Payout History</h2>
            <div style={{ overflowX: 'auto' }}>
              <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: '15px' }}>
                <thead>
                  <tr style={{ borderBottom: '2px solid #e2e8f0' }}>
                    <th style={{ padding: '12px', textAlign: 'left' }}>Time</th>
                    <th style={{ padding: '12px', textAlign: 'left' }}>User</th>
                    <th style={{ padding: '12px', textAlign: 'right' }}>Amount</th>
                    <th style={{ padding: '12px', textAlign: 'left' }}>Tx Hash</th>
                  </tr>
                </thead>
                <tbody>
                  {payouts.slice(0, 5).map((payout, idx) => (
                    <tr key={idx} style={{ borderBottom: '1px solid #e2e8f0' }}>
                      <td style={{ padding: '12px' }}>
                        {formatDate(new Date(payout.ExecutedAt))}
                      </td>
                      <td style={{ padding: '12px' }}>
                        {payout.UserAddress ? `${payout.UserAddress.slice(0, 6)}...${payout.UserAddress.slice(-4)}` : 'Unknown'}
                      </td>
                      <td style={{ padding: '12px', textAlign: 'right', color: '#48bb78', fontWeight: 'bold' }}>
                        ${payout.Amount?.toFixed(2)}
                      </td>
                      <td style={{ padding: '12px' }}>
                        {payout.TxHash ? (
                          <a href={`https://sepolia.etherscan.io/tx/${payout.TxHash}`}
                             target="_blank"
                             rel="noopener noreferrer"
                             style={{ color: '#667eea' }}>
                            {payout.TxHash.substring(0, 10)}...
                          </a>
                        ) : 'Pending'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
