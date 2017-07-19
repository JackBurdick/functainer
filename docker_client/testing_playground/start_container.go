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

			/*
				exposedCadvPort := map[nat.Port]struct{}{"8080/tcp": {}}
				configOptions := container.Config{Image: strings.Join(image.RepoTags, ""), ExposedPorts: exposedCadvPort}

				portBindings := map[nat.Port][]nat.PortBinding{
					"8080/tcp": {{HostIP: "0.0.0.0", HostPort: "8080"}}}

				hostConfig := container.HostConfig{
					PublishAllPorts: true,
					PortBindings:    portBindings,
				}

				produces:
					- 0.0.0.0:8080->8080/tcp

				CONTAINER ID        IMAGE                          COMMAND                CREATED             STATUS              PORTS                    NAMES
				ea494e9408a7        jackburdick/cosineexp:latest   "/bin/sh -c ./goapp"   2 minutes ago       Up 2 minutes        0.0.0.0:8080->8080/tcp   jacksss


			*/

			// look into this helper function
			// portMap, bindingMap, err := nat.ParsePortSpecs([]string{"1234/tcp", "2345/udp"})

			exposedCadvPort := map[nat.Port]struct{}{"8080/tcp": {}}
			configOptions := container.Config{Image: strings.Join(image.RepoTags, ""), ExposedPorts: exposedCadvPort}

			networkConfig := network.NetworkingConfig{}

			portBindings := map[nat.Port][]nat.PortBinding{
				"8080/tcp": {{HostIP: "0.0.0.0", HostPort: "8080"}}}

			hostConfig := container.HostConfig{
				PublishAllPorts: true,
				PortBindings:    portBindings,
			}

			// TODO: Devise a progamatic way of producting a container name.
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

}
