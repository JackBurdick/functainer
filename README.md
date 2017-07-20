# Dataduit

## Included
Working Directory
- `./Complete/`
    - `./container/`
        - Holds the main functionality
    - `./use_container/`
        - Calls the container function (`./container/`)
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