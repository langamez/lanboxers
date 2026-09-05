package game

import (
	"sync"

	"github.com/langamez/lanboxers/domain"
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
			Area:      domain.Area{},
			Health:    domain.MaxHealth,
			Situation: domain.Situation{SituationType: domain.Idle},
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

func Snapshot(boxer domain.BaseBoxer) domain.BaseBoxer {
	return domain.BaseBoxer{
		Area:      boxer.Area,
		Name:      boxer.Name,
		Color:     boxer.Color,
		Health:    boxer.Health,
		LastHit:   boxer.LastHit,
		Position:  boxer.Position,
		Direction: boxer.Direction,
		Situation: boxer.Situation,
	}
}

func (g *Game) UpdateBoxer(
	id domain.PlayerID,
	toBoxer domain.BaseBoxer,
	parts map[domain.BodyPart][]domain.Direction,
) {
	boxer := g.Boxers[id]
	boxer.Lock.Lock()
	defer boxer.Lock.Unlock()

	g.Boxers[id].BaseBoxer.Position = toBoxer.Position
	g.Boxers[id].BaseBoxer.Situation = toBoxer.Situation
	g.Boxers[id].BaseBoxer.Direction = toBoxer.Direction

	g.RenderChans[id] <- domain.RenderCommand{Boxer: toBoxer, Parts: parts}
}
