package game

import (
	"context"
	"lanBox/domain"
	"lanBox/render"
)

type Game struct {
	Cage   domain.Cage
	Boxers domain.Boxers

	Context context.Context
	Cancel  context.CancelFunc
}

func NewGame(width, height int) *Game {
	ctx, cancel := context.WithCancel(context.Background())

	cage := NewCage(width, height)
	boxers := NewBoxers(cage)
	RecalculateCollisionAreas(boxers)

	return &Game{
		Cage:   cage,
		Boxers: boxers,

		Context: ctx,
		Cancel:  cancel,
	}
}

func (g *Game) CloseGame() {
	render.ClearScreen()
	g.Cancel()
}

func (g *Game) RenderConverter() domain.RenderInfo {
	return domain.RenderInfo{
		Cage:   g.Cage,
		Boxers: g.Boxers,
	}
}
