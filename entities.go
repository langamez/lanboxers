package main

type Position struct {
	X int // horizontal position
	Y int // vertical position
}

type Boxer struct {
	Color string
	Pose  Position
}

type Cage struct {
	Color string
	Limit Position
}
