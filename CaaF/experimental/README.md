## Experimental CaaF

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