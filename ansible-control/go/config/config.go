package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Kafka KafkaConfig `json:"kafka"`
}

// KafkaConfig represents Kafka configuration
type KafkaConfig struct {
	Broker        string `json:"broker"`
	ConsumerTopic string `json:"consumer_topic"`
	ProducerTopic string `json:"producer_topic"`
}

// LoadConfig loads configuration from config.json file
func LoadConfig() (*Config, error) {
	// Get the config file path (relative to the executable or source)
	configPath := getConfigPath()

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.Kafka.Broker == "" {
		return nil, fmt.Errorf("kafka.broker is required in config file")
	}
	if config.Kafka.ConsumerTopic == "" {
		return nil, fmt.Errorf("kafka.consumer_topic is required in config file")
	}
	if config.Kafka.ProducerTopic == "" {
		return nil, fmt.Errorf("kafka.producer_topic is required in config file")
	}

	return &config, nil
}

// getConfigPath returns the path to config.json
// Tries multiple locations:
// 1. config/config.json (relative to executable)
// 2. ./config/config.json (relative to current directory)
func getConfigPath() string {
	// Try relative to executable first
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		configPath := filepath.Join(exeDir, "config", "config.json")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
		// Also try one level up (if exe is in bin/)
		parentConfigPath := filepath.Join(filepath.Dir(exeDir), "ansible-control", "go", "config", "config.json")
		if _, err := os.Stat(parentConfigPath); err == nil {
			return parentConfigPath
		}
	}

	// Try relative to current working directory
	wd, err := os.Getwd()
	if err == nil {
		configPath := filepath.Join(wd, "config", "config.json")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
		// Try go/config/config.json
		goConfigPath := filepath.Join(wd, "go", "config", "config.json")
		if _, err := os.Stat(goConfigPath); err == nil {
			return goConfigPath
		}
	}

	// Default fallback
	return "config/config.json"
}
