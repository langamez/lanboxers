package game

import (
	"lanBox/domain"
)

func NewCage(width, height int) domain.Cage {
	return domain.Cage{
		Color: "\033[44m",
		Area: domain.Area{
			Min: domain.Position{
				X: domain.WallLength + 1,
				Y: domain.WallLength + 1,
			},
			Max: domain.Position{
				X: (width - 10) - domain.WallLength,
				Y: (height - 10) - domain.WallLength,
			},
		},
	}
}
