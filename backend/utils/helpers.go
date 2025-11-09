package utils

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	RPC struct {
		URL             string `yaml:"url"`
		ContractAddress string `yaml:"contract_address"`
		PrivateKey      string `yaml:"private_key"`
	} `yaml:"rpc"`

	Detector struct {
		ThresholdPercent float64 `yaml:"threshold_percent"`
		WindowMinutes    int     `yaml:"window_minutes"`
	} `yaml:"detector"`

	Chainlink struct {
		BtcUsdFeed     string `yaml:"btc_usd_feed"`
		UpdateInterval int    `yaml:"update_interval"`
	} `yaml:"chainlink"`

	Mode string `yaml:"mode"`
}

var AppConfig *Config

// LoadConfig loads configuration from YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	AppConfig = &config
	return &config, nil
}

// LogInfo prints info level log
func LogInfo(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

// LogError prints error level log
func LogError(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// LogDebug prints debug level log
func LogDebug(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}
