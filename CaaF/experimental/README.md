## Experimental CaaF

**NOTE**
- The ports are a bit strange.. currently calling `8080` from "api", this doesn't feel right.. I think it should really be `8000`
- Everything is currently gziped -- I don't know the limitations of this.
- A config file will be needed

## Running
1. build docker image
    - `docker image build -t jackburdick/cosineexp .`
2. Run the docker image
    - `docker run -d --rm -p 8000:8080 jackburdick/cosineexp`
3. Navigate to main file
    - target: `dataduit\CaaF\experimental\container`
4. Run `main.go` for testing
    - `go run main.go`



## Resources
- `Dockerfile`
    - [multi-stage-docker-builds-for-creating-tiny-go-images](https://medium.com/travis-on-docker/multi-stage-docker-builds-for-creating-tiny-go-images-e0e1867efe5a)
- `fixtures`
    - `/en_stopwords.json` [link](https://github.com/6/stopwords-json)
    - `/en_punctuation.json` custom: python [string.punctuation](https://docs.python.org/2/library/string.html#string.punctuation) + extras

### Goal
1. Build function
2. Write main API
3. Containerize
4. Call From go
    - HTTP

### LOOK INTO
1. is there a way to build/start/stop containers from w/in go?