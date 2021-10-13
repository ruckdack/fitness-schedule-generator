package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func ReadJson(fileLocation string) *ConfigPlan {
	jsonFile, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var configPlan ConfigPlan
	byteJson, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	json.Unmarshal(byteJson, &configPlan)
	// read weekdays into internal time format
	configPlan.Weekdays = readWeekdays(&byteJson)
	err = checkConfig(&configPlan)
	if err == nil {
		return &configPlan
	}
	log.Fatal(err.Error())
	return nil
}

func readWeekdays(byteJson *[]byte) map[time.Weekday]string {
	type PlanConfigWeekdaysOnly struct {
		Weekdays map[string]string `json:"weekdays"`
	}
	var configPlanWeekdaysOnly PlanConfigWeekdaysOnly
	json.Unmarshal(*byteJson, &configPlanWeekdaysOnly)
	weekdayMap := make(map[time.Weekday]string)
	for key, value := range configPlanWeekdaysOnly.Weekdays {
		_, found := Find(getWeekdays(DAYS_OF_WEEK), key)
		if !found {
			log.Fatal("weekdays are not valid")
		}
		weekdayMap[DAYS_OF_WEEK[key]] = value
	}
	return weekdayMap
}

func checkConfig(configPlan *ConfigPlan) error {
	// check if startDate is a valid date
	_, err := time.Parse(TIME_LAYOUT, configPlan.StartDate)
	if err != nil {
		return nil
	}
	// check if weekdays use valid splits
	splitNames := make([]string, len(configPlan.Splits))
	for idx, split := range configPlan.Splits {
		splitNames[idx] = split.Name
	}
	for weekday, val := range configPlan.Weekdays {
		_, found := Find(splitNames, val)
		if !found {
			return errors.New(val + " on " + weekday.String() + " is not defined")
		}
	}
	// check if exercise identifiers in splits are defined
	availableExercises := make([]string, len(configPlan.Exercises))
	for idx, exercise := range configPlan.Exercises {
		availableExercises[idx] = exercise.Name
	}
	for _, split := range configPlan.Splits {
		for _, variations := range split.Executions {
			if len(variations) > 2 {
				return errors.New("no more than 2 variations allowed")
			}
			for _, exercise := range variations {
				_, found := Find(availableExercises, exercise)
				if !found {
					return errors.New(exercise + " is not defined")
				}
			}
		}
	}
	// check if muscle identifiers in exercises are defined
	availableMuscles := make([]string, len(configPlan.Muscles))
	for idx, muscle := range configPlan.Muscles {
		availableMuscles[idx] = muscle.Name
	}
	for _, exercise := range configPlan.Exercises {
		_, found := Find(availableMuscles, exercise.Target)
		if !found {
			return errors.New(exercise.Target + " is not defined")
		}
	}
	return nil
}

func getWeekdays(m map[string]time.Weekday) []string {
	keys := make([]string, len(m))
	idx := 0
	for key := range m {
		keys[idx] = key
		idx++
	}
	return keys
}