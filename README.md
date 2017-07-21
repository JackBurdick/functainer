# Dataduit

## Included
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

### Note
- latest docker build may be needed
    - information here: https://github.com/moby/moby/issues/33946

### Links
- [lucid chart](https://www.lucidchart.com/invitations/accept/e348323e-11f7-45a9-b830-77764936fdb8)

### Useful Resources
- I'll look into a dropbox/pocket/drive folder for dumping ideas

### In Progress/Future
1. Config file
    - specifications
        - managing ports/"pipeline"
    - building
2. dataduit wrapper
    - [dataduit container (API)]<-->[function(s)]
        - main "API"
            - traffic routed from program, through here, to functions

### Stretch considerations
1. Function scaling?
    - scale up testing functions when applicable?
2. Helpers for sending "lots" of data over the line -- multiple gzips
    - theory
        - shard set
        - sent in pieces
        - rebuild
    

### Currently have
- dataduit.com
- dataduit.io
- datadu.it
- @dataduit