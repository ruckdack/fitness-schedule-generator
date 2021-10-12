package main

import "time"

var WeekNames = [5]string{"intro", "accumulation1", "accumulation2", "overreaching", "deload"}

type Plan [5]*Week

type Week []*WorkoutDay

type WorkoutDay struct {
	Date    string
	Weekday time.Weekday
	Split string
	Exercises []*Exercise
}

type Exercise struct {
	Name string
	Weight float32
	Reps int
	Sets int
	Rir int
}