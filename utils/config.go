package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port          int    `json:"port"`
	ClipboardSize int    `json:"clipboard_size"`
	HistoryFile   string `json:"history_file"`
	ServerKey     string `json:"server_key"`
	ServerCert    string `json:"server_cert"`
	Secret        string `json:"clipboard_secret"`
}

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
	return &config, nil
}
