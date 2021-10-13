package main

import (
	"encoding/json"
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

// spreads towards beginning
func SpreadStart(amount int, parts int, idx int) int {
	partsArr := make([]int, parts)
	for i := 0; i < amount; i++ {
		partsArr[i % parts]++
	}
	return partsArr[idx]
}

func SpreadEnd(amount int, parts int, idx int) int {
	return SpreadStart(amount, parts, parts - 1 - idx)
}

func ConvertIntoPrettyJSON(emp interface{}) string {
	empJSON, err := json.MarshalIndent(emp, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(empJSON)
}

func Reverse(arr []interface{}) []interface{} {
	newArr := make([]interface{}, len(arr))
    for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		newArr[i], newArr[j] = arr[j], arr[i]
	}
	return newArr
}