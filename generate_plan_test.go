package main

import (
	"math"
	"strings"
	"testing"
	"time"
)

const configJson string = `{"weekdays":{"mon":"upper","tue":"lower","wed":"upper","fri":"lower","sat":"upper"},"start_date":"2021-10-15","splits":[{"name":"upper","executions":[["incline benchpress (dumbbell)","benchpress"],["high row (pulley)"],["lateral raise (dumbbell)","lateral raise (machine)"],["lat pulldown"],["incline curl (dumbbell)","curl machine"],["triceps extension","overhead triceps extension"]]},{"name":"lower","executions":[["iwas für quads"],["leg extension"],["RDL","hip thrust"],["laying leg curl"],["seated calf raise"],["standing calf raise"],["crunch machine"]]}],"exercises":[{"name":"incline benchpress (dumbbell)","initial_1rm":107,"reps":9,"target":"chest"},{"name":"benchpress","initial_1rm":100,"reps":8,"target":"chest"},{"name":"high row (pulley)","initial_1rm":100,"reps":8,"target":"lower traps / rhomboids"},{"name":"lateral raise (dumbbell)","initial_1rm":100,"reps":8,"target":"delts"},{"name":"lateral raise (machine)","initial_1rm":100,"reps":8,"target":"delts"},{"name":"lat pulldown","initial_1rm":100,"reps":8,"target":"lat"},{"name":"incline curl (dumbbell)","initial_1rm":100,"reps":8,"target":"biceps"},{"name":"curl machine","initial_1rm":100,"reps":8,"target":"biceps"},{"name":"triceps extension","initial_1rm":100,"reps":8,"target":"triceps"},{"name":"overhead triceps extension","initial_1rm":100,"reps":8,"target":"triceps"},{"name":"crunch machine","initial_1rm":100,"reps":8,"target":"abs"},{"name":"standing calf raise","initial_1rm":100,"reps":8,"target":"calves"},{"name":"seated calf raise","initial_1rm":100,"reps":8,"target":"calves"},{"name":"laying leg curl","initial_1rm":100,"reps":8,"target":"hamstrings"},{"name":"hip thrust","initial_1rm":100,"reps":8,"target":"glute"},{"name":"RDL","initial_1rm":100,"reps":8,"target":"hamstrings"},{"name":"leg extension","initial_1rm":100,"reps":8,"target":"quads"},{"name":"iwas für quads","initial_1rm":100,"reps":8,"target":"quads"}],"muscles":[{"name":"quads","sets":[12,14,14,16,8]},{"name":"chest","sets":[12,14,14,16,8]},{"name":"hamstrings","sets":[12,14,14,16,8]},{"name":"delts","sets":[12,14,14,16,8]},{"name":"calves","sets":[12,14,14,16,8]},{"name":"glute","sets":[12,14,14,16,8]},{"name":"lower traps / rhomboids","sets":[12,14,14,16,8]},{"name":"lat","sets":[12,14,14,16,8]},{"name":"delts","sets":[12,14,14,16,8]},{"name":"abs","sets":[12,14,14,16,8]},{"name":"biceps","sets":[12,14,14,16,8]},{"name":"triceps","sets":[12,14,14,16,8]}]}`

func TestSets(t *testing.T) {
	configPlan := ReadJson(configJson)
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
	configPlan := ReadJson(configJson)
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
				if (haveInitialOneRM - wantInitialOnrRM) > 0.25 {
					t.Error(exercise.Name + " 1RM is not correct")
				}
			}
		} 
	}
}

func TestDeloadOneRM(t *testing.T) {
	configPlan := ReadJson(configJson)
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

func TestDates(t *testing.T) {
	configPlan := ReadJson(configJson)
	plan := GeneratePlan(configPlan)
	for _, week := range(plan) {
		for _, workoutDay := range(week) {
			date, err := time.Parse(TIME_LAYOUT, workoutDay.Date)
			if err != nil {
				t.Error(workoutDay.Date + " cannot be parsed")
			}
			if strings.ToLower(date.Weekday().String()) != workoutDay.Weekday {
				t.Error(workoutDay.Date + " does not match the weekday:\nwant: " + strings.ToLower(date.Weekday().String()) + "\nhave:" + workoutDay.Weekday)
			}
		} 
	}
}