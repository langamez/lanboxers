package game

import (
	"lanBox/domain"
)

func (g *Game) HandleAction(
	act domain.Act,
	player domain.PlayerID,
) {
	switch act {
	case domain.DoPunch:
		g.BoxerPunch(
			player,
			domain.Left,
		)
	case -domain.DoPunch:
		g.BoxerPunch(
			player,
			domain.Right,
		)
	case domain.DoMoveDown:
		g.BoxerMove(
			player,
			domain.Lower,
		)
	case -domain.DoMoveDown:
		g.BoxerMove(
			player,
			domain.Upper,
		)
	case domain.DoMoveLeft:
		g.BoxerMove(
			player,
			domain.Left,
		)
	case -domain.DoMoveLeft:
		g.BoxerMove(
			player,
			domain.Right,
		)
	case domain.Quit, -domain.Quit:
		g.CloseGame()
	}
}
