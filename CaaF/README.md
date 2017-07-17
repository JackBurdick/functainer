#CaaF: **C**ontainer **a**s **a** **F**unction


### Ideal case
```golang
...

// create CaaF
ddCosine := dataduit.build(xxxxxxxxx_config.yaml)

inDirectory := "./input_directory"

// convert input into specified format (JSON?)
data, err := magicHelper(inDirectory)
if err != nil {
    fmt.Fprintf(w, "Unable to convert input: %v", err)
}

// use CaaF
/*
    *. start/run container
    1. http req w/data as body to container on specified port
    2. container (`ddCosine`) processes results
    3. Results are returned

*/
results, err := ddCosine(data)
if err != nil {
    fmt.Fprintf(w, "Unable to obtain results: %v", err)
}

```

### Potential problems
- How big can the request body be?
- How do we manage ports?
- Is this even the "right" way to approach this?