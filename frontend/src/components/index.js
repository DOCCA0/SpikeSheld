import React from 'react';

export const WalletConnect = ({ isConnected, account, onConnect, onDisconnect, loading }) => {
  if (!isConnected) {
    return (
      <div className="wallet-section">
        <h2>Connect Your Wallet</h2>
        <p>Connect your wallet to purchase spike insurance</p>
        <button 
          className="button" 
          onClick={onConnect}
          disabled={loading}
        >
          {loading ? <span className="loading"></span> : 'Connect Wallet'}
        </button>
      </div>
    );
  }

  return (
    <div className="wallet-section">
      <h2>Wallet Connected</h2>
      <div className="account-info">
        <div className="account-address">
          {account.substring(0, 6)}...{account.substring(account.length - 4)}
        </div>
      </div>
      <button 
        className="button button-danger" 
        onClick={onDisconnect}
      >
        Disconnect
      </button>
    </div>
  );
};

export const PolicyCard = ({ policy, index }) => {
  const getPolicyStatus = () => {
    if (policy.claimed) return 'claimed';
    if (policy.expiryTime < new Date()) return 'expired';
    if (policy.active) return 'active';
    return 'inactive';
  };

  const status = getPolicyStatus();

  return (
    <div className={`policy-card ${status}`}>
      <h4>
        Policy #{index}
        <span className={`status-badge status-${status}`} style={{ marginLeft: '10px' }}>
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
        <span className="policy-value">{policy.purchaseTime.toLocaleString()}</span>
      </div>
      <div className="policy-detail">
        <span className="policy-label">Expires:</span>
        <span className="policy-value">{policy.expiryTime.toLocaleString()}</span>
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
};

export const InsuranceInfo = () => {
  return (
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
        <h3>Threshold</h3>
        <div className="value">10% Drop</div>
      </div>
    </div>
  );
};

export const ErrorMessage = ({ message }) => {
  if (!message) return null;
  
  return (
    <div className="error-message">
      ❌ {message}
    </div>
  );
};

export const SuccessMessage = ({ message }) => {
  if (!message) return null;
  
  return (
    <div className="success-message">
      ✅ {message}
    </div>
  );
};
