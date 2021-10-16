package main

import (
	"fmt"
	"log"
)

func main() {
	jsonString, err := ReadPipe()
	if err != nil {
		log.Fatal(err.Error())
	}
	configPlan, err := ReadJson(jsonString)
	if (err != nil) {
		log.Fatal(err)
		return
	}
	plan := GeneratePlan(configPlan)
	fmt.Println(ConvertIntoPrettyJSON(plan))
}