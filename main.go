package main

import (
	"fmt"
	"log"
)

func main() {
	planConfig := *ReadJson()
	log.Println("read json config successfully")
	fmt.Printf("%+v\n", planConfig)
}