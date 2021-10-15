package main

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

func ReadJson(jsonString string) *ConfigPlan {
	var configPlan ConfigPlan
	byteJson := []byte(jsonString)
	err := json.Unmarshal(byteJson, &configPlan)
	if (err != nil) {
		log.Fatal("json structure is not correct")
		return nil
	}
	readWeekdays(&configPlan)
	err = checkConfig(&configPlan)
	if err == nil {
		return &configPlan
	}
	log.Fatal(err.Error())
	return nil
}

func readWeekdays(configPlan *ConfigPlan) {
	configPlan.Weekdays = make(map[time.Weekday]string)
	for key, value := range configPlan.ConfigWeekdays {
		_, found := FindString(getWeekdays(DAYS_OF_WEEK), key)
		if !found {
			log.Fatal("weekdays are not valid")
		}
		configPlan.Weekdays[DAYS_OF_WEEK[key]] = value
	}
}

func checkConfig(configPlan *ConfigPlan) error {
	regex, _ := regexp.Compile(ALLOWED_STRING_REGEX)
	// check if all fields are valid
	if (len(configPlan.Weekdays) == 0 || configPlan.StartDate == "" || len(configPlan.Splits) == 0 || len(configPlan.Exercises) == 0 || len(configPlan.Muscles) == 0) {
		return errors.New("json fields are missing")
	}
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
		_, found := FindString(splitNames, val)
		if !found {
			return errors.New(val + " on " + strings.ToLower(weekday.String()) + " is not defined")
		}
	}
	// check if split string fields are matching allowed string regex
	for _, split := range configPlan.Splits {
		if !regex.MatchString(split.Name) {
			return errors.New("identifier \"" + split.Name + "\" is not allowed")
		}
	}
	// check if muscle string fields are matching allowed string regex
	for _, muscle := range configPlan.Muscles {
		if !regex.MatchString(muscle.Name) {
			return errors.New("identifier \"" + muscle.Name + "\" is not allowed")
		}
	}
	// check if exercise string fields are matching allowed string regex
	availableExercises := make([]string, len(configPlan.Exercises))
	for idx, exercise := range configPlan.Exercises {
		availableExercises[idx] = exercise.Name
	}
	for _, exercise := range availableExercises {
		if !regex.MatchString(exercise) {
			return errors.New("identifier \"" + exercise + "\" is not allowed")
		}
	}
	// check if exercise identifiers in splits are defined
	for _, split := range configPlan.Splits {
		for _, superset := range split.Supersets {
			for _, variations := range superset {
				if len(variations) > 2 {
					return errors.New("in split " + split.Name + ": no more than 2 variations allowed")
				}
				for _, exercise := range variations {
					_, found := FindString(availableExercises, exercise)
					if !found {
						return errors.New(exercise + " is not defined")
					}
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
		_, found := FindString(availableMuscles, exercise.Target)
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