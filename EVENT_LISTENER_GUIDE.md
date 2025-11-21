# Event Listener Implementation Guide

## Summary
Successfully implemented blockchain event listener for SpikeShield backend to monitor and sync PolicyPurchased and PayoutExecuted contract events.

## Components Added

### 1. Event Listener Service (backend/eventlistener/listener.go)
- Monitors InsurancePool contract events
- Polls blockchain at configurable intervals
- Processes PolicyPurchased and PayoutExecuted events
- Stores events in database with deduplication
- Tracks sync progress to avoid reprocessing

### 2. Database Schema Updates (backend/db/schema.sql)
- `contract_events` table: Stores all blockchain events
- `sync_state` table: Tracks last synced block number
- Added indexes for performance

### 3. Configuration (backend/config.yaml)
```yaml
eventlistener:
  enabled: true
  poll_interval: 30
```

### 4. Main Integration (backend/main.go)
- Automatic startup when enabled
- Background goroutine for event polling
- Graceful shutdown handling

## Features

1. **Automatic Synchronization**: Polls blockchain every 30 seconds (configurable)
2. **Event Handlers**: 
   - PolicyPurchased: Creates policy records
   - PayoutExecuted: Updates policy status and creates payout records
3. **Idempotent**: Prevents duplicate event processing
4. **Resumable**: Tracks last synced block, resumes on restart

## Usage

Start backend with event listener enabled:
```bash
cd backend
go run main.go
```

The listener will automatically:
- Sync historical events (last 1000 blocks on first run)
- Poll for new events continuously
- Store data in database tables

## Database Tables

- `contract_events`: Raw event data
- `policies`: User insurance policies
- `payouts`: Payout execution records
- `sync_state`: Synchronization checkpoint

## Testing

Monitor logs for event synchronization:
- Look for "🎧 Event listener started" message
- Query `contract_events` table to verify event storage- Check for "

