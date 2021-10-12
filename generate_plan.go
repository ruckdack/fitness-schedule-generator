package main

import (
	"sort"
	"time"
)

func GeneratePlan(configPlan *ConfigPlan) Plan {
	workoutsPerWeek := len(configPlan.Weekdays)
	//frequencyPerExercise := getFrequencyPerExercise(configPlan)
	trainigWeekdays := getTrainingWeekdaysInOrder(configPlan)

	variationCountPerSplit := make([]int, len(configPlan.Splits))
	generateWeek := func(weekIndex int) Week {
		
		generateWorkoutDay := func(workoutInWeek int) WorkoutDay {

			getSplitIndex := func(splitName string) int {
				for idx, split := range configPlan.Splits {
					if split.Name == splitName {
						return idx
					}
				}
				// should not happen if config checker works
				return 0
			}

			splitName := configPlan.Weekdays[trainigWeekdays[workoutInWeek]]
			splitIndex := getSplitIndex(splitName)
			executions := configPlan.Splits[splitIndex].Executions
			exercises := make([]Exercise, len(executions))
			// select correct exercise from variations
			for idx, variations := range executions {
				var name string
				if len(variations) == 2 {
					name = variations[variationCountPerSplit[splitIndex]]
				} else {
					name = variations[0]
				}
				exercises[idx] = Exercise{
					Name: name,
					Weight: 0,
					Reps: 0,
					Sets: 0,
					Rir: 0,
				}
			}
			variationCountPerSplit[splitIndex] = (variationCountPerSplit[splitIndex] + 1) % 2
			return WorkoutDay{
				Date: "bla",
				Weekday: trainigWeekdays[workoutInWeek].String(),
				WeekType: WeekTypes[weekIndex],
				Split: splitName,
				Exercises: exercises,
			}
		}
		week := make(Week, workoutsPerWeek)
		for workoutInWeek := range trainigWeekdays {
			week[workoutInWeek] = generateWorkoutDay(workoutInWeek)
		}
		return week
	}

	var plan Plan // wtf?
	for weekIndex := range plan {
		plan[weekIndex] = generateWeek(weekIndex)
	}
	return plan
}

func getFrequencyPerExercise(configPlan *ConfigPlan) map[string]int {
	m := make(map[string]int)
	for _, split := range configPlan.Weekdays {
		m[split] = m[split] + 1
	}
	return m
}

func getTrainingWeekdaysInOrder(configPlan *ConfigPlan) []time.Weekday {
	intWeekdays := make([]int, len(configPlan.Weekdays))
	idx := 0
	for weekday := range configPlan.Weekdays {
		intWeekdays[idx] = int(weekday)
		idx++
	}
	sort.Ints(intWeekdays)
	weekdays := make([]time.Weekday, len(intWeekdays))
	idx = 0
	for _, intWeekday := range intWeekdays {
		weekdays[idx] = time.Weekday(intWeekday)
		idx++
	}
	return weekdays
}