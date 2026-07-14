package main

import (
	"context"
	"sync"
	"unicode/utf8"
)

var (
	WallLength             = 3
	IdleBoxerWidth         = 2
	IdleBoxerBehindLength  = 0
	IdleBoxerForwardLength = utf8.RuneCountInString(BodyChars[Idle][Arm][Left][Left])
	AllBodyParts           = map[BodyPart][]HorDirection{Head: {Left}, Shoulder: {Left, Right}, Arm: {Left, Right}}
)

type Act int
type Bounds bool
type BodyPart int
type IsOpponent bool
type HorDirection bool
type VerDirection bool
type SituationType int
type Boxers map[IsOpponent]*Boxer

type Position struct {
	X int // horizontal position
	Y int // vertical position
}

type BaseBoxer struct {
	Name         string
	Health       int
	LastHit      int64 //timestamp
	Color        string
	Position     Position
	Opponent     IsOpponent
	Direction    HorDirection
	Situation    SituationType
	SituationDir HorDirection
	Area         map[Bounds]*Position
}

type Boxer struct {
	*BaseBoxer
	Lock sync.Mutex
}

type Event struct {
	Act   Act
	Actor IsOpponent
}

type Cage struct {
	Color string
	Area  map[Bounds]*Position
}

type GameInfo struct {
	Cage   Cage
	Boxers Boxers
	Cancel context.CancelFunc
}

const (
	MaxHealthPoint int          = 100
	HeadHitHp      int          = 3
	ShoulderHitHp  int          = 1
	Min            Bounds       = false
	Max            Bounds       = true
	Main           IsOpponent   = false
	Opponent       IsOpponent   = true
	Upper          VerDirection = false
	Lower          VerDirection = true
	Left           HorDirection = false
	Right          HorDirection = true

	Head BodyPart = iota + 1
	Shoulder
	Arm

	Idle SituationType = iota + 1
	PunchInit
	Punch
	PunchPeak
	HeadHit
	HeadInitHit
	ShoulderHit
	ShoulderInitHit

	DoMoveLeft Act = iota + 1
	DoMoveDown
	DoPunch
	Quit
)

var MainChan = make(chan Event)

var BodyChars = map[SituationType]map[BodyPart]map[HorDirection]map[HorDirection]string{
	Idle: {
		Head: {
			Left:  {Left: "(:/)"},
			Right: {Left: "(/:)"},
		},
		Arm: {
			Left:  {Left: "╭╭══O", Right: "╰╰══O"},
			Right: {Right: "O══╮╮", Left: "O══╯╯"},
		},
		Shoulder: {
			Left:  {Left: "\\\\", Right: "//"},
			Right: {Right: "//", Left: "\\\\"},
		},
	},
	Punch: {
		Shoulder: {
			Left:  {Left: "\\\\══O", Right: "//══O"},
			Right: {Right: "O══//", Left: "O══\\\\"},
		},
	},
	PunchInit: {
		Shoulder: {
			Left:  {Left: "\\\\═O", Right: "//═O"},
			Right: {Right: "O═//", Left: "O═\\\\"},
		},
	},
	HeadHit: {
		Head: {
			Left:  {Left: "**"},
			Right: {Right: "**"},
		},
		Arm: {
			Left:  {Left: "* ", Right: "* "},
			Right: {Right: " *", Left: " *"},
		},
		Shoulder: {
			Left:  {Left: "*", Right: "*"},
			Right: {Right: "*", Left: "*"},
		},
	},
	HeadInitHit: {
		Head: {
			Left:  {Left: "*"},
			Right: {Right: "*"},
		},
		Shoulder: {
			Left:  {Left: "*", Right: "*"},
			Right: {Right: "*", Left: "*"},
		},
	},
	ShoulderHit: {
		Head: {
			Left:  {Left: "* "},
			Right: {Left: " *"},
		},
		Arm: {
			Left:  {Left: "* ", Right: "* "},
			Right: {Right: " *", Left: " *"},
		},
		Shoulder: {
			Left:  {Left: "**", Right: "**"},
			Right: {Right: "**", Left: "**"},
		},
	},
	ShoulderInitHit: {
		Shoulder: {
			Left:  {Left: "*", Right: "*"},
			Right: {Right: "*", Left: "*"},
		},
	},
}

var Positions = map[SituationType]map[BodyPart]map[HorDirection]Position{
	Idle: {
		Head:     {Left: {X: 0, Y: 0}},
		Arm:      {Right: {X: 0, Y: +2}, Left: {X: 0, Y: -2}},
		Shoulder: {Right: {X: 0, Y: +1}, Left: {X: 0, Y: -1}},
	},
	Punch: {
		Shoulder: {Right: {X: +1, Y: +1}, Left: {X: +1, Y: -1}},
	},
	PunchInit: {
		Shoulder: {Right: {X: 0, Y: +1}, Left: {X: 0, Y: -1}},
	},
	HeadHit: {
		Head:     {Left: {X: -2, Y: 0}},
		Arm:      {Right: {X: -3, Y: +2}, Left: {X: -3, Y: -2}},
		Shoulder: {Right: {X: -2, Y: +1}, Left: {X: -2, Y: -1}},
	},
	HeadInitHit: {
		Head:     {Left: {X: -1, Y: 0}},
		Shoulder: {Right: {X: -2, Y: +1}, Left: {X: -2, Y: -1}},
	},
	ShoulderHit: {
		Head:     {Left: {X: -2, Y: 0}},
		Arm:      {Right: {X: -2, Y: +2}, Left: {X: -2, Y: -2}},
		Shoulder: {Right: {X: -2, Y: +1}, Left: {X: -2, Y: -1}},
	},
	ShoulderInitHit: {
		Shoulder: {Right: {X: -1, Y: +1}, Left: {X: -1, Y: -1}},
	},
}
