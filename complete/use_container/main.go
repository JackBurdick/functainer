package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jdkato/prose/tokenize"
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
		}
	}
}

func main() {

	// Constants, these will hopefully eventually come from a YAML file.
	imgHandle := "jackburdick/automated"
	imgTag := imgHandle + ":latest"
	pathToDockerfile := "../container/"
	pathToCreatedTarDir := "./archive/archive"

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
	if err != nil {
		fmt.Printf("error > build container: %v\n", err)
	}

	// Start the container
	startContainerByID(contID, cli)

	// ----------------------------------- use container

	// Input data (directory of files to compare)
	inputDir := "./input/"

	// URL endpoints.
	URL := "http://localhost:8000/cosineSim"
	//URL := "http://localhost:8000/"

	// Create map from input directory.
	fileMap, err := createMap(inputDir)
	if err != nil {
		fmt.Println(err)
	}

	// Build json data from map.
	jsonData, err := json.Marshal(fileMap)
	if err != nil {
		fmt.Println(err)
	}

	// Gzip json data.
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(jsonData)
	if err != nil {
		fmt.Printf("can't gzip: %v\n", err)
	}
	if err := zw.Close(); err != nil {
		fmt.Printf("can't close zw: %v\n", err)
	}

	// Create http request to container. (using JSON+gzip)
	/*
		NOTE: if sending plain json, this was working
			- req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
	*/

	req, err := http.NewRequest("POST", URL, &buf)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("X-Custom-Header", "CaaF")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	// Get response.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Read response.
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response:", string(body))

		// TODO: this currently isn't working - API return value will need to
		// be changed depending on whether this unmarshal is attempted.
		// resultMap := make(map[string]map[string]float64)
		//json.Unmarshal(body, &resultMap)
		//fmt.Println(resultMap)
	} else {
		fmt.Println("response Statuscode:", resp.StatusCode)
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
	}

	// ----------------- stop container

	// Stop the container.
	stopContainerByID(contID, cli)

	// Remove the container.
	removeContainerByID(contID, cli)

	// Delete the image.
	deleteImageByTag(imgTag, images, cli)

	// Prune container.
	// jack := make(map[string]string)
	// jack["slippery"] = "fish"
	cPruneFilter := filters.NewArgs()
	//cPruneFilter.Add("label", "slippery=fish")
	//cPruneFilter.Add("dangling", "true")
	_, err = cli.ContainersPrune(context.Background(), cPruneFilter)
	if err != nil {
		fmt.Printf("Error prune container: %v\n", err)
	}

	// Prune image.
	_, err = cli.ImagesPrune(context.Background(), cPruneFilter)
	if err != nil {
		fmt.Printf("Error prune container: %v\n", err)
	}

}

func createMap(dPath string) (map[string][]string, error) {
	fileMap := make(map[string][]string)

	// Create map of filename to tokenized content.
	dFiles, _ := ioutil.ReadDir(dPath)
	for _, f := range dFiles {
		b, err := ioutil.ReadFile(dPath + f.Name())
		if err != nil {
			fmt.Println(err)
		}

		// Convert bytes to string, then use 3rd party to tokenize.
		fileMap[f.Name()] = tokenize.TextToWords(string(b))
	}

	return fileMap, nil
}