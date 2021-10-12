package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func PrintPretty(emp interface{}) {
	//MarshalIndent
	empJSON, err := json.MarshalIndent(emp, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(empJSON))
}