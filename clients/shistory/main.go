package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

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

	client := utils.Client(config)
	rurl := fmt.Sprintf("%s/history", utils.ServerUrl(config))
	resp, err := client.Get(rurl)
	if err != nil {
		fmt.Println("Error fetching history:", err)
		return
	}
	defer resp.Body.Close()

	// generate our secure key for decoding our clipboard items
	skey := utils.Generate32ByteKey(config.Secret)

	if resp.StatusCode == http.StatusOK {
		var history []string
		json.NewDecoder(resp.Body).Decode(&history)
		fmt.Println("Clipboard History:")
		for i, item := range history {
			// encrypted item
			//             fmt.Printf("%d: %s\n", i+1, item)
			// decrypting our item using our secure key
			data, err := utils.Decrypt(skey, item)
			if err == nil {
				fmt.Printf("%d: %s\n", i+1, string(data))
			} else {
				fmt.Println("ERROR:", err)
			}
		}
	} else {
		fmt.Println("Failed to fetch history:", resp.Status)
	}
}
