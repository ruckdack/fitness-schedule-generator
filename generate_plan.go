package main

import (
	"math"
	"sort"
	"time"
)

const ONE_RM_INCREASE float64 = 0.01

func GeneratePlan(configPlan *ConfigPlan) Plan {
	workoutsPerWeek := len(configPlan.Weekdays)
	trainigWeekdays := getTrainingWeekdaysInOrder(configPlan)
	variationCountPerSplit := make([]int, len(configPlan.Splits))

	generateWeek := func(weekIndex int) Week {
		
		generateWorkoutDay := func(workoutInWeek int) WorkoutDay {

			splitName := configPlan.Weekdays[trainigWeekdays[workoutInWeek]]
			getSplitIndex := func() int {
				for idx, split := range configPlan.Splits {
					if split.Name == splitName {
						return idx
					}
				}
				// should not happen if config checker works
				return 0
			}
			splitIndex := getSplitIndex()
			executions := configPlan.Splits[splitIndex].Executions
			exercises := make([]Exercise, len(executions))

			// select correct exercise from variations
			for idx, variations := range executions {
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
				weight := math.Pow(1 + ONE_RM_INCREASE, float64(weekIndex % 4)) * configPlan.Exercises[idx].InitialOneRM
				weight /= 1 + float64(reps + RirMapping[weekIndex]) / 30
				weight = math.Round(4 * weight) / 4
				// sets is added in a second round as we don't know yet which exercises will be selected for the entire week
				exercises[idx] = Exercise{
					Name: name,
					Target: target,
					Reps: reps,
					Weight: weight,
					Rir: RirMapping[weekIndex],
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

		// now we can add sets

		frequencyPerMuscle := make(map[string]int)
		for _, workoutDay := range(week) {
			// prevents double adding when a muscle is hit twice in a day
			musclesHit := make(map[string]bool)
			for _, exercise := range(workoutDay.Exercises) {
				musclesHit[exercise.Target] = true
			}
			for muscle := range musclesHit {
				frequencyPerMuscle[muscle]++
			}
		}
		// we need to know in which iteration of the frequency we are to select the correct amount of sets (later in cycle => round sets towards end)
		frequencyIterationPerMuscle := make(map[string]int)

		for workoutDayIdx, workoutDay := range(week) {

			exercisesPerMusclePerDay := make(map[string]int)
			for _, exercise := range workoutDay.Exercises {
				exercisesPerMusclePerDay[exercise.Target]++
			}
			// same as frequency iteration, only difference is that we round towards beginning
			iterationsThisDayPerMuscle := make(map[string]int) 
			
			musclesHit := make(map[string]bool)
			for exerciseIdx, exercise := range(week[workoutDayIdx].Exercises) {
				setsForMusclePerWeek := func() int {
					for _, muscle := range(configPlan.Muscles) {
						if muscle.Name == exercise.Target {
							return muscle.Sets[weekIndex]
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

	var plan Plan
	for weekIndex := range plan {
		plan[weekIndex] = generateWeek(weekIndex)
	}
	return plan
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