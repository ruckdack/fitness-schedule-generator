package main

import "time"

type Exercise struct {
	Name         string  `json:"name"`
	InitialOneRM float32 `json:"initial_1rm"`
	Reps         int     `json:"reps"`
}

type Execution struct {
	Variations	[]string `json:"variations"`
	Sets		[5]int 	 `json:"sets"`
}

type Split struct {
	Name      string       `json:"name"`
	Executions []Execution `json:"executions"`
}

type PlanConfig struct {
	Weekdays  map[time.Weekday]string
	StartDate string     `json:"start_date"`
	Splits    []Split    `json:"splits"`
	Exercises []Exercise `json:"exercises"`
}
