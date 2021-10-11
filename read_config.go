package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const timeLayout = "2006-01-02"
var daysOfWeek = map[string]time.Weekday {
	"mon": time.Monday,
	"tue": time.Tuesday,
	"wed": time.Wednesday,
	"thu": time.Thursday,
	"fri": time.Friday,
	"sat": time.Saturday,
	"sun": time.Sunday,
}

func ReadJson() *PlanConfig {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var planConfig PlanConfig
	byteJson, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	json.Unmarshal(byteJson, &planConfig)
	// read weekdays into internal time format
	planConfig.Weekdays = readWeekdays(&byteJson)
	if isConfigValid(&planConfig) {
		return &planConfig
	}
	log.Fatal("weekdays or start date is not valid")
	return nil
}

func readWeekdays(byteJson *[]byte) map[time.Weekday]string {
	type PlanConfigWeekdaysOnly struct {
		Weekdays map[string]string `json:"weekdays"`
	}
	var planConfigWeekdaysOnly PlanConfigWeekdaysOnly
	json.Unmarshal(*byteJson, &planConfigWeekdaysOnly)
	weekdayMap := make(map[time.Weekday]string)
	for key, value := range planConfigWeekdaysOnly.Weekdays {
		_, found := Find(getWeekdays(daysOfWeek), key)
		if !found {
			log.Fatal("weekdays are not valid")
		}
		weekdayMap[daysOfWeek[key]] = value
	}
	return weekdayMap
}

func isConfigValid(planConfig *PlanConfig) bool {
	// check if startDate is a valid date
	_, err := time.Parse(timeLayout, planConfig.StartDate)
	if err != nil {
		return false
	}
	// check if weekdays use valid splits
	splitNames := make([]string, len(planConfig.Splits))
	for idx, split := range planConfig.Splits {
		splitNames[idx] = split.Name
	}
	for _, val := range planConfig.Weekdays {
		_, found := Find(splitNames, val)
		if !found {
			return false
		}
	}
	return true
}

func getWeekdays(m map[string]time.Weekday) []string {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	return keys
}