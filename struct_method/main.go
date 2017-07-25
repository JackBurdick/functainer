package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jdkato/prose/tokenize"
	"github.com/jhoonb/archivex"
	"github.com/spf13/viper"
)

// Config contains all configuration information -- this is a first draft (jack)
type config struct {
	pathToDockerfile string
	endPointName     string
	hostIP           string
	hostPort         string
	containerName    string
	userName         string
	imgName          string
	imgHandle        string
	tarDir           string
	imgTag           string
	//inputPath        string
}

type created struct {
	contID  string
	tarPath string
	url     string
}

type ddContainer struct {
	ddConfig  config
	ddCreated created
	cntx      context.Context
	cli       *client.Client
}

// set configuration for the container
func (dd *ddContainer) configDD(configPath string) error {

	var c config

	// Set the path for the config file
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	var ok bool
	c.pathToDockerfile, ok = viper.Get("model.ddDir").(string)
	if !ok {
		fmt.Printf("error retriving pathToDockerfile from config\n")
	}
	c.endPointName, ok = viper.Get("network.host.endpoint").(string)
	if !ok {
		fmt.Printf("error retriving endPointName from config\n")
	}
	// c.inputPath, ok = viper.Get("input.file.path").(string)
	// if !ok {
	// 	fmt.Printf("error retriving inputPath from config\n")
	// }
	c.hostIP, ok = viper.Get("network.host.ip").(string)
	if !ok {
		fmt.Printf("error retriving hostIP from config\n")
	}
	c.hostPort, ok = viper.Get("network.host.port").(string)
	if !ok {
		fmt.Printf("error retriving hostPort from config\n")
	}
	c.containerName, ok = viper.Get("container.name").(string)
	if !ok {
		fmt.Printf("error retriving containerName from config\n")
	}
	c.userName, ok = viper.Get("container.image.user").(string)
	if !ok {
		fmt.Printf("error retriving userName from config\n")
	}
	c.imgName, ok = viper.Get("container.image.img").(string)
	if !ok {
		fmt.Printf("error retriving imgName from config\n")
	}
	c.imgHandle = c.userName + "/" + c.imgName
	c.tarDir, ok = viper.Get("tar.dir").(string)
	if !ok {
		fmt.Printf("error retriving tarDir from config\n")
	}
	c.imgTag = c.imgHandle + ":latest"

	dd.ddConfig = c

	return nil
}

// build, start, run container
func (dd *ddContainer) startDD() error {
	// Create tar of all container related files.
	tarPath, err := createTar(dd.ddConfig.tarDir, dd.ddConfig.pathToDockerfile)
	if err != nil {
		fmt.Printf("Unable to create tar: %v", err)
	}
	dd.ddCreated.tarPath = tarPath

	// Create docker cli environment.
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	dd.cli = cli

	// TODO: look into this variable. Though it is being used "correctly", I'm
	// not sure of its details.
	cntx := context.Background()
	dd.cntx = cntx

	// Build image from the tar file.
	buildImageFromTar(cntx, dd.ddCreated.tarPath, dd.ddConfig.imgHandle, dd.cli)

	// TODO: see if I can do this without the loop.
	images, err := dd.cli.ImageList(dd.cntx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	// Build container and get container id.
	contID, err := buildContainerFromImage(dd.cntx, dd.ddConfig.imgTag, dd.ddConfig.hostIP, dd.ddConfig.hostPort, dd.ddConfig.containerName, images, dd.cli)
	if err != nil {
		fmt.Printf("error > build container: %v\n", err)
	}
	dd.ddCreated.contID = contID

	// Start the container
	startContainerByID(dd.cntx, dd.ddCreated.contID, dd.cli)

	// set API endpoint.
	dd.ddCreated.url = "http://" + dd.ddConfig.hostIP + ":" + dd.ddConfig.hostPort + "/" + dd.ddConfig.endPointName
	//dd.ddCreated.url := "http://" + dd.ddConfig.hostIP + ":" + dd.ddConfig.hostPort + "/"

	return nil
}

// run input through container
func (dd *ddContainer) useDD(inputPath string) (string, error) {
	var result string
	// Create map from input directory.
	// TODO: handle input information (file vs dir)
	// TODO: handle input preprocessing function.
	fileMap, err := createMap(inputPath)
	// fileMap, err := createSudokuInput(c.inputPath)
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

	// Create new request.
	req, err := http.NewRequest("POST", dd.ddCreated.url, &buf)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("X-Custom-Header", "CaaF")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	// Send Request and get response.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Read response.
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		result = string(body)

		// TODO: unmarshal to JSON currently isn't working - API return value
		// will need to be changed depending on whether this unmarshal is attempted.
		// resultMap := make(map[string]map[string]float64)
		//json.Unmarshal(body, &resultMap)
		//fmt.Println(resultMap)
	} else {
		// TODO: create return error
		fmt.Println("response Statuscode:", resp.StatusCode)
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		result = "Nothing to see here"
	}

	return result, nil
}

// stop, delete container and image
func (dd *ddContainer) endDD() error {

	// Stop the container.
	stopContainerByID(dd.cntx, dd.ddCreated.contID, dd.cli)

	// Remove the container.
	removeContainerByID(dd.cntx, dd.ddCreated.contID, dd.cli)

	// TODO: see if I can do this without the loop.
	images, err := dd.cli.ImageList(dd.cntx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	// Delete the image.
	deleteImageByTag(dd.cntx, dd.ddConfig.imgTag, images, dd.cli)

	return nil
}

// config, start, stop, end
func (dd *ddContainer) completeDD(inputPath string) (string, error) {
	err := dd.startDD()
	if err != nil {
		fmt.Printf("Error starting container: %v\n", err)
	}

	res, err := dd.useDD(inputPath)
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}
	//fmt.Printf("Response: %v", res)

	dd.endDD()
	return res, nil
}

