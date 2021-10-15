package main

import (
	"math"
	"sort"
	"strings"
	"time"
)

func GeneratePlan(configPlan *ConfigPlan) Plan {
	var plan Plan

	for weekIdx := range plan {
		plan[weekIdx] = generateWeek(configPlan, weekIdx)
	}
	addDates(&plan, configPlan)
	return plan
}

func generateWeek(configPlan *ConfigPlan, weekIdx int) Week {
	variationCountPerSplit := make([]int, len(configPlan.Splits))
	workoutsPerWeek := len(configPlan.Weekdays)
	workoutWeekdays := getRotatedTrainingWeekdays(configPlan)
	
	generateWorkoutDay := func(workoutInWeek int) WorkoutDay {

		splitName := configPlan.Weekdays[workoutWeekdays[workoutInWeek]]
		splitIndex := func() int {
			for idx, split := range configPlan.Splits {
				if split.Name == splitName {
					return idx
				}
			}
			// should not happen if config checker works
			return 0
		}()
		configSupersets := configPlan.Splits[splitIndex].Supersets
		supersets := make([]Superset, len(configSupersets))

		// select correct exercise from variations
		for supersetIdx, superset := range configSupersets {
			exercisesInSuperset := make([]Exercise, len(superset))
			for _, variations := range superset {
				name := func() string {
					if len(variations) == 2 {
						return variations[variationCountPerSplit[splitIndex]]
					} 
					return variations[0]
				}()
				target, reps := func() (string, int) {
					for _, exercise := range configPlan.Exercises {
						if exercise.Name == name {
							return exercise.Target, exercise.Reps
						}
					}
					// should not happen if config checker works
					return "", 0
				}()
				weight := math.Pow(1 + ONE_RM_INCREASE, float64(weekIdx % len(WEEK_TYPES)-1)) * configPlan.Exercises[idx].InitialOneRM
				weight /= 1 + float64(reps + RIR_MAPPING[weekIdx]) / 30
				weight = math.Round(4 * weight) / 4
				// sets is added in a second round as we don't know yet which exercises will be selected for the entire week
				exercisesInSuperset[supersetIdx] = Exercise{
					Name: name,
					Target: target,
					Reps: reps,
					Weight: weight,
					Rir: RIR_MAPPING[weekIdx],
				}
			}
			supersets[supersetIdx] = exercisesInSuperset
		}
		variationCountPerSplit[splitIndex] = (variationCountPerSplit[splitIndex] + 1) % 2
		return WorkoutDay{
			Weekday: strings.ToLower(workoutWeekdays[workoutInWeek].String()),
			WeekType: WEEK_TYPES[weekIdx],
			Split: splitName,
			Supersets: supersets,
		}
	}
	week := make(Week, workoutsPerWeek)
	for workoutInWeek := range workoutWeekdays {
		week[workoutInWeek] = generateWorkoutDay(workoutInWeek)
	}
	return addSets(week, weekIdx, configPlan)
}

func addSets(week Week, weekIdx int, configPlan *ConfigPlan) Week {
	frequencyPerMuscle := make(map[string]int)
	for _, workoutDay := range(week) {
		// prevents double adding when a muscle is hit twice in a day
		musclesHit := make(map[string]bool)
		for _, superset := range(workoutDay.Supersets) {
			for _, exercise := range(superset) {
				musclesHit[exercise.Target] = true
			}
		}
		for muscle := range musclesHit {
			frequencyPerMuscle[muscle]++
		}
	}
	// we need to know in which iteration of the frequency we are to select the correct amount of sets (later in cycle => round sets towards end)
	frequencyIterationPerMuscle := make(map[string]int)

	for workoutDayIdx, workoutDay := range(week) {

		exercisesPerMusclePerDay := make(map[string]int)
		for _, superset := range workoutDay.Supersets {
			for _, exercise := range superset {
				exercisesPerMusclePerDay[exercise.Target]++
			}
		}
		// same as frequency iteration, only difference is that we round towards beginning
		iterationsThisDayPerMuscle := make(map[string]int) 
		
		musclesHit := make(map[string]bool)
		for exerciseIdx, exercise := range(week[workoutDayIdx].Exercises) {
			setsForMusclePerWeek := func() int {
				for _, muscle := range(configPlan.Muscles) {
					if muscle.Name == exercise.Target {
						return muscle.Sets[weekIdx]
					}
					// should not happen if config checker works
				}
				return 0
			}()
			workoutDaySets := SpreadEnd(setsForMusclePerWeek, frequencyPerMuscle[exercise.Target], frequencyIterationPerMuscle[exercise.Target])
			exerciseSets := SpreadStart(workoutDaySets, exercisesPerMusclePerDay[exercise.Target], iterationsThisDayPerMuscle[exercise.Target])
			week[workoutDayIdx].Exercises[exerciseIdx].Sets = exerciseSets

			musclesHit[exercise.Target] = true
			iterationsThisDayPerMuscle[exercise.Target]++
		}
		for muscle := range musclesHit {
			frequencyIterationPerMuscle[muscle]++
		}
	}
	return week
}

func getWorkoutWeekdaysInOrder(configPlan *ConfigPlan) []time.Weekday {
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

// finds the earliest possible workout date after/on the start date
func getEarliestWorkoutDate(configPlan *ConfigPlan) time.Time {
	// no error should occur if read_config works (checks for validaty)
	date, _ := time.Parse(TIME_LAYOUT, configPlan.StartDate)
	workoutWeekdays := getWorkoutWeekdaysInOrder(configPlan)
	for i := 0; i < len(workoutWeekdays); i++ {
		_, isWorkoutDay := FindWeekday(workoutWeekdays, date.Weekday())
		if isWorkoutDay {
			break
		}
		date = NextDay(date)
	}
	// TODO this might cause some issues if config has no weekdays as the loop does not run
	return date
}

func getRotatedTrainingWeekdays(configPlan *ConfigPlan) []time.Weekday {
	earliestWorkoutDay := getEarliestWorkoutDate(configPlan).Weekday()
	workoutDays := getWorkoutWeekdaysInOrder(configPlan)
	idx, _ := FindWeekday(workoutDays, earliestWorkoutDay)
	return append(workoutDays[idx:], workoutDays[:idx]...)
}

func addDates(plan *Plan, configPlan *ConfigPlan) {
	currentDate := getEarliestWorkoutDate(configPlan)
	for weekIdx := range(plan) {
		for workoutDayIdx, workoutDay  := range(plan[weekIdx]) {
			for workoutDay.Weekday != strings.ToLower(currentDate.Weekday().String()) {
				currentDate = NextDay(currentDate)
			}
			plan[weekIdx][workoutDayIdx].Date = currentDate.Format(TIME_LAYOUT)
		}
	}
}