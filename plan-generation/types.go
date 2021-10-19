package main

const ONE_RM_INCREASE float64 = 0.01

var WEEK_TYPES = [5]string{"intro", "accumulation1", "accumulation2", "overreaching", "deload"}
var RIR_MAPPING = [5]int{2, 1, 1, 0, 4}

type Plan [5]Week

type Week []WorkoutDay

type WorkoutDay struct {
	Date      string     `json:"date"`
	Weekday   string     `json:"weekday"`
	WeekType  string     `json:"weektype"`
	Split     string     `json:"split"`
	Supersets []Superset `json:"exercises"`
}

type Superset []Exercise

type Exercise struct {
	Name   string  `json:"name"`
	Target string  `json:"target"`
	Weight float64 `json:"weight"`
	Reps   int     `json:"reps"`
	Sets   int     `json:"sets"`
	Rir    int     `json:"rir"`
}