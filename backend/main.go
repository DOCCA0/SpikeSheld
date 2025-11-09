package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"spikeshield/api"
	"spikeshield/datafeed"
	"spikeshield/db"
	"spikeshield/detector"
	"spikeshield/utils"
)

func main() {
	// Parse command line flags
	mode := flag.String("mode", "replay", "Mode: replay or live")
	symbol := flag.String("symbol", "BTCUSDT", "Trading symbol")
	startTime := flag.String("start", "", "Start time for replay mode (e.g., 2021-05-19T00:00:00)")
	endTime := flag.String("end", "", "End time for replay mode (e.g., 2021-05-19T01:00:00)")
	configPath := flag.String("config", "config.yaml", "Path to config file")
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

	// Create detector
	det := detector.NewDetector(*symbol, config.Detector.ThresholdPercent, config.Detector.WindowMinutes)

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
		runReplayMode(*symbol, *startTime, *endTime, det, onSpikeDetected)
	case "live":
		runLiveMode(config, *symbol, det, onSpikeDetected)
	default:
		utils.LogError("Invalid mode: %s (use 'replay' or 'live')", *mode)
		os.Exit(1)
	}
}

// runReplayMode processes historical data from CSV
func runReplayMode(symbol, startStr, endStr string, det *detector.Detector, callback func(*db.Spike)) {
	utils.LogInfo("📊 Running in REPLAY mode")

	// Parse time range
	var start, end time.Time
	var err error

	if startStr == "" {
		// Default to a sample time range
		start = time.Date(2021, 5, 19, 0, 0, 0, 0, time.UTC)
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			start, err = time.Parse("2006-01-02T15:04:05", startStr)
			if err != nil {
				utils.LogError("Failed to parse start time: %v", err)
				os.Exit(1)
			}
		}
	}

	if endStr == "" {
		end = start.Add(24 * time.Hour)
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			end, err = time.Parse("2006-01-02T15:04:05", endStr)
			if err != nil {
				utils.LogError("Failed to parse end time: %v", err)
				os.Exit(1)
			}
		}
	}

	utils.LogInfo("Time range: %s to %s", start.Format(time.RFC3339), end.Format(time.RFC3339))

	// Load CSV data
	csvPath := fmt.Sprintf("../data/%s_%s.csv", symbol, start.Format("2006-01-02"))
	
	// Check if file exists, otherwise use default
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		csvPath = "../data/btcusdt_2021-05-19.csv"
		utils.LogInfo("Using default CSV: %s", csvPath)
	}

	feed := datafeed.NewReplayFeed(csvPath, symbol, start, end)
	if err := feed.LoadAndStore(); err != nil {
		utils.LogError("Failed to load replay data: %v", err)
		os.Exit(1)
	}

	// Run detection on loaded data
	spikes, err := det.DetectAllInRange(start, end)
	if err != nil {
		utils.LogError("Detection failed: %v", err)
		os.Exit(1)
	}

	utils.LogInfo("Detection complete: found %d spike(s)", len(spikes))

	// Trigger payouts for detected spikes
	for _, spike := range spikes {
		callback(spike)
	}

	utils.LogInfo("✅ Replay mode completed")
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

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start live feed in goroutine
	go func() {
		if err := liveFeed.Start(ctx); err != nil {
			utils.LogError("Live feed error: %v", err)
		}
	}()

	// Start continuous monitoring
	go det.ContinuousMonitor(30*time.Second, callback)

	utils.LogInfo("Live mode running... Press Ctrl+C to stop")

	// Wait for shutdown signal
	<-sigChan
	utils.LogInfo("Shutting down...")
	cancel()
	time.Sleep(1 * time.Second)
}
