package game

import (
	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/sprites"
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
	boxerCopy.Direction = boxerCopy.Direction.Opposite()
	opponentCopy.Direction = opponentCopy.Direction.Opposite()
}

func (g *Game) BoxerMove(playerID domain.PlayerID, direction domain.Direction) {
	boxer := g.Boxers[playerID]
	opponent := g.Boxers[playerID.Opposite()]

	boxerCopy := Snapshot(boxer.BaseBoxer)
	opponentCopy := Snapshot(opponent.BaseBoxer)

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
		g.UpdateBoxer(playerID.Opposite(), opponentCopy, sprites.AllBodyParts)
	}
	g.UpdateBoxer(playerID, boxerCopy, sprites.AllBodyParts)

	// todo: make this cleaner
	if override {
		RecalculateCollisionAreas(g.Boxers)
	}
}
