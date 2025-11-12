import React, { useEffect, useState } from 'react';
import { useContract } from '../hooks/useContract';

/**
 * PayoutNotification component
 * Listens for payout events and displays notifications when user receives payouts
 */
const PayoutNotification = () => {
  const { listenForPayouts, isConnected } = useContract();
  const [payouts, setPayouts] = useState([]);

  useEffect(() => {
    if (!isConnected) return;

    // Listen for payout events
    const unsubscribe = listenForPayouts((payoutData) => {
      // Add new payout to the list
      setPayouts(prev => [payoutData, ...prev]);

      // Show browser notification if permitted
      if ('Notification' in window && Notification.permission === 'granted') {
        new Notification('💰 Insurance Payout Received!', {
          body: `You received ${payoutData.amount} USDT for policy #${payoutData.policyId}`,
          icon: '/logo192.png'
        });
      }
    });

    // Request notification permission
    if ('Notification' in window && Notification.permission === 'default') {
      Notification.requestPermission();
    }

    // Cleanup listener on unmount
    return () => {
      if (unsubscribe) unsubscribe();
    };
  }, [isConnected, listenForPayouts]);

  if (!isConnected || payouts.length === 0) return null;

  return (
    <div style={{
      position: 'fixed',
      top: '20px',
      right: '20px',
      zIndex: 1000,
      maxWidth: '400px'
    }}>
      {payouts.map((payout, index) => (
        <div
          key={`${payout.txHash}-${index}`}
          style={{
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            color: 'white',
            padding: '15px',
            borderRadius: '10px',
            marginBottom: '10px',
            boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
            animation: 'slideIn 0.3s ease-out'
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
            <span style={{ fontSize: '24px', marginRight: '10px' }}>🎉</span>
            <strong>Payout Received!</strong>
          </div>
          <div style={{ fontSize: '14px', opacity: 0.9 }}>
            Amount: <strong>{payout.amount} USDT</strong>
          </div>
          <div style={{ fontSize: '12px', opacity: 0.8, marginTop: '5px' }}>
            Policy #{payout.policyId} | Block #{payout.blockNumber}
          </div>
        </div>
      ))}
    </div>
  );
};

export default PayoutNotification;
