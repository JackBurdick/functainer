package main

// TODO: Error handling needs to be implemented.

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jhoonb/archivex"
)

// createTar creates a tar of the Dockerfile directory.
func createTar(pathToCreatedTarDir string, pathToDockerfile string) (string, error) {
	tar := new(archivex.TarFile)
	tar.Create(pathToCreatedTarDir)
	tar.AddAll(pathToDockerfile, false)
	tar.Close()
	return pathToCreatedTarDir + ".tar", nil
}

// buildImageFromTar creates the tar of the Dockerfile and directory.
func buildImageFromTar(tarPath string, imgHandle string, cli *client.Client) {
	dockerBuildContext, err := os.Open(tarPath)
	defer dockerBuildContext.Close()
	buildOptions := types.ImageBuildOptions{Tags: []string{imgHandle}}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	// NOTE: This needs to be here to ensure the image is built before we
	// start/run it. There is likely a more elegant way to handle this.
	_, err = ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

// buildContainerFromImage builds the container from the image.
func buildContainerFromImage(imgTag string, images []types.ImageSummary, cli *client.Client) (string, error) {
	var contID string
	for _, image := range images {

		// Select specified image by the repo tag.
		if strings.Join(image.RepoTags, "") == imgTag {

			// Create the container from the image.
			// TODO: Devise a progamatic way of producting a container name.
			// I'm not even sure if the container name is assigned right now.
			exposedPort := map[nat.Port]struct{}{"8080/tcp": {}}
			jack := make(map[string]string)
			jack["slippery"] = "fish"
			configOptions := container.Config{Image: strings.Join(image.RepoTags, ""), ExposedPorts: exposedPort, Labels: jack}
			networkConfig := network.NetworkingConfig{}
			portBindings := map[nat.Port][]nat.PortBinding{
				"8080/tcp": {{HostIP: "127.0.0.1", HostPort: "8000"}}}
			hostConfig := container.HostConfig{
				PublishAllPorts: true,
				PortBindings:    portBindings,
			}
			containerName := "dunnoman"
			createResponse, err := cli.ContainerCreate(context.Background(), &configOptions, &hostConfig, &networkConfig, containerName)
			if err != nil {
				fmt.Println(err)
			}
			contID = createResponse.ID
		}
	}
	return contID, nil
}

// startContainerByID starts the container.
func startContainerByID(contID string, cli *client.Client) {
	if contID != "" {
		err := cli.ContainerStart(context.Background(), contID, types.ContainerStartOptions{})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Error: no container created")
	}
}

// stopContainerByID stops the container.
func stopContainerByID(contID string, cli *client.Client) {

	// TODO: I have a stop time here, when the stoptime was nil, the processes
	// took noticeably long.. I need to investigate what, if any, problems this
	// causes.
	stopTime := time.Duration(100) * time.Millisecond
	err := cli.ContainerStop(context.Background(), contID, &stopTime)
	if err != nil {
		fmt.Println("ERROR: can't stop container")
	}
}

// removeContainerByID removes the container.
func removeContainerByID(contID string, cli *client.Client) {
	// TODO: Weigh the advantages of using the `force: true` flag here
	err := cli.ContainerRemove(context.Background(), contID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		fmt.Println("ERROR: can't remove container")
	}
}

// buildContainerFromImage
func deleteImageByTag(imgTag string, images []types.ImageSummary, cli *client.Client) {
	// TODO: Can the image loop be removed?
	for _, image := range images {
		if strings.Join(image.RepoTags, "") == imgTag {
			imgID := strings.TrimLeft(image.ID, "sha256")
			imgID = strings.TrimLeft(imgID, ":")
			_, err := cli.ImageRemove(context.Background(), imgID, types.ImageRemoveOptions{})
			if err != nil {
				fmt.Printf("ERROR: image %v not deleted\n", imgID)
			}
			//fmt.Printf("Image deleted: %v\n", imgID)
			//fmt.Println(deleteResponse)
		}
	}
}

func main() {

	// Constants, these will hopefully eventually come from a YAML file.
	imgHandle := "jackburdick/automated"
	imgTag := imgHandle + ":latest"
	pathToDockerfile := "../../CaaF/Experimental/container/"
	pathToCreatedTarDir := "./test_arch/archieve"

	// Create tar.
	tarPath, err := createTar(pathToCreatedTarDir, pathToDockerfile)
	if err != nil {
		fmt.Printf("Unable to create tar: %v", err)
	}

	// Create docker cli environment.
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// Build image from the tar file.
	buildImageFromTar(tarPath, imgHandle, cli)

	// TODO: see if I can do this without the loop.
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	// Build container and get container id.
	contID, err := buildContainerFromImage(imgTag, images, cli)

	// Start the container
	startContainerByID(contID, cli)

	// Stop the container.
	stopContainerByID(contID, cli)

	// Remove the container.
	removeContainerByID(contID, cli)

	// Delete the image.
	deleteImageByTag(imgTag, images, cli)

	// Prune container.
	//cli.ContainersPrune()
	// jack := make(map[string]string)
	// jack["slippery"] = "fish"
	cPruneFilter := filters.NewArgs()
	//cPruneFilter.Add("label", "slippery=fish")
	//cPruneFilter.Add("dangling", "true")
	cPruneResult, err := cli.ContainersPrune(context.Background(), cPruneFilter)
	if err != nil {
		fmt.Printf("Error prune container: %v\n", err)
	}
	fmt.Println(cPruneResult)

	// Prune image.
	//cli.ImagesPrune()
	iPruneResult, err := cli.ImagesPrune(context.Background(), cPruneFilter)
	if err != nil {
		fmt.Printf("Error prune container: %v\n", err)
	}
	fmt.Println(iPruneResult)

}
