package main

import "testing"

func TestSets(t *testing.T) {
	configPlan := ReadJson("config.json")
	plan := GeneratePlan(configPlan)
	for weekIdx, week := range(plan) {
		haveSetsPerMuscle := make(map[string]int)
		for _, workoutDay := range(week) {
			for _, exercise := range(workoutDay.Exercises) {
				haveSetsPerMuscle[exercise.Target] += exercise.Sets
			}
		}
		wantSetsPerMuscle := make(map[string]int)
		for _, muscle := range(configPlan.Muscles) {
			wantSetsPerMuscle[muscle.Name] = muscle.Sets[weekIdx]
		}
		for muscle := range(wantSetsPerMuscle) {
			if wantSetsPerMuscle[muscle] != haveSetsPerMuscle[muscle] {
				t.Error(muscle + " is wrong")
			}
		}
		for muscle := range(haveSetsPerMuscle) {
			if wantSetsPerMuscle[muscle] != haveSetsPerMuscle[muscle] {
				t.Error(muscle + " is wrong")
			}
		}
	}
}