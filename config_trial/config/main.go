package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// read YAML

func main() {
	viper.SetConfigFile("./example_config.yml")
	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
}
