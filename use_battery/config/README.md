# config
config files that specify container usage

## Future advancements:
- a flag (inside .yml or elsewhere) for local/default or specified IP would be nice

## Essential Components
```yml
---
# path to the container creation file (relative to this file)
model:
  ddDir: "relative/path/to/container/"
  
# created container information
# container.name is the name of the created container.
# imgHandle is the name of the created image( userName + "/" + imgName )
container:
  name: "<container_name>"
  image:
    user: "jackburdick"
    img: "automated"

# network information
# host.ip is the specified IP to host the container.
# host.port is the port that is exposed to the user/can be called from the API.
# NOTE: the `endpoint` must match the containers `main.go` api functionality
network:
  host:
    ip: "127.0.0.1"
    port: "8000"
    endpoint: "<functionName>"

tar:
  dir: "./archive/archive"
...
```