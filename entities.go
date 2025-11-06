package main

type Direction int
type DirectionType int
type SituationType int

type Position struct {
	X int // horizontal position
	Y int // vertical position
}

type Boxer struct {
	Color     string
	Pose      Position
	Situation SituationType
	Direction DirectionType
}

type Cage struct {
	Color string
	Limit Position
}

const (
	Idle SituationType = 1

	HeadHitInit SituationType = 2
	HeadHit     SituationType = 3

	ShoulderHitInit SituationType = 4
	ShoulderHit     SituationType = 5

	UpPunchInit SituationType = 6
	UpPunch     SituationType = 7

	DownPunchInit SituationType = 8
	DownPunch     SituationType = 9
)

const (
	Up    Direction     = 1
	Down  Direction     = 2
	Left  DirectionType = 1
	Right DirectionType = 2
)

type BodySituation struct {
	head     map[DirectionType]string
	arm      map[DirectionType]map[Direction]string
	shoulder map[DirectionType]map[Direction]string

	headPose     map[DirectionType]Position
	armPose      map[DirectionType]map[Direction]Position
	shoulderPose map[DirectionType]map[Direction]Position
}

var Situations = map[SituationType]*BodySituation{
	Idle: {
		head: map[DirectionType]string{
			Right: "(:", Left: ":)"},
		headPose: map[DirectionType]Position{
			Right: {X: 0, Y: 0},
			Left:  {X: 0, Y: 0},
		},
		arm: map[DirectionType]map[Direction]string{
			Right: {Up: "O══╮╮", Down: "O══╯╯"},
			Left:  {Up: "╭╭══O", Down: "╰╰══O"}},
		shoulder: map[DirectionType]map[Direction]string{
			Right: {Up: "()", Down: "()"},
			Left:  {Up: "()", Down: "()"}},
		armPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -4, Y: -2}, Down: {X: -4, Y: 2}},
			Left:  {Up: {X: 1, Y: -2}, Down: {X: 1, Y: 2}},
		},
		shoulderPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: 0, Y: 1}, Down: {X: 0, Y: -1}},
			Left:  {Up: {X: 0, Y: -1}, Down: {X: 0, Y: 1}},
		},
	},
	HeadHitInit: {
		head: map[DirectionType]string{
			Right: "): *", Left: "* :("},
		headPose: map[DirectionType]Position{
			Right: {X: 0, Y: 0},
			Left:  {X: -2, Y: 0},
		},
		arm: map[DirectionType]map[Direction]string{
			Right: {Up: "O══╮╮", Down: "O══╯╯"},
			Left:  {Up: "╭╭══O", Down: "╰╰══O"}},
		shoulder: map[DirectionType]map[Direction]string{
			Right: {Up: "()", Down: "()"},
			Left:  {Up: "()", Down: "()"}},
		armPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -4, Y: -2}, Down: {X: -4, Y: 2}},
			Left:  {Up: {X: 1, Y: -2}, Down: {X: 1, Y: 2}},
		},
		shoulderPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: 0, Y: -1}, Down: {X: 0, Y: 1}},
			Left:  {Up: {X: 0, Y: -1}, Down: {X: 0, Y: 1}},
		},
	},
	HeadHit: {
		head: map[DirectionType]string{
			Right: "): * *", Left: "* * :("},
		headPose: map[DirectionType]Position{
			Right: {X: 0, Y: 0},
			Left:  {X: -4, Y: 0},
		},
		arm: map[DirectionType]map[Direction]string{
			Right: {Up: "O══╮╮", Down: "O══╯╯"},
			Left:  {Up: "╭╭══O", Down: "╰╰══O"}},
		shoulder: map[DirectionType]map[Direction]string{
			Right: {Up: "()  *", Down: "()  *"},
			Left:  {Up: "*  ()", Down: "*  ()"}},
		armPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -4, Y: -2}, Down: {X: -4, Y: 2}},
			Left:  {Up: {X: 1, Y: -2}, Down: {X: 1, Y: 2}},
		},
		shoulderPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: 0, Y: 1}, Down: {X: 0, Y: -1}},
			Left:  {Up: {X: -3, Y: -1}, Down: {X: -3, Y: 1}},
		},
	},
	ShoulderHitInit: {
		head: map[DirectionType]string{
			Right: "):", Left: ":("},
		headPose: map[DirectionType]Position{
			Right: {X: 0, Y: 0},
			Left:  {X: 0, Y: 0},
		},
		arm: map[DirectionType]map[Direction]string{
			Right: {Up: "O══╮╮", Down: "O══╯╯"},
			Left:  {Up: "╭╭══O", Down: "╰╰══O"}},
		shoulder: map[DirectionType]map[Direction]string{
			Right: {Up: "() *", Down: "() *"},
			Left:  {Up: "* ()", Down: "* ()"}},
		armPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -4, Y: -2}, Down: {X: -4, Y: 2}},
			Left:  {Up: {X: 1, Y: -2}, Down: {X: 1, Y: 2}},
		},
		shoulderPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: 0, Y: 1}, Down: {X: 0, Y: -1}},
			Left:  {Up: {X: -2, Y: 1}, Down: {X: -2, Y: -1}},
		},
	},
	ShoulderHit: {
		head: map[DirectionType]string{
			Right: "):  *", Left: "*  :("},
		headPose: map[DirectionType]Position{
			Right: {X: 0, Y: 0},
			Left:  {X: -3, Y: 0},
		},
		arm: map[DirectionType]map[Direction]string{
			Right: {Up: "O══╮╮", Down: "O══╯╯"},
			Left:  {Up: "╭╭══O", Down: "╰╰══O"}},
		shoulder: map[DirectionType]map[Direction]string{
			Right: {Up: "() * *", Down: "() * *"},
			Left:  {Up: "* * ()", Down: "* * ()"}},
		armPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -4, Y: -2}, Down: {X: -4, Y: 2}},
			Left:  {Up: {X: 1, Y: -2}, Down: {X: 1, Y: 2}},
		},
		shoulderPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: 0, Y: 1}, Down: {X: 0, Y: -1}},
			Left:  {Up: {X: -4, Y: 1}, Down: {X: -4, Y: -1}},
		},
	},
	UpPunchInit: {
		head: map[DirectionType]string{
			Right: ">:", Left: ":<"},
		headPose: map[DirectionType]Position{
			Right: {X: 0, Y: 0},
			Left:  {X: 0, Y: 0},
		},
		arm: map[DirectionType]map[Direction]string{
			Right: {Up: "", Down: "O══╯╯"},
			Left:  {Up: "", Down: "╰╰══O"}},
		shoulder: map[DirectionType]map[Direction]string{
			Right: {Up: "O══()", Down: "()"},
			Left:  {Up: "()══O", Down: "()"}},
		armPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -1, Y: -2}, Down: {X: -4, Y: 2}},
			Left:  {Up: {X: 1, Y: 2}, Down: {X: 1, Y: 2}},
		},
		shoulderPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -3, Y: -1}, Down: {X: 0, Y: 1}},
			Left:  {Up: {X: 0, Y: -1}, Down: {X: 0, Y: 1}},
		},
	},
	UpPunch: {
		head: map[DirectionType]string{
			Right: ">:", Left: ":<"},
		headPose: map[DirectionType]Position{
			Right: {X: 0, Y: 0},
			Left:  {X: 0, Y: 0},
		},
		arm: map[DirectionType]map[Direction]string{
			Right: {Up: "", Down: "O══╯╯"},
			Left:  {Up: "", Down: "╰╰══O"}},
		shoulder: map[DirectionType]map[Direction]string{
			Right: {Up: "O═══()", Down: "()"},
			Left:  {Up: "()═══O", Down: "()"}},
		armPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -1, Y: -2}, Down: {X: -1, Y: 2}},
			Left:  {Up: {X: -1, Y: 2}, Down: {X: -2, Y: 2}},
		},
		shoulderPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -5, Y: -1}, Down: {X: 1, Y: 1}},
			Left:  {Up: {X: 1, Y: -1}, Down: {X: -1, Y: 1}},
		},
	},
	DownPunchInit: {
		head: map[DirectionType]string{
			Right: ">:", Left: ":<"},
		headPose: map[DirectionType]Position{
			Right: {X: 0, Y: 0},
			Left:  {X: 0, Y: 0},
		},
		arm: map[DirectionType]map[Direction]string{
			Right: {Up: "O══╮╮", Down: ""},
			Left:  {Up: "╭╭══O", Down: ""}},
		shoulder: map[DirectionType]map[Direction]string{
			Right: {Up: "()", Down: "O══()"},
			Left:  {Up: "()", Down: "()══O"}},
		armPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -4, Y: -2}, Down: {X: 2, Y: 2}},
			Left:  {Up: {X: 1, Y: -2}, Down: {X: 1, Y: -2}},
		},
		shoulderPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: 0, Y: -1}, Down: {X: -4, Y: 1}},
			Left:  {Up: {X: 0, Y: -1}, Down: {X: 0, Y: 1}},
		},
	},
	DownPunch: {
		head: map[DirectionType]string{
			Right: ">:", Left: ":<"},
		headPose: map[DirectionType]Position{
			Right: {X: 0, Y: 0},
			Left:  {X: 0, Y: 0},
		},
		arm: map[DirectionType]map[Direction]string{
			Right: {Up: "O══╮╮", Down: ""},
			Left:  {Up: "╭╭══O", Down: ""}},
		shoulder: map[DirectionType]map[Direction]string{
			Right: {Up: "()", Down: "O═══()"},
			Left:  {Up: "()", Down: "()═══O"}},
		armPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: -1, Y: -2}, Down: {X: -1, Y: 2}},
			Left:  {Up: {X: -2, Y: -2}, Down: {X: 1, Y: 2}},
		},
		shoulderPose: map[DirectionType]map[Direction]Position{
			Right: {Up: {X: 1, Y: -1}, Down: {X: -5, Y: 1}},
			Left:  {Up: {X: -1, Y: -1}, Down: {X: 1, Y: 1}},
		},
	},
}
