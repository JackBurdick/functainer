package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("Command: %v\n", container.Command)
		fmt.Printf("Created: %v\n", container.Created)
		fmt.Printf("HostConfig: %v\n", container.HostConfig)
		fmt.Printf("ID: %v\n", container.ID)
		fmt.Printf("Image: %v\n", container.Image)
		fmt.Printf("ImageID: %v\n", container.ImageID[:10])
		fmt.Printf("Labels: %v\n", container.Labels)
		fmt.Printf("Mounts: %v\n", container.Mounts)
		fmt.Printf("Names: %v\n", container.Names)
		fmt.Printf("NetworkSettings: %v\n", container.NetworkSettings)
		fmt.Printf("Ports: %v\n", container.Ports)
		fmt.Printf("SizeRootFs: %v\n", container.SizeRootFs)
		fmt.Printf("SizeRw: %v\n", container.SizeRw)
		fmt.Printf("State: %v\n", container.State)
		fmt.Printf("Status: %v\n", container.Status)
	}
}
