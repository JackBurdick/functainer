package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	//fileMap := make(map[string][]string)

	data, err := ioutil.ReadFile("./test.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print("data:  ", string(data))
	var slice []string
	err = json.Unmarshal(data, &slice)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("slice: %q\n", slice)

	// decoder := json.NewDecoder(jsonData)
	// var JSONdata []byte
	// err := decoder.Decode(&JSONdata)

	// jack, err := json.Marshal(jsonData)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(jsonData)

	// var c []byte
	// json.Unmarshal(jsonData, c)

	//var x map[string]interface{}
	// var x map[string][]string
	// //var x string
	// json.Unmarshal(jsonData, &x)
	// err = json.NewDecoder().Decode(&fileMap)
	// if err != nil {
	// 	http.Error(w, err.Error(), 400)
	// 	return
	// }

	//fmt.Printf("%v\n", x)
}
