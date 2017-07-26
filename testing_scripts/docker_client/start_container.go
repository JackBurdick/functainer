package main

// TODO: expose ports

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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

			// look into this helper function
			// portMap, bindingMap, err := nat.ParsePortSpecs([]string{"1234/tcp", "2345/udp"})

			exposedCadvPort := map[nat.Port]struct{}{"8080/tcp": {}}
			configOptions := container.Config{Image: strings.Join(image.RepoTags, ""), ExposedPorts: exposedCadvPort}

			networkConfig := network.NetworkingConfig{}

			//could also use 0.0.0.0, in place of 127.0.0.1 for local host
			// https://stackoverflow.com/questions/38834434/how-are-127-0-0-1-0-0-0-0-and-localhost-different
			portBindings := map[nat.Port][]nat.PortBinding{
				"8080/tcp": {{HostIP: "127.0.0.1", HostPort: "8000"}}}

			hostConfig := container.HostConfig{
				PublishAllPorts: true,
				PortBindings:    portBindings,
			}

			// TODO: Devise a progamatic way of producting a container name.

			// Create container
			containerName := "jacksss"
			createResponse, err := cli.ContainerCreate(context.Background(), &configOptions, &hostConfig, &networkConfig, containerName)
			if err != nil {
				fmt.Println(err)
			}
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

	// Stop the container
	// TODO: I have a stop time here, when the stoptime was nil, the processes
	// took noticeably long.. I need to investigate what, if any, problems this
	// causes.
	stopTime := time.Duration(100) * time.Millisecond
	err = cli.ContainerStop(context.Background(), contID, &stopTime)
	if err != nil {
		fmt.Println("ERROR: can't stop container")
	}
	fmt.Printf("id: %v, stopped?\n", contID)

	// TODO: Weigh the advantages of using the `force: true` flag here
	err = cli.ContainerRemove(context.Background(), contID, types.ContainerRemoveOptions{})
	if err != nil {
		fmt.Println("ERROR: can't remove container")
	}
	fmt.Printf("id: %v, removed?\n", contID)

	// Delete image
	// TODO: Can the image loop be removed?
	for _, image := range images {
		if strings.Join(image.RepoTags, "") == "jackburdick/cosineexp:latest" {
			fmt.Printf("image.ID: %v\n", image.ID)

			imgID := strings.TrimLeft(image.ID, "sha256")
			imgID = strings.TrimLeft(imgID, ":")
			deleteResponse, err := cli.ImageRemove(context.Background(), imgID, types.ImageRemoveOptions{})
			if err != nil {
				fmt.Printf("ERROR: image %v not deleted\n", imgID)
			}
			fmt.Printf("Image deleted: %v\n", imgID)
			fmt.Println(deleteResponse)
		}
	}

}
