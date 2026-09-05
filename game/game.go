package game

import (
	"context"

	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/render"
)

type Game struct {
	Cage        domain.Cage
	Boxers      domain.Boxers
	RenderChans domain.RenderChannels

	Context context.Context
	Cancel  context.CancelFunc
}

func NewGame(width, height int) *Game {
	ctx, cancel := context.WithCancel(context.Background())

	cage := NewCage(width, height)
	boxers := NewBoxers(cage)
	renderChans := NewRenderChannels()
	RecalculateCollisionAreas(boxers)

	return &Game{
		Cage:        cage,
		Boxers:      boxers,
		RenderChans: renderChans,

		Context: ctx,
		Cancel:  cancel,
	}
}

func (g *Game) CloseGame() {
	render.ClearScreen()
	g.Cancel()
}

func (g *Game) LoseGame(playerID domain.PlayerID) {
	render.LoseEffect(g.Boxers[playerID].Color, g.Cage.Area.Max)
	g.CloseGame()
}

func (g *Game) RenderConverter() domain.RenderInfo {
	return domain.RenderInfo{
		Cage:   g.Cage,
		Boxers: g.Boxers,
	}
}
