package main

func GeneratePlan(configPlan *ConfigPlan) Plan {
	workoutsPerWeek := len(configPlan.Weekdays)
	frequencyPerExercise := getFrequencyPerExercise(configPlan)

	func generateWeek(weekIndex int) Week {

		func generateWorkoutDay(workoutInWeek int) *WorkoutDay {

		}

		var week Week

		for workoutInWeek := range workoutsPerWeek {
			week[workoutInWeek] = generateWorkoutDay(workoutInWeek)
		}
		return week
	}

	plan := [5]Week
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

func getMaxVariation(configPlan *ConfigPlan) {
	currentMax = 0
	for _, split := range configPlan.Splits {
		for _, execution := range split.Executions {
			currentMax = math.Max(currentMax, len(execution.Variations))
		}
	}
	return currentMax
}