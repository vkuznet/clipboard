package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	utils "github.com/vkuznet/clipboard/utils"
)

func main() {
	configPath := flag.String("config", utils.ConfigLocation(), "configuration file")
	item := flag.String("item", "", "clipboard item")
	flag.Parse()

	// Handle positional arguments (data passed without -item flag)
	inputData := utils.DataInput(*item)

	config, err := utils.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	client := utils.Client(config)
	rurl := fmt.Sprintf("%s/paste?item=%s", utils.ServerUrl(config), inputData)
	resp, err := client.Get(rurl)
	if err != nil {
		fmt.Println("Error pasting data:", err)
		return
	}
	defer resp.Body.Close()

	// generate our secure key for decoding our clipboard items
	skey := utils.Generate32ByteKey(config.Secret)

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		// decrypt our encrypted item
		if data, err := utils.Decrypt(skey, string(body)); err == nil {
			fmt.Printf("%s\n", data)
		} else {
			fmt.Println("ERROR:", err)
		}
	} else {
		fmt.Println("Failed to paste data:", resp.Status)
	}
}
