package datafeed

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"spikeshield/db"
	"spikeshield/utils"
)

// ReplayFeed reads historical price data from CSV file
type ReplayFeed struct {
	FilePath string
	Symbol   string
	Start    time.Time
	End      time.Time
}

// NewReplayFeed creates a new replay feed instance
func NewReplayFeed(filePath, symbol string, start, end time.Time) *ReplayFeed {
	return &ReplayFeed{
		FilePath: filePath,
		Symbol:   symbol,
		Start:    start,
		End:      end,
	}
}

// LoadAndStore reads CSV and stores price data in database
func (rf *ReplayFeed) LoadAndStore() error {
	file, err := os.Open(rf.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	utils.LogInfo("Loading price data from %s", rf.FilePath)
	count := 0

	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}

		// Parse CSV: timestamp, open, high, low, close, volume
		if len(record) < 6 {
			continue
		}

		timestamp, err := parseTimestamp(record[0])
		if err != nil {
			utils.LogError("Failed to parse timestamp: %v", err)
			continue
		}

		// Filter by time range
		if timestamp.Before(rf.Start) || timestamp.After(rf.End) {
			continue
		}

		open, _ := strconv.ParseFloat(record[1], 64)
		high, _ := strconv.ParseFloat(record[2], 64)
		low, _ := strconv.ParseFloat(record[3], 64)
		close, _ := strconv.ParseFloat(record[4], 64)
		volume, _ := strconv.ParseFloat(record[5], 64)

		priceData := &db.PriceData{
			Timestamp: timestamp,
			Symbol:    rf.Symbol,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}

		if err := db.InsertPrice(priceData); err != nil {
			utils.LogError("Failed to insert price: %v", err)
			continue
		}

		count++
	}

	utils.LogInfo("Loaded %d price records", count)
	return nil
}

// parseTimestamp converts string to time.Time
// Supports multiple formats: "2006-01-02 15:04:05", "2006-01-02T15:04:05Z", Unix timestamp
func parseTimestamp(s string) (time.Time, error) {
	// Try ISO format first
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	// Try Unix timestamp
	if timestamp, err := strconv.ParseInt(s, 10, 64); err == nil {
		// Check if milliseconds or seconds
		if timestamp > 1e12 {
			return time.Unix(0, timestamp*int64(time.Millisecond)), nil
		}
		return time.Unix(timestamp, 0), nil
	}

	return time.Time{}, fmt.Errorf("unsupported timestamp format: %s", s)
}
