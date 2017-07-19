package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
)

// IMPORTANT: I had to change the following;
/*
	"github.com/moby/moby/api/types"
	"github.com/moby/moby/client"
-- to
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

------------------------------ Unrelated
-- current error message:
{"errorDetail":
		{"message":"invalid from flag value build-env: repository sha256 not found: does not exist or no pull access"},
		"error":"invalid from flag value build-env: repository sha256 not found: does not exist or no pull access"}
*/

func main() {

	// Create tar.
	tar := new(archivex.TarFile)
	tar.Create("./test_arch/archieve")
	tar.AddAll("../../CaaF/Experimental/container/", false)
	tar.Close()

	// Initialize client.
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	dockerBuildContext, err := os.Open("./test_arch/archieve.tar")
	defer dockerBuildContext.Close()
	//defaultHeaders := map[string]string{"User-Agent": ""}

	buildOptions := types.ImageBuildOptions{}

	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	// Output build response information.
	fmt.Printf("buildResponse.OSType: %s \n", buildResponse.OSType)
	response, err := ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Println(string(response))
}