// createTar creates a tar of the Dockerfile directory.
func createTar(pathToCreatedTarDir string, pathToDockerfile string) (string, error) {
	tar := new(archivex.TarFile)
	tar.Create(pathToCreatedTarDir)
	tar.AddAll(pathToDockerfile, false)
	tar.Close()
	return pathToCreatedTarDir + ".tar", nil
}

// buildImageFromTar creates the tar of the Dockerfile and directory.
// NOTE: the required files need to be placed in the root of the directory where
// the Dockerfile is located `pathToDockerfile`.  The current implementation
// does not support the addition of directories.
func buildImageFromTar(cntx context.Context, tarPath string, imgHandle string, cli *client.Client) {
	dockerBuildContext, err := os.Open(tarPath)
	defer dockerBuildContext.Close()
	buildOptions := types.ImageBuildOptions{Tags: []string{imgHandle}}
	buildResponse, err := cli.ImageBuild(cntx, dockerBuildContext, buildOptions)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	// NOTE: This needs to be here to ensure the image is built before we
	// start/run it. There is likely a more elegant way to handle this.
	_, err = ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

// buildContainerFromImage builds the container from the image and returns the
// created container id.
func buildContainerFromImage(cntx context.Context, imgTag string, hostIP string, hostPort string, containerName string, images []types.ImageSummary, cli *client.Client) (string, error) {
	var contID string
	for _, image := range images {

		// Select specified image by the repo tag.
		if strings.Join(image.RepoTags, "") == imgTag {

			// Create the container from the image.
			// TODO: Devise a progamatic way of producting a container name.
			// I'm not even sure if the container name is assigned right now.
			// TODO: Determine how exposedPort and port Bindings are different.
			exposedPort := map[nat.Port]struct{}{"8080/tcp": {}}
			jack := make(map[string]string)
			jack["slippery"] = "fish"
			configOptions := container.Config{Image: strings.Join(image.RepoTags, ""), ExposedPorts: exposedPort, Labels: jack}
			networkConfig := network.NetworkingConfig{}
			portBindings := map[nat.Port][]nat.PortBinding{
				"8080/tcp": {{HostIP: hostIP, HostPort: hostPort}}}
			hostConfig := container.HostConfig{
				PublishAllPorts: true,
				PortBindings:    portBindings,
			}

			createResponse, err := cli.ContainerCreate(cntx, &configOptions, &hostConfig, &networkConfig, containerName)
			if err != nil {
				fmt.Println(err)
			}
			contID = createResponse.ID
		}
	}
	return contID, nil
}

// startContainerByID starts the container by the specified container id.
func startContainerByID(cntx context.Context, contID string, cli *client.Client) {
	if contID != "" {
		err := cli.ContainerStart(cntx, contID, types.ContainerStartOptions{})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Error: no container created")
	}
}

// stopContainerByID stops the container by the specified container id.
func stopContainerByID(cntx context.Context, contID string, cli *client.Client) {

	// TODO: I have a stop time here, when the stoptime was nil, the processes
	// took noticeably long.. I need to investigate what, if any, problems this
	// causes.
	stopTime := time.Duration(100) * time.Millisecond
	err := cli.ContainerStop(cntx, contID, &stopTime)
	if err != nil {
		fmt.Println("ERROR: can't stop container")
	}
}

// removeContainerByID removes the container by the specified container id.
func removeContainerByID(cntx context.Context, contID string, cli *client.Client) {

	// TODO: Weigh the advantages of using the `Force: true` flag here
	err := cli.ContainerRemove(cntx, contID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		fmt.Println("ERROR: can't remove container")
	}
}

// buildContainerFromImage accepts a list of the current images and the
// specified image tag. If the image is present, the image is removed.
func deleteImageByTag(cntx context.Context, imgTag string, images []types.ImageSummary, cli *client.Client) {

	// TODO: Errors should be handled/returned.
	// TODO: Can the image loop be removed?
	for _, image := range images {
		if strings.Join(image.RepoTags, "") == imgTag {
			imgID := strings.TrimLeft(image.ID, "sha256")
			imgID = strings.TrimLeft(imgID, ":")

			// TODO: Weigh the advantages of using the `Force: true` flag here
			_, err := cli.ImageRemove(cntx, imgID, types.ImageRemoveOptions{Force: true})
			if err != nil {
				fmt.Printf("ERROR: image %v not deleted\n", imgID)
			}
		}
	}
}

func main() {
	configPath := "./cosine_config.yml"
	inputPath := "../complete/use_container/input/"

	// ----------------------- component call
	var dd ddContainer

	err := dd.configDD(configPath)
	if err != nil {
		fmt.Printf("Error with dd conf: %v\n", err)
	}

	err = dd.startDD()
	if err != nil {
		fmt.Printf("Error starting container: %v\n", err)
	}

	res, err := dd.useDD(inputPath)
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}
	fmt.Printf("Response A: %v\n", res)
	fmt.Printf("\nDONE with multistage call\n\n")

	dd.endDD()

	// --------------------------- single call
	var jj ddContainer

	// initialize
	err = jj.configDD(configPath)
	if err != nil {
		fmt.Printf("Error with jj conf: %v\n", err)
	}

	// use
	resB, err := jj.completeDD(inputPath)
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}

	// examine results
	fmt.Printf("Response B: %v", resB)

	fmt.Println("DONE")

}

// createMap is a helper that accepts a path to a directory and creates the
// input data for the model.
// NOTE: this may/may not be included in functionality.  It will likely fall on
// the user to create the specified input data.
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
