package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	utils "github.com/vkuznet/clipboard/utils"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	data := flag.String("data", "", "Data to copy")
	flag.Parse()

	if *data == "" {
		fmt.Println("Please provide data to copy using the -data flag")
		return
	}

	config, err := utils.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	client := utils.Client(config)
	rurl := fmt.Sprintf("%s/copy?data=%s", utils.ServerUrl(config), url.QueryEscape(*data))
	resp, err := client.Get(rurl)
	//     resp, err := http.Get(fmt.Sprintf("%s/copy?data=%s", *server, url.QueryEscape(*data)))
	if err != nil {
		fmt.Println("Error copying data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Data successfully copied")
	} else {
		fmt.Println("Failed to copy data:", resp.Status)
	}
}
