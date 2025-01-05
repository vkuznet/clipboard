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
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	item := flag.String("item", "", "clipboard item")
	flag.Parse()

	config, err := utils.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	client := utils.Client(config)
	rurl := fmt.Sprintf("%s/paste?item=%s", utils.ServerUrl(config), *item)
	resp, err := client.Get(rurl)
	//     resp, err := http.Get(fmt.Sprintf("%s/paste", *server))
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
