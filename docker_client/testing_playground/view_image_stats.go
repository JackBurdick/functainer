package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {

		// Select specified image by the repo tag.
		if strings.Join(image.RepoTags, "") == "jackburdick/cosineexp:latest" {
			fmt.Printf("Command: %v\n", image.Containers)
			fmt.Printf("Created: %v\n", image.Created)
			fmt.Printf("ID: %v\n", image.ID)
			fmt.Printf("Labels: %v\n", image.Labels)
			fmt.Printf("ParentID: %v\n", image.ParentID)
			fmt.Printf("RepoDigests: %v\n", image.RepoDigests)
			fmt.Printf("RepoTags: %v\n", image.RepoTags)
			fmt.Printf("SharedSize: %v\n", image.SharedSize)
			fmt.Printf("Size: %v\n", image.Size)
			fmt.Printf("VirtualSize: %v\n", image.VirtualSize)
			fmt.Printf("----------------------------------------\n")
		}
	}
}

// func (cli *Client) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
