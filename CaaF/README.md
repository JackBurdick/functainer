# CaaF
**C**ontainer **a**s **a** **F**unction

#### Warnings:
- all docker images are pruned
- all docker containers are pruned
- files are written (`.tar` in `./archive/archive/`)

## Use
To run;
1. Make sure docker is running
2. Navigate to `./use_container/`
3. `go run main.go`

### Known issues
- All containers are pruned
    - target filter is not working as expected
- All images are pruned
    - target filter is not working as expected
- Unable to marshal the response from plain text -> JSON -> MAP

#### TODO/Look into
- Eventually, the functionality needs to be packaged up so that it can be imported
- "fair weather functionality" -- error handling is non-existent, this would fair miserably if/when something goes wrong.
- The preprocessing function situation is a sore spot
    - The basic issue is that users will need to write a function to get their raw input into the desired input (which could be indicated in the config/documentation) -- this is fine. The issue arrises when desiding how to handle "common" cases (i.e. when using cosine similarity, it'd be nice to just specify a directory to run CS on all the included documents).. Thinking on this some more, it may not be an issue. Examples could just be shown in the documentation for copy/pasting according to the individuals needs.
- gzip return data?
- only use gzip if above a threshold?
- archivex does not support tar (sub)directories
    - can the standard lib be used instead?


Working Directory
- `./container/`
    - holds the main functionality / is meant to be modular
- `./use_container/`
    - `./archive/archive.tar`
        - .tar (created/overwritten when building a container)
    - `./input/`
        - example input files
    - `./config/`
        - holds config files
    - `./main.go`
        - **entry point - all functionality lives here**

### approximate image information
```
REPOSITORY                  TAG                 IMAGE ID            CREATED              SIZE
jackburdick/cosineexp       latest              a21e1999784a        About a minute ago   10.1MB
```