package domain

type Action int

const (
	DoMoveLeft Action = iota + 1
	DoMoveDown
	DoPunch
	Quit
)
