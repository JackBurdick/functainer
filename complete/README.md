# Working Directory

## Known issues
- All containers are pruned
    - target filter is not working as expected
- All images are pruned
    - target filter is not working as expected
- Unable to marshal the response from plain text -> JSON -> MAP

### Look into
- gzip return data?
- only use gzip if above a threshold?
- archivex does not support tar (sub)directories
    - can the standard lib be used instead?
- api functionality
    - move api to wrapper?
        - stuct/method calls?

## Use
To run;
1. Make sure docker is running
2. Navigate to `./use_container/`
3. `go run main.go`

Working Directory
- `./Complete/`
    - `./dd_cosineSim/`
        - Holds the main functionality / is meant to be modular
        - to use, you must specify the path in `./use_container/main.go`
    - `./use_container/`
        - Calls the container function (`./<dd_container>/`)
Experimental
- Experimental scripts and files are held in this directory
    - `./CaaF/`
    - `docker_client`

### image information
```
REPOSITORY                  TAG                 IMAGE ID            CREATED              SIZE
jackburdick/cosineexp       latest              a21e1999784a        About a minute ago   10.1MB
```

### Key components
```golang

// input information

// URL endpoint to container function

// build JSON data
//....

// gzip JSON data
// ....

// Build and execute request
// - the POST request sends the gzipped json data
// ...
req, err := http.NewRequest("POST", URL, &buf)
req.Header.Set("Content-Encoding", "gzip")
resp, err := client.Do(req)


// display body
body, _ := ioutil.ReadAll(resp.Body)

```