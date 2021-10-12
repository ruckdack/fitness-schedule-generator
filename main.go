package main

import (
	"os"
)

func main() {
	configFileLocation := os.Args[1]
	configPlan := ReadJson(configFileLocation)
	plan := GeneratePlan(configPlan)
	PrintPretty(plan)
}