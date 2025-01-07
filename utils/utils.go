package utils

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// DataInput provides data input string captured via options
func DataInput(data string) string {
	var inputData string

	// Check if -data flag is provided
	if data != "" {
		inputData = data
	} else if len(flag.Args()) > 0 { // Check for positional arguments
		inputData = flag.Args()[0]
	} else { // Check if data is piped via stdin
		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Println("Error checking stdin:", err)
			os.Exit(1)
		}

		if (stat.Mode() & os.ModeCharDevice) == 0 { // Data is piped
			reader := bufio.NewReader(os.Stdin)
			pipeData, err := reader.ReadString('\n')
			if err != nil && err.Error() != "EOF" {
				fmt.Println("Error reading stdin:", err)
				os.Exit(1)
			}
			inputData = strings.TrimSpace(pipeData) // Trim extra spaces or newlines
		}
	}

	// If no data was provided
	if inputData == "" {
		fmt.Println("No data provided")
		os.Exit(1)
	}

	return inputData

}
