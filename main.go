package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	info, err := os.Stdin.Stat()
    if err != nil {
        panic(err)
    }
	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("Pipe something into this command")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	var jsonString string
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		jsonString += line
	}
	configPlan := ReadJson(jsonString)
	plan := GeneratePlan(configPlan)
	fmt.Println(ConvertIntoPrettyJSON(plan))
}