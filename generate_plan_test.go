package main

import (
	"math"
	"testing"
)

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

func TestOneRM(t *testing.T) {
	configPlan := ReadJson("config.json")
	plan := GeneratePlan(configPlan)
	for weekIdx, week := range(plan) {
		for _, workoutDay := range(week) {
			for _, exercise := range(workoutDay.Exercises) {
				haveInitialOneRM := exercise.Weight * (1 + (exercise.Weight + float64(exercise.Reps)) / 30)
				haveInitialOneRM /= math.Pow(1 + ONE_RM_INCREASE, float64(weekIdx % (len(plan)-1)))
				wantInitialOnrRM := func() float64 {
					for _, exercise := range(configPlan.Exercises) {
						return exercise.InitialOneRM
					}
					return 0
				}()
				if (haveInitialOneRM - wantInitialOnrRM) < 0.25 {
					t.Error(exercise.Name + " 1RM is not correct")
				}
			}
		} 
	}
}

func TestDeloadOneRM(t *testing.T) {
	configPlan := ReadJson("config.json")
	plan := GeneratePlan(configPlan)
	for _, workoutDay := range(plan[len(plan)-1]) {
		for _, exercise := range(workoutDay.Exercises) {
			haveInitialOneRM := exercise.Weight * (1 + (exercise.Weight + float64(exercise.Reps)) / 30)
			wantInitialOnrRM := func() float64 {
				for _, exercise := range(configPlan.Exercises) {
					return exercise.InitialOneRM
				}
				return 0
			}()
			if (haveInitialOneRM - wantInitialOnrRM) < 0.25 {
				t.Error(exercise.Name + " 1RM is not correct")
			}
		}
	}
}