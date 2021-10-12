package main

import "fmt"

func main() {
	configPlan := ReadJson("config.json")
	plan := GeneratePlan(configPlan)
	fmt.Println(ConvertIntoPrettyJSON(plan))
}