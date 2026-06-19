package main

import (
	"sync"
	"unicode/utf8"
)

var (
	WallLength             = 3
	IdleBoxerWidth         = 2
	IdleBoxerBehindLength  = 0
	IdleBoxerForwardLength = utf8.RuneCountInString(BodyChars[Idle][Arm][Left][Up])
	AllBodyParts           = map[BodyPart][]VerDirection{Head: {Up}, Shoulder: {Up, Down}, Arm: {Up, Down}}
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
	Name      string
	Health    int
	LastHit   int64 //timestamp
	Color     string
	Position  Position
	Opponent  IsOpponent
	Direction HorDirection
	Situation SituationType
	Area      map[Bounds]*Position
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
}

const (
	MaxHealthPoint int          = 100
	HeadHitHp      int          = 3
	ShoulderHitHp  int          = 1
	Min            Bounds       = false
	Max            Bounds       = true
	Main           IsOpponent   = false
	Opponent       IsOpponent   = true
	Up             VerDirection = false
	Down           VerDirection = true
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

var OppChan = make(chan Event)
var MainChan = make(chan Event)
var Router = map[IsOpponent]chan Event{
	Main:     MainChan,
	Opponent: OppChan,
}

var BodyChars = map[SituationType]map[BodyPart]map[HorDirection]map[VerDirection]string{
	Idle: {
		Head: {
			Left:  {Up: "(:/)"},
			Right: {Up: "(/:)"},
		},
		Arm: {
			Left:  {Up: "╭╭══O", Down: "╰╰══O"},
			Right: {Up: "O══╮╮", Down: "O══╯╯"},
		},
		Shoulder: {
			Left:  {Up: "\\\\", Down: "//"},
			Right: {Up: "//", Down: "\\\\"},
		},
	},
	Punch: {
		Shoulder: {
			Left:  {Up: "\\\\══O", Down: "//══O"},
			Right: {Up: "O══//", Down: "O══\\\\"},
		},
	},
	PunchInit: {
		Shoulder: {
			Left:  {Up: "\\\\═O", Down: "//═O"},
			Right: {Up: "O═//", Down: "O═\\\\"},
		},
	},
	HeadHit: {
		Head: {
			Left:  {Up: "**"},
			Right: {Up: "**"},
		},
		Arm: {
			Left:  {Up: "* ", Down: "* "},
			Right: {Up: " *", Down: " *"},
		},
		Shoulder: {
			Left:  {Up: "*", Down: "*"},
			Right: {Up: "*", Down: "*"},
		},
	},
	HeadInitHit: {
		Head: {
			Left:  {Up: "*"},
			Right: {Up: "*"},
		},
		Shoulder: {
			Left:  {Up: "*", Down: "*"},
			Right: {Up: "*", Down: "*"},
		},
	},
	ShoulderHit: {
		Head: {
			Left:  {Up: "* "},
			Right: {Up: " *"},
		},
		Arm: {
			Left:  {Up: "* ", Down: "* "},
			Right: {Up: " *", Down: " *"},
		},
		Shoulder: {
			Left:  {Up: "**", Down: "**"},
			Right: {Up: "**", Down: "**"},
		},
	},
	ShoulderInitHit: {
		Shoulder: {
			Left:  {Up: "*", Down: "*"},
			Right: {Up: "*", Down: "*"},
		},
	},
}

var Positions = map[SituationType]map[BodyPart]map[VerDirection]Position{
	Idle: {
		Head:     {Up: {X: 0, Y: 0}},
		Arm:      {Down: {X: 0, Y: +2}, Up: {X: 0, Y: -2}},
		Shoulder: {Down: {X: 0, Y: +1}, Up: {X: 0, Y: -1}},
	},
	Punch: {
		Shoulder: {Down: {X: +1, Y: +1}, Up: {X: +1, Y: -1}},
	},
	PunchInit: {
		Shoulder: {Down: {X: 0, Y: +1}, Up: {X: 0, Y: -1}},
	},
	HeadHit: {
		Head:     {Up: {X: -2, Y: 0}},
		Arm:      {Down: {X: -3, Y: +2}, Up: {X: -3, Y: -2}},
		Shoulder: {Down: {X: -2, Y: +1}, Up: {X: -2, Y: -1}},
	},
	HeadInitHit: {
		Head:     {Up: {X: -1, Y: 0}},
		Shoulder: {Down: {X: -2, Y: +1}, Up: {X: -2, Y: -1}},
	},
	ShoulderHit: {
		Head:     {Up: {X: -2, Y: 0}},
		Arm:      {Down: {X: -2, Y: +2}, Up: {X: -2, Y: -2}},
		Shoulder: {Down: {X: -2, Y: +1}, Up: {X: -2, Y: -1}},
	},
	ShoulderInitHit: {
		Shoulder: {Down: {X: -1, Y: +1}, Up: {X: -1, Y: -1}},
	},
}
