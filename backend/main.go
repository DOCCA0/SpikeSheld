package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"spikeshield/api"
	"spikeshield/datafeed"
	"spikeshield/db"
	"spikeshield/detector"
	"spikeshield/eventlistener"
	"spikeshield/utils"
)

func main() {
	// Parse command line flags
	mode := flag.String("mode", "replay", "Mode: replay or live")
	symbol := flag.String("symbol", "BTCUSDT", "Trading symbol")
	configPath := flag.String("config", "config.yaml", "Path to config file")
	apiPort := flag.String("api-port", "8080", "API server port")
	flag.Parse()

	utils.LogInfo("🚀 SpikeShield Starting...")
	utils.LogInfo("Mode: %s, Symbol: %s", *mode, *symbol)

	// Load configuration
	config, err := utils.LoadConfig(*configPath)
	if err != nil {
		utils.LogError("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Connect to database
	if err := db.Connect(config); err != nil {
		utils.LogError("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize insert notification channel
	db.InitNotifier()

	// Start API server in background
	apiServer := api.NewServer(":" + *apiPort)
	go func() {
		if err := apiServer.Start(); err != nil {
			utils.LogError("API server error: %v", err)
		}
	}()

	// Start event listener in background if enabled
	var listener *eventlistener.EventListener
	if config.EventListener.Enabled {
		pollInterval := time.Duration(config.EventListener.PollInterval) * time.Second
		listener, err = eventlistener.NewEventListener(config.RPC.URL, config.RPC.ContractAddress, pollInterval)
		if err != nil {
			utils.LogError("Failed to create event listener: %v", err)
		} else {
			utils.LogInfo("Starting event listener (poll interval: %ds)", config.EventListener.PollInterval)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				if err := listener.Start(ctx); err != nil && err != context.Canceled {
					utils.LogError("Event listener error: %v", err)
				}
			}()
			defer listener.Close()
		}
	}

	// Create detector
	det := detector.NewDetector(*symbol, config.Detector.ThresholdPercent, config.Detector.BodyRatioMax)

	// Create payout service
	payoutSvc, err := api.NewPayoutService(config.RPC.URL, config.RPC.ContractAddress, config.RPC.PrivateKey)
	if err != nil {
		utils.LogError("Failed to create payout service: %v", err)
		// Continue without payout service for demo
		payoutSvc = nil
	}
	if payoutSvc != nil {
		defer payoutSvc.Close()
	}

	// Spike callback - triggers payout when spike is detected
	onSpikeDetected := func(spike *db.Spike) {
		utils.LogInfo("🔔 Spike callback triggered")
		if payoutSvc != nil {
			if err := payoutSvc.ExecutePayout(spike); err != nil {
				utils.LogError("Failed to execute payout: %v", err)
			}
		} else {
			utils.LogInfo("Payout service not available, skipping payout execution")
		}
	}

	// Run based on mode
	switch *mode {
	case "replay":
		runReplayMode(det, onSpikeDetected)
	case "live":
		runLiveMode(config, *symbol, det, onSpikeDetected)
	default:
		utils.LogError("Invalid mode: %s (use 'replay' or 'live')", *mode)
		os.Exit(1)
	}
}

// monitorDatabaseInserts listens for new rows in the database and triggers detection
func monitorDatabaseInserts(det *detector.Detector, callback func(*db.Spike)) {
	utils.LogInfo("👀 Monitoring database for new inserts via channel...")

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	utils.LogInfo("Monitoring active. Press Ctrl+C to stop")

	// Listen for insert notifications
	for {
		select {
		case <-db.InsertNotifier:
			utils.LogInfo("📥 New data inserted, checking for spikes...")
			// Run detection on latest price
			spike, err := det.CheckForSpike()
			if err != nil {
				utils.LogError("Detection failed: %v", err)
				continue
			}

			// Trigger callback if spike detected
			if spike != nil {
				callback(spike)
			}

		case <-sigChan:
			utils.LogInfo("Shutting down...")
			return
		}
	}
}

// runReplayMode waits for external script to insert data row by row
func runReplayMode(det *detector.Detector, callback func(*db.Spike)) {
	utils.LogInfo("📊 Running in REPLAY mode")
	utils.LogInfo("Waiting for external script to insert data from CSV...")

	// Just monitor database inserts (external script will feed data)
	monitorDatabaseInserts(det, callback)
}

// runLiveMode monitors real-time price from Chainlink
func runLiveMode(config *utils.Config, symbol string, det *detector.Detector, callback func(*db.Spike)) {
	utils.LogInfo("⚡ Running in LIVE mode")

	// Create live feed
	liveFeed, err := datafeed.NewLiveFeed(
		config.RPC.URL,
		config.Chainlink.BtcUsdFeed,
		symbol,
		config.Chainlink.UpdateInterval,
	)
	if err != nil {
		utils.LogError("Failed to create live feed: %v", err)
		os.Exit(1)
	}
	defer liveFeed.Close()

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start live feed in background to populate database
	go func() {
		if err := liveFeed.Start(ctx); err != nil {
			utils.LogError("Live feed error: %v", err)
		}
	}()

	utils.LogInfo("Live feed started, now monitoring database...")

	// Monitor database inserts (same as replay mode)
	monitorDatabaseInserts(det, callback)

	// Cancel live feed context on exit
	cancel()
}
