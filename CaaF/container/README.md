# container
Containers, to be used as functions, live in this directory.

## TODO:
There should be a README in each container that;
- specifies the input
- explains how it works
- shows the output
- includes an example usage


## Basic container necessities
All containers are currently go based, but I am working on python
- `main.go`
    - API/http handler - calls `<function>`
- `<function>.go`
    - main functionality
- `Dockerfile`
    - specifies the container creation


## Current containers
- `./cosineSim/`
    - performs cosine similarity on files
- `./sudokuSolver/`
    - solves sudoku puzzles