package domain

type Act int

const (
	DoMoveLeft Act = iota + 1
	DoMoveDown
	DoPunch
	Quit
)
