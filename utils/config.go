package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config represents configuration object
type Config struct {
	Port          int    `json:"port"`
	ClipboardSize int    `json:"clipboard_size"`
	HistoryFile   string `json:"history_file"`
	ServerKey     string `json:"server_key"`
	ServerCert    string `json:"server_cert"`
	Secret        string `json:"clipboard_secret"`
}

// LoadConfig loads configuration from given file
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	if config.Secret == "" {
		// generate persistent secret with given number of bytes
		// it will be persistent across sessions
		if secret, err := GenerateSecret(32); err == nil {
			config.Secret = secret
		} else {
			log.Fatal(err)
		}
	}
	return &config, nil
}

// ConfigLocation provides default location of configuration file in user's home area
func ConfigLocation() string {
	return fmt.Sprintf("%s/.clipboard/config.json", os.Getenv("HOME"))
}
