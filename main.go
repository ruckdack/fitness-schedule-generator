package main

import (
	"fmt"
	"os"
)

func main() {
	configFileLocation := os.Args[1]
	planConfig := *ReadJson(configFileLocation)
	fmt.Printf("%+v\n", planConfig)
}