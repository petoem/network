package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	configPath := flag.String("config", "", "path to a json file containing the config")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "error no config file path provided")
		os.Exit(1)
	}

	configRaw, err := os.ReadFile(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("could not read config file: %w", err))
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(configRaw, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("could not parse json: %w", err))
		os.Exit(1)
	}

	err = NewAnnouncer(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("announcer: %w", err))
		os.Exit(1)
	}
}
