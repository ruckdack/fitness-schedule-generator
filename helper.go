package main

import (
	"encoding/json"
	"log"
	"time"
)

func FindString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func FindWeekday(slice []time.Weekday, val time.Weekday) (int, bool) {
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

func ReverseInt(arr []int) []int {
	newArr := make([]int, len(arr))
    for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		newArr[i], newArr[j] = arr[j], arr[i]
	}
	return newArr
}

func NextDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}