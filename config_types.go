package main

import "time"

type ConfigExercise struct {
	Name         string  `json:"name"`
	InitialOneRM float64 `json:"initial_1rm"`
	Reps         int     `json:"reps"`
	Target	     string  `json:"target"`
}

type ConfigMuscle struct {
	Name string `json:"name"`
	Sets [5]int `json:"sets"`
}

type ConfigVariations []string

type ConfigSplit struct {
	Name      string       `json:"name"`
	Executions []ConfigVariations `json:"executions"`
}

type ConfigPlan struct {
	Weekdays  map[time.Weekday]string
	StartDate string     `json:"start_date"`
	Splits    []ConfigSplit    `json:"splits"`
	Exercises []ConfigExercise `json:"exercises"`
	Muscles   []ConfigMuscle `json:"muscles"`
}
