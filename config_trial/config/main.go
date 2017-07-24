package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// read YAML

func main() {
	viper.SetConfigFile("./dd_example.yml")
	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	pathToDockerfile := viper.Get("pathToDockerfile")
	fmt.Printf("pathToDockerfile: %v\n", pathToDockerfile)

	endPointName := viper.Get("endPointName")
	fmt.Printf("endPointName: %v\n", endPointName)

	inputPath := viper.Get("inputPath")
	fmt.Printf("inputPath: %v\n", inputPath)

	hostIP := viper.Get("hostIP")
	fmt.Printf("hostIP: %v\n", hostIP)

	hostPort := viper.Get("hostPort")
	fmt.Printf("hostPort: %v\n", hostPort)

	containerName := viper.Get("containerName")
	fmt.Printf("containerName: %v\n", containerName)

	userName := viper.Get("userName")
	fmt.Printf("userName: %v\n", userName)

	imgName := viper.Get("imgName")
	fmt.Printf("imgName: %v\n", imgName)

	uN := fmt.Sprintf("%v", userName)
	iN := fmt.Sprintf("%v", imgName)

	imgHandle := uN + "/" + iN
	fmt.Printf("imgHandle: %v\n", imgHandle)
}
