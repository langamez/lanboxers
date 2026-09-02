package domain

import (
	"sync"
)

type Area struct {
	Min Position
	Max Position
}

type Boxer struct {
	BaseBoxer
	Lock sync.Mutex
}

type BaseBoxer struct {
	Name         string
	Health       int
	LastHit      int64
	Color        string
	Position     Position
	Direction    Direction
	Situation    Situation
	SituationDir Direction
	Area         Area
}

type PunchDetail struct {
	Attacker  *Boxer
	Direction Direction
	Power     bool
}

func (b *BaseBoxer) AddHealth(amount int) {
	b.Health += amount
	if b.Health > MaxHealth {
		b.Health = MaxHealth
	}
}
