## Experimental directory to call CaaF

## Use
Please see parent directory


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

// Build and execute requst
// ...
req, err := http.NewRequest("POST", URL, &buf)
req.Header.Set("Content-Encoding", "gzip")
resp, err := client.Do(req)


// display body
body, _ := ioutil.ReadAll(resp.Body)
fmt.Println("response Body:", string(body))

```