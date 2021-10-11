package main

import "time"

type Exercise struct {
	Name         string  `json:"name"`
	InitialOneRM float32 `json:"initial_1rm"`
	Reps         int     `json:"reps"`
	Sets         []int   `json:"sets"`
}

type ExerciseVariation []Exercise

type Split struct {
	Name      string              `json:"name"`
	Exercises []ExerciseVariation `json:"exercises"`
}

type PlanConfig struct {
	Weekdays  map[time.Weekday]string
	StartDate string  `json:"start_date"`
	Splits    []Split `json:"splits"`
}