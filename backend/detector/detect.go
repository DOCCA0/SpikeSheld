package detector

import (
	"fmt"
	"time"

	"spikeshield/db"
	"spikeshield/utils"
)

// Detector monitors price changes and detects spikes
type Detector struct {
	ThresholdPercent float64
	WindowMinutes    int
	Symbol           string
}

// NewDetector creates a new detector instance
func NewDetector(symbol string, thresholdPercent float64, windowMinutes int) *Detector {
	return &Detector{
		ThresholdPercent: thresholdPercent,
		WindowMinutes:    windowMinutes,
		Symbol:           symbol,
	}
}

// CheckForSpike analyzes recent prices and detects if a spike occurred
func (d *Detector) CheckForSpike() (*db.Spike, error) {
	// Get latest price
	latest, err := db.GetLatestPrice(d.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest price: %w", err)
	}

	// Get price from N minutes ago
	windowStart := latest.Timestamp.Add(-time.Duration(d.WindowMinutes) * time.Minute)
	prices, err := db.GetPricesBetween(d.Symbol, windowStart, latest.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical prices: %w", err)
	}

	if len(prices) < 2 {
		// Not enough data to detect spike
		return nil, nil
	}

	// Find the highest price in the window
	var maxPrice float64
	var maxPriceTime time.Time
	
	for _, p := range prices {
		if p.Close > maxPrice {
			maxPrice = p.Close
			maxPriceTime = p.Timestamp
		}
	}

	// Calculate drop percentage
	dropPercent := ((maxPrice - latest.Close) / maxPrice) * 100

	utils.LogDebug("Price check: max=$%.2f (at %s), current=$%.2f, drop=%.2f%%",
		maxPrice, maxPriceTime.Format(time.RFC3339), latest.Close, dropPercent)

	// Check if drop exceeds threshold
	if dropPercent >= d.ThresholdPercent {
		spike := &db.Spike{
			Timestamp:   latest.Timestamp,
			Symbol:      d.Symbol,
			PriceBefore: maxPrice,
			PriceAfter:  latest.Close,
			DropPercent: dropPercent,
		}

		// Save spike to database
		if err := db.InsertSpike(spike); err != nil {
			return nil, fmt.Errorf("failed to insert spike: %w", err)
		}

		utils.LogInfo("🚨 SPIKE DETECTED! %s dropped %.2f%% from $%.2f to $%.2f",
			d.Symbol, dropPercent, maxPrice, latest.Close)

		return spike, nil
	}

	return nil, nil
}

// ContinuousMonitor runs spike detection in a loop
func (d *Detector) ContinuousMonitor(checkInterval time.Duration, onSpike func(*db.Spike)) {
	utils.LogInfo("Starting continuous spike monitoring for %s (threshold: %.2f%%)", d.Symbol, d.ThresholdPercent)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		spike, err := d.CheckForSpike()
		if err != nil {
			utils.LogError("Detection error: %v", err)
			continue
		}

		if spike != nil && onSpike != nil {
			// Trigger callback when spike is detected
			onSpike(spike)
		}
	}
}

// DetectAllInRange analyzes all price data in a time range (for replay mode)
func (d *Detector) DetectAllInRange(start, end time.Time) ([]*db.Spike, error) {
	utils.LogInfo("Analyzing price data from %s to %s", start.Format(time.RFC3339), end.Format(time.RFC3339))

	prices, err := db.GetPricesBetween(d.Symbol, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get prices: %w", err)
	}

	if len(prices) < 2 {
		return nil, fmt.Errorf("insufficient price data")
	}

	var spikes []*db.Spike
	windowSize := d.WindowMinutes

	// Scan through all prices looking for spikes
	for i := windowSize; i < len(prices); i++ {
		currentPrice := prices[i]
		
		// Look back at window
		var maxPrice float64
		for j := i - windowSize; j < i; j++ {
			if prices[j].Close > maxPrice {
				maxPrice = prices[j].Close
			}
		}

		// Calculate drop
		dropPercent := ((maxPrice - currentPrice.Close) / maxPrice) * 100

		if dropPercent >= d.ThresholdPercent {
			spike := &db.Spike{
				Timestamp:   currentPrice.Timestamp,
				Symbol:      d.Symbol,
				PriceBefore: maxPrice,
				PriceAfter:  currentPrice.Close,
				DropPercent: dropPercent,
			}

			if err := db.InsertSpike(spike); err != nil {
				utils.LogError("Failed to insert spike: %v", err)
				continue
			}

			spikes = append(spikes, spike)
			utils.LogInfo("Spike detected at %s: %.2f%% drop from $%.2f to $%.2f",
				spike.Timestamp.Format(time.RFC3339), dropPercent, maxPrice, currentPrice.Close)
		}
	}

	utils.LogInfo("Analysis complete: found %d spike(s)", len(spikes))
	return spikes, nil
}
