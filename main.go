package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	jsonString, err := ReadPipe()
	if err != nil {
		log.Fatal(err.Error())
	}
	configPlan := ReadJson(jsonString)
	plan := GeneratePlan(configPlan)
	fmt.Println(ConvertIntoPrettyJSON(plan))
}

func ReadPipe() (string, error) {
	info, err := os.Stdin.Stat()
    if err != nil {
        panic(err)
    }
	if info.Mode()&os.ModeCharDevice != 0 {
		return "", errors.New("pipe something into this command")
	}
	reader := bufio.NewReader(os.Stdin)
	var str string
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		str += line
	}
	return str, nil
}