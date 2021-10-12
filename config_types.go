package main

import "time"

type ConfigExercise struct {
	Name         string  `json:"name"`
	InitialOneRM float32 `json:"initial_1rm"`
	Reps         int     `json:"reps"`
}

type ConfigExecution struct {
	Variations	[]string `json:"variations"`
	Sets		[5]int 	 `json:"sets"`
}

type ConfigSplit struct {
	Name      string       `json:"name"`
	Executions []ConfigExecution `json:"executions"`
}

type ConfigPlan struct {
	Weekdays  map[time.Weekday]string
	StartDate string     `json:"start_date"`
	Splits    []ConfigSplit    `json:"splits"`
	Exercises []ConfigExercise `json:"exercises"`
}
