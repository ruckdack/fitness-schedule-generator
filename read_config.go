package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

//const jsonString = `{"weekdays":{"mon":"upper","tue":"lower","wed":"upper","fri":"lower","sat":"upper"},"start_date":"2021-10-15","splits":[{"name":"upper","executions":[["inclinebenchpress(dumbbell)","benchpress"],["highrow(pulley)"],["lateralraise(dumbbell)","lateralraise(machine)"],["latpulldown"],["inclinecurl(dumbbell)","curlmachine"],["tricepsextension","overheadtricepsextension"]]},{"name":"lower","executions":[["iwasfürquads"],["legextension"],["RDL","hipthrust"],["layinglegcurl"],["seatedcalfraise"],["standingcalfraise"],["crunchmachine"]]}],"exercises":[{"name":"inclinebenchpress(dumbbell)","initial_1rm":107,"reps":9,"target":"chest"},{"name":"benchpress","initial_1rm":100,"reps":8,"target":"chest"},{"name":"highrow(pulley)","initial_1rm":100,"reps":8,"target":"lowertraps/rhomboids"},{"name":"lateralraise(dumbbell)","initial_1rm":100,"reps":8,"target":"delts"},{"name":"lateralraise(machine)","initial_1rm":100,"reps":8,"target":"delts"},{"name":"latpulldown","initial_1rm":100,"reps":8,"target":"lat"},{"name":"inclinecurl(dumbbell)","initial_1rm":100,"reps":8,"target":"biceps"},{"name":"curlmachine","initial_1rm":100,"reps":8,"target":"biceps"},{"name":"tricepsextension","initial_1rm":100,"reps":8,"target":"triceps"},{"name":"overheadtricepsextension","initial_1rm":100,"reps":8,"target":"triceps"},{"name":"crunchmachine","initial_1rm":100,"reps":8,"target":"abs"},{"name":"standingcalfraise","initial_1rm":100,"reps":8,"target":"calves"},{"name":"seatedcalfraise","initial_1rm":100,"reps":8,"target":"calves"},{"name":"layinglegcurl","initial_1rm":100,"reps":8,"target":"hamstrings"},{"name":"hipthrust","initial_1rm":100,"reps":8,"target":"glute"},{"name":"RDL","initial_1rm":100,"reps":8,"target":"hamstrings"},{"name":"legextension","initial_1rm":100,"reps":8,"target":"quads"},{"name":"iwasfürquads","initial_1rm":100,"reps":8,"target":"quads"}],"muscles":[{"name":"quads","sets":[12,14,14,16,8]},{"name":"chest","sets":[12,14,14,16,8]},{"name":"hamstrings","sets":[12,14,14,16,8]},{"name":"delts","sets":[12,14,14,16,8]},{"name":"calves","sets":[12,14,14,16,8]},{"name":"glute","sets":[12,14,14,16,8]},{"name":"lowertraps/rhomboids","sets":[12,14,14,16,8]},{"name":"lat","sets":[12,14,14,16,8]},{"name":"delts","sets":[12,14,14,16,8]},{"name":"abs","sets":[12,14,14,16,8]},{"name":"biceps","sets":[12,14,14,16,8]},{"name":"triceps","sets":[12,14,14,16,8]}]}`

func ReadJson(jsonString string) *ConfigPlan {
	// jsonFile, err := os.Open(fileLocation)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil
	// }
	var configPlan ConfigPlan
	// byteJson, _ := ioutil.ReadAll(jsonFile)
	// jsonFile.Close()
	byteJson := []byte(jsonString)
	json.Unmarshal(byteJson, &configPlan)
	// read weekdays into internal time format
	configPlan.Weekdays = readWeekdays(&byteJson)
	err := checkConfig(&configPlan)
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
		_, found := FindString(getWeekdays(DAYS_OF_WEEK), key)
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
		_, found := FindString(splitNames, val)
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
				_, found := FindString(availableExercises, exercise)
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