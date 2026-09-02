package domain

type Bound int
type BodyPart int
type PlayerID int
type Direction int
type Situation int
type VerticalDirection int

const (
	Min Bound = iota
	Max
)

const (
	PlayerMain PlayerID = iota
	PlayerOpponent
)

func (p PlayerID) Opposite() PlayerID {
	return 1 - p
}

const (
	Left Direction = iota
	Right
	Upper
	Lower
)

func (d Direction) Opposite() Direction {
	switch d {
	case Left:
		return Right
	case Right:
		return Left
	case Upper:
		return Lower
	case Lower:
		return Upper
	}
	return d
}

const (
	Arm BodyPart = iota + 1
	Head
	Shoulder
)

const (
	Idle Situation = iota + 1
	PunchInit
	Punch
	PunchPeak
	HeadHit
	HeadInitHit
	ShoulderHit
	ShoulderInitHit
)
