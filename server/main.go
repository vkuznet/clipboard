package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	utils "github.com/vkuznet/clipboard/utils"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	config, err := utils.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	// Initialize secure clipboard
	clipboard, err := utils.NewSecureClipboard(config)
	if err != nil {
		fmt.Println("Error initializing secure clipboard:", err)
		os.Exit(1)
	}

	// Load persistent history
	if err := clipboard.LoadHistory(config.HistoryFile); err != nil && !os.IsNotExist(err) {
		fmt.Println("Failed to load clipboard history:", err)
		os.Exit(1)
	}
	defer clipboard.SaveHistory(config.HistoryFile)

	http.HandleFunc("/copy", clipboard.HandleCopy)
	http.HandleFunc("/paste", clipboard.HandlePaste)
	http.HandleFunc("/history", clipboard.HandleHistory)

	address := ":" + strconv.Itoa(config.Port)
	fmt.Printf("Server running on port %d\n", config.Port)

	if config.ServerKey != "" && config.ServerCert != "" {
		server := &http.Server{
			Addr: address,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}
		err = server.ListenAndServeTLS(config.ServerCert, config.ServerKey)
	} else {
		err = http.ListenAndServe(address, nil)
	}

	if err != nil {
		fmt.Println("Server error:", err)
		os.Exit(1)
	}
}
