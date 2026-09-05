package game

import (
	"github.com/langamez/lanboxers/domain"
	"sync"
)

func NewBoxer(
	name string,
	color string,
	position domain.Position,
	direction domain.Direction,
) *domain.Boxer {
	return &domain.Boxer{
		Lock: sync.Mutex{},
		BaseBoxer: domain.BaseBoxer{
			Name:      name,
			Color:     color,
			Position:  position,
			Direction: direction,
			Situation: domain.Idle,
			Area:      domain.Area{},
			Health:    domain.MaxHealth,
		},
	}
}

func NewBoxers(cage domain.Cage) domain.Boxers {
	return domain.Boxers{
		domain.PlayerMain: NewBoxer(
			"User1",
			"\033[31m",
			domain.Position{
				X: cage.Area.Max.X / 3,
				Y: cage.Area.Max.Y / 2,
			},
			domain.Left,
		),

		domain.PlayerOpponent: NewBoxer(
			"User2",
			"\033[32m",
			domain.Position{
				X: (cage.Area.Max.X / 3) * 2,
				Y: cage.Area.Max.Y / 2,
			},
			domain.Right,
		),
	}
}
