package main

var WeekTypes = [5]string{"intro", "accumulation1", "accumulation2", "overreaching", "deload"}

type Plan [5]Week

type Week []WorkoutDay

type WorkoutDay struct {
	Date    string
	Weekday string//time.Weekday
	WeekType string
	Split string
	Exercises []Exercise
}

type Exercise struct {
	Name string
	Weight float32
	Reps int
	Sets int
	Rir int
}