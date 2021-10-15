package main

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

const configJson string = `{    "weekdays": {        "mon": "upper",        "tue": "lower",        "wed": "upper",        "fri": "lower",        "sat": "upper"    },    "start_date": "2021-10-16",    "splits": [        {            "name": "upper",            "exercises": [                [                    ["incline benchpress (dumbbell)", "benchpress"],                    ["high row (pulley)"]                ],                [                    ["lateral raise (dumbbell)", "lateral raise (machine)"],                    ["lat pulldown"]                ],                [                    ["incline curl (dumbbell)", "curl machine"],                    ["triceps extension", "overhead triceps extension"]                ]            ]        },        {            "name": "lower",            "exercises": [                [["leg press"]],                [["RDL", "hip thrust"]],                [["leg extension"], ["laying leg curl"]],                [["seated calf raise"], ["standing calf raise"]],                [["crunch machine"]]            ]        }    ],    "exercises": [        {            "name": "incline benchpress (dumbbell)",            "initial_1rm": 10,            "reps": 9,            "target": "chest"        },        {            "name": "benchpress",            "initial_1rm": 200,            "reps": 8,            "target": "chest"        },        {            "name": "high row (pulley)",            "initial_1rm": 300,            "reps": 8,            "target": "lower traps / rhomboids"        },        {            "name": "lateral raise (dumbbell)",            "initial_1rm": 400,            "reps": 8,            "target": "delts"        },        {            "name": "lateral raise (machine)",            "initial_1rm": 500,            "reps": 8,            "target": "delts"        },        {            "name": "lat pulldown",            "initial_1rm": 600,            "reps": 8,            "target": "lat"        },        {            "name": "incline curl (dumbbell)",            "initial_1rm": 700,            "reps": 8,            "target": "biceps"        },        {            "name": "curl machine",            "initial_1rm": 800,            "reps": 8,            "target": "biceps"        },        {            "name": "triceps extension",            "initial_1rm": 900,            "reps": 8,            "target": "triceps"        },        {            "name": "overhead triceps extension",            "initial_1rm": 1000,            "reps": 8,            "target": "triceps"        },        {            "name": "crunch machine",            "initial_1rm": 1100,            "reps": 8,            "target": "abs"        },        {            "name": "standing calf raise",            "initial_1rm": 1200,            "reps": 8,            "target": "calves"        },        {            "name": "seated calf raise",            "initial_1rm": 1300,            "reps": 8,            "target": "calves"        },        {            "name": "laying leg curl",            "initial_1rm": 1400,            "reps": 8,            "target": "hamstrings"        },        {            "name": "hip thrust",            "initial_1rm": 1500,            "reps": 8,            "target": "glute"        },        {            "name": "RDL",            "initial_1rm": 1600,            "reps": 8,            "target": "hamstrings"        },        {            "name": "leg extension",            "initial_1rm": 1700,            "reps": 8,            "target": "quads"        },        {            "name": "leg press",            "initial_1rm": 1800,            "reps": 8,            "target": "quads"        }    ],    "muscles": [        {            "name": "quads",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "chest",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "hamstrings",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "delts",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "calves",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "glute",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "lower traps / rhomboids",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "lat",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "delts",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "abs",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "biceps",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "triceps",            "sets": [12, 14, 14, 16, 8]        }    ]}`

func TestSets(t *testing.T) {
	configPlan, err := ReadJson(configJson)
	if (err != nil) {
		t.Error(err)
		return
	}
	plan := GeneratePlan(configPlan)
	for weekIdx, week := range(plan) {
		haveSetsPerMuscle := make(map[string]int)
		for _, workoutDay := range(week) {
			for _, superset := range(workoutDay.Supersets) {
				for _, exercise := range(superset) {
					haveSetsPerMuscle[exercise.Target] += exercise.Sets
				}
			}
		}
		wantSetsPerMuscle := make(map[string]int)
		for _, muscle := range(configPlan.Muscles) {
			wantSetsPerMuscle[muscle.Name] = muscle.Sets[weekIdx]
		}
		for muscle := range(wantSetsPerMuscle) {
			if wantSetsPerMuscle[muscle] != haveSetsPerMuscle[muscle] {
				t.Error(muscle + " volume is wrong\nwant: " + fmt.Sprint(wantSetsPerMuscle[muscle]) + "\nhave: " + fmt.Sprint(haveSetsPerMuscle[muscle]))
			}
		}
	}
}

func TestOneRM(t *testing.T) {
	configPlan, err := ReadJson(configJson)
	if (err != nil) {
		t.Error(err)
		return
	}
	plan := GeneratePlan(configPlan)
	for weekIdx, week := range(plan) {
		for _, workoutDay := range(week) {
			for _, superset := range(workoutDay.Supersets) {
				for _, exercise := range(superset) {
					haveInitialOneRM := exercise.Weight * (1 + float64(exercise.Reps + exercise.Rir) / 30)
				haveInitialOneRM /= math.Pow(1 + ONE_RM_INCREASE, float64(weekIdx % (len(plan)-1)))
				wantInitialOneRM := func() float64 {
					for _, ex := range(configPlan.Exercises) {
						if (exercise.Name == ex.Name) {
							return ex.InitialOneRM
						}
					}
					return 0
				}()
				wantString := fmt.Sprint(wantInitialOneRM)
				haveString := fmt.Sprint(haveInitialOneRM)
				if (haveInitialOneRM - wantInitialOneRM) > 0.25 {
					t.Error(exercise.Name + " 1RM is not correct\nwant: " + wantString + "\nhave: " + haveString)
				}
				}
			}
		} 
	}
}

func TestDeloadOneRM(t *testing.T) {
	configPlan, err := ReadJson(configJson)
	if (err != nil) {
		t.Error(err)
		return
	}
	plan := GeneratePlan(configPlan)
	for _, workoutDay := range(plan[len(plan)-1]) {
		for _, superset := range(workoutDay.Supersets) {
			for _, exercise := range(superset) {
				haveInitialOneRM := exercise.Weight * (1 + float64(exercise.Reps + exercise.Rir) / 30)
				wantInitialOneRM := func() float64 {
					for _, ex := range(configPlan.Exercises) {
						if (exercise.Name == ex.Name) {
							return ex.InitialOneRM
						}
					}
					return 0
				}()
				wantString := fmt.Sprint(wantInitialOneRM)
				haveString := fmt.Sprint(haveInitialOneRM)
				if (haveInitialOneRM - wantInitialOneRM) > 0.25 {
					t.Error(exercise.Name + " 1RM is not correct\nwant: " + wantString + "\nhave: " + haveString)
				}
			}
		}
	}
}

func TestDates(t *testing.T) {
	configPlan, err := ReadJson(configJson)
	if (err != nil) {
		t.Error(err)
		return
	}
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