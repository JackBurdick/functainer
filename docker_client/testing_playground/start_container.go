package main

// TODO: expose ports

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// ---------- Build container from image
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	var contID string
	for _, image := range images {

		// Select specified image by the repo tag.
		if strings.Join(image.RepoTags, "") == "jackburdick/cosineexp:latest" {

			// Create container from image
			//cli.ContainerCreate(context.Background())
			configOptions := container.Config{Image: strings.Join(image.RepoTags, "")}
			hostConfig := container.HostConfig{PublishAllPorts: true}
			networkConfig := network.NetworkingConfig{}

			// TODO: Devise a progamatic way of producting a container name.
			containerName := "jacksss"
			createResponse, err := cli.ContainerCreate(context.Background(), &configOptions, &hostConfig, &networkConfig, containerName)
			if err != nil {
				fmt.Println(err)
			}
			// fmt.Printf("createResponse: %v", createResponse)
			// fmt.Printf("err: %v", err)
			contID = createResponse.ID
		}
	}

	// Start the container
	if contID != "" {
		err = cli.ContainerStart(context.Background(), contID, types.ContainerStartOptions{})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Error: no container created")
	}

}
