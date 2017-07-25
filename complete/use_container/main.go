package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	"github.com/jhoonb/archivex"
	"github.com/spf13/viper"
)

// Config contains all configuration information -- this is a first draft (jack)
type Config struct {
	pathToDockerfile string
	endPointName     string
	inputPath        string
	hostIP           string
	hostPort         string
	containerName    string
	userName         string
	imgName          string
	imgHandle        string
	tarDir           string
	imgTag           string
}

// TODO: need to rework the input methodology
// -- I think the input "helper" should live it its own file and the input

// config necessary components are;
// - model
//		-- the input type needs to be standardized (json?)
// - Dockerfile
// 		- builds container
// - required files/fixtures
//		-- for example, in cosineSimilarity stopwords+punctuation
// - API (main.go)
// 		- This will likely become the main `dataduit` wrapper
//			-- calls main function
//			-- builds+starts+runs+stops+removes container
// 				-- sets up ports
func createConfig(configPath string) (Config, error) {
	var c Config

	// Set the path for the config file
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	c.pathToDockerfile = viper.Get("model.ddDir").(string)
	c.endPointName = viper.Get("network.host.endpoint").(string)
	c.inputPath = viper.Get("input.file.path").(string)
	c.hostIP = viper.Get("network.host.ip").(string)
	c.hostPort = viper.Get("network.host.port").(string)
	c.containerName = viper.Get("container.name").(string)
	c.userName = viper.Get("container.image.user").(string)
	c.imgName = viper.Get("container.image.img").(string)
	c.imgHandle = c.userName + "/" + c.imgName
	c.tarDir = viper.Get("tar.dir").(string)
	c.imgTag = c.imgHandle + ":latest"

	return c, nil
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

// main creates and uses the container.
func main() {

	configFilePath := "./sudoku_config.yml"
	c, err := createConfig(configFilePath)
	if err != nil {
		fmt.Printf("Error creating config: %v\n", err)
	}

	// TODO: is there a better way to handle this? There likely is, the latest
	// tag may be needed for grabbing the image, but I don't think it's needed
	// for creating the image.
	//imgTag := imgHandle + ":latest"

	// TODO: this could be sent to a temp directory and should also maybe be
	// deleted after use.
	//pathToCreatedTarDir := "./archive/archive"

	// Create tar.
	tarPath, err := createTar(c.tarDir, c.pathToDockerfile)
	if err != nil {
		fmt.Printf("Unable to create tar: %v", err)
	}

	// Create docker cli environment.
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// TODO: look into this variable. Though it is being used "correctly", I'm
	// not sure of its details.
	cntx := context.Background()

	// Build image from the tar file.
	buildImageFromTar(cntx, tarPath, c.imgHandle, cli)

	// TODO: see if I can do this without the loop.
	images, err := cli.ImageList(cntx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	// Build container and get container id.
	contID, err := buildContainerFromImage(cntx, c.imgTag, c.hostIP, c.hostPort, c.containerName, images, cli)
	if err != nil {
		fmt.Printf("error > build container: %v\n", err)
	}

	// Start the container
	startContainerByID(cntx, contID, cli)

	// ----------------------------------- use container

	// URL endpoints.
	URL := "http://" + c.hostIP + ":" + c.hostPort + "/" + c.endPointName
	//URL := "http://" + c.hostIP + ":" + c.hostPort + "/"

	// Create map from input directory.
	// fileMap, err := createMap(inputDir)
	fileMap, err := createSudokuInput(c.inputPath)
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
		NOTE: if sending plain json, this works;
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
	if resp.StatusCode == http.StatusOK {
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
	stopContainerByID(cntx, contID, cli)

	// Remove the container.
	removeContainerByID(cntx, contID, cli)

	// Delete the image.
	deleteImageByTag(cntx, c.imgTag, images, cli)

	// TODO: I'm really unsure how to handle this pruning situation
	// I don't want to prune any images not created by dataduit.
	// The other problem here is that there is the 'multistage' building
	// occuring, so I'm not sure how to make a label for the build stage image.
	//timeString := "30s"
	// Prune containers with the created label.
	cPruneFilter := filters.NewArgs()
	//cPruneFilter.Add("label", "slippery:fish")
	//cPruneFilter.Add("until", timeString)
	_, err = cli.ContainersPrune(cntx, cPruneFilter)
	if err != nil {
		fmt.Printf("Error prune container: %v\n", err)
	}

	// Prune image.
	iPruneFilter := filters.NewArgs()
	//iPruneFilter.Add("until", timeString)
	_, err = cli.ImagesPrune(cntx, iPruneFilter)
	if err != nil {
		fmt.Printf("Error prune image: %v\n", err)
	}

}

// ------------ Helper function related to loading data into dd_container fmt
// TODO: this functionality is dd_container specific.
// createMap is a helper that accepts a path to a directory and creates the
// input data for the model.
// NOTE: this may/may not be included in functionality.  It will likely fall on
// the user to create the specified input data.
// func createMap(dPath string) (map[string][]string, error) {
// 	fileMap := make(map[string][]string)

// 	// Create map of filename to tokenized content.
// 	dFiles, _ := ioutil.ReadDir(dPath)
// 	for _, f := range dFiles {
// 		b, err := ioutil.ReadFile(dPath + f.Name())
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		// Convert bytes to string, then use 3rd party to tokenize.
// 		fileMap[f.Name()] = tokenize.TextToWords(string(b))
// 	}

// 	return fileMap, nil
// }

// -------------------------------------------

// crossIndex 'crosses' two strings such that the two individual values from
// each string join together to create a new value.  For example, if string one
// is "ABC" and string two is "123", the resulting return value will be;
// ["A1","A2","A3","B1","B2","B3","C1","C2","C3"].
func crossIndex(A string, N string) []string {
	var ks []string
	for _, a := range A {
		for _, n := range N {
			ks = append(ks, (string(a) + string(n)))
		}
	}
	return ks
}

func createSudokuInput(fPath string) (map[string][]string, error) {
	sudokuMap := make(map[string][]string)

	data, err := ioutil.ReadFile(fPath)
	if err != nil {
		return sudokuMap, err
	}

	// Global board information.  The Sudoku board is assumed to be a standard
	// 9x9 (A-I)x(1-9) grid -- where the first index (upper left) would be `A1`
	// and the last index (lower right) would be `I9`.
	rows := "ABCDEFGHI"
	cols := "123456789"
	inds := crossIndex(rows, cols)

	// Convert the string representing the board into a grid(map) that maps a
	// key (index) to the values (label for the box, or possible label for the
	// box). for instance, if we know A1=7, map['A1'] = '7', but if the given
	// index is empty (B2, as an example), the corresponding value would be
	// '123456789' (map['B2'] = '123456789') NOTE: i acts as an increment for
	// every target character found.
	i := 0
	for _, c := range data {
		switch string(c) {
		case "_":
			sudokuMap[inds[i]] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
			i++
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			sudokuMap[inds[i]] = []string{string(c)}
			i++
		case "\n", " ", "\r":
			continue
		default:
			return sudokuMap, fmt.Errorf("unexpected value (%v) in Sudoku input", c)
		}
	}

	return sudokuMap, nil
}
