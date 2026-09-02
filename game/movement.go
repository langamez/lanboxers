package game

import (
	"lanBox/domain"
	"lanBox/render"
	"lanBox/sprites"
)

func (g *Game) resolvePass(
	boxerCopy *domain.BaseBoxer,
	opponentCopy *domain.BaseBoxer,
) {
	switch boxerCopy.Direction {
	case domain.Left:
		boxerCopy.Position.X += domain.IdleBoxerForwardLength + 5
		opponentCopy.Position.X -= domain.IdleBoxerForwardLength + 5
	case domain.Right:
		boxerCopy.Position.X -= domain.IdleBoxerForwardLength + 5
		opponentCopy.Position.X += domain.IdleBoxerForwardLength + 5
	}

	//render.PrintOn(domain.Position{X: 10, Y: 10}, boxerCopy.Color, fmt.Sprintf("%d", boxerCopy.Position.X))
	//render.PrintOn(domain.Position{X: 10, Y: 11}, opponentCopy.Color, fmt.Sprintf("%d", opponentCopy.Position.X))
	//render.Frame(20)

	boxerCopy.Direction = boxerCopy.Direction.Opposite()
	opponentCopy.Direction = opponentCopy.Direction.Opposite()
}

func (g *Game) BoxerMove(playerID domain.PlayerID, direction domain.Direction) {
	boxer := g.Boxers[playerID]
	opponent := g.Boxers[playerID.Opposite()]

	boxerCopy := render.Snapshot(boxer.BaseBoxer)
	opponentCopy := render.Snapshot(opponent.BaseBoxer)

	switch direction {
	case domain.Upper:
		boxerCopy.Position.Y--
	case domain.Lower:
		boxerCopy.Position.Y++
	case domain.Left:
		boxerCopy.Position.X--
	case domain.Right:
		boxerCopy.Position.X++
	default:
		return
	}

	if g.CageCollision(boxerCopy) {
		return
	}
	collision, override := g.BoxerCollision(playerID, boxerCopy)
	if collision {
		return
	} else if override {
		g.resolvePass(&boxerCopy, &opponentCopy)
		render.BoxerFrame(opponent, opponentCopy, sprites.AllBodyParts)
	}
	render.BoxerFrame(boxer, boxerCopy, sprites.AllBodyParts)

	// todo: make this cleaner
	if override {
		RecalculateCollisionAreas(g.Boxers)
	}
}
