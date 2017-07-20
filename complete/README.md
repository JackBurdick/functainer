## Experimental directory to call CaaF

## Use
To run;
1. Make sure docker is running
2. Navigate to `./use_container/`
3. `go run main.go`

### image information
```
REPOSITORY                  TAG                 IMAGE ID            CREATED              SIZE
jackburdick/cosineexp       latest              a21e1999784a        About a minute ago   10.1MB
```

### Key components
```golang

// input information
inputDir := "./input/"

// URL endpoint to container function
URL := "http://localhost:8080/cosineSim"


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