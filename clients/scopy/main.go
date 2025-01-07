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
	configPath := flag.String("config", utils.ConfigLocation(), "configuration file")
	data := flag.String("data", "", "Data to copy")
	flag.Parse()

	// Handle positional arguments (data passed without -data flag)
	var inputData string
	if *data != "" {
		inputData = *data
	} else if len(flag.Args()) > 0 {
		inputData = flag.Args()[0]
	} else {
		fmt.Println("No data provided")
		os.Exit(1)
	}

	config, err := utils.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	client := utils.Client(config)
	rurl := fmt.Sprintf("%s/copy?data=%s", utils.ServerUrl(config), url.QueryEscape(inputData))
	resp, err := client.Get(rurl)
	if err != nil {
		fmt.Println("Error copying data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		//         fmt.Println("Data successfully copied to secure clipboard")
	} else {
		fmt.Println("Failed to copy data:", resp.Status)
	}
}
