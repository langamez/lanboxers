package game

import (
	"github.com/langamez/lanboxers/domain"
)

func RecalculateCollisionAreas(boxers domain.Boxers) {
	for _, b := range boxers {
		switch b.Direction {
		case domain.Left:
			b.Area.Max = domain.Position{
				X: domain.IdleBoxerForwardLength,
				Y: domain.IdleBoxerWidth,
			}

			b.Area.Min = domain.Position{
				X: domain.IdleBoxerBehindLength,
				Y: -domain.IdleBoxerWidth,
			}

		case domain.Right:
			b.Area.Max = domain.Position{
				X: domain.IdleBoxerBehindLength - 1,
				Y: domain.IdleBoxerWidth,
			}

			b.Area.Min = domain.Position{
				X: -domain.IdleBoxerForwardLength + 1,
				Y: -domain.IdleBoxerWidth,
			}
		}
	}
}

func (g *Game) CageCollision(boxer domain.BaseBoxer) bool {
	boxerMin := boxer.Area.Min.Add(boxer.Position)
	boxerMax := boxer.Area.Max.Add(boxer.Position)

	cageMin := g.Cage.Area.Min
	cageMax := g.Cage.Area.Max

	return boxerMin.X <= cageMin.X ||
		boxerMax.X >= cageMax.X ||
		boxerMin.Y <= cageMin.Y ||
		boxerMax.Y >= cageMax.Y
}

func verticalOverlap(
	minA, maxA,
	minB, maxB domain.Position,
) bool {
	return !(minA.Y <= maxB.Y &&
		maxA.Y >= minB.Y)
}

func horizontalCollision(
	direction domain.Direction,
	minA, maxA,
	minB, maxB domain.Position,
) bool {
	switch direction {
	case domain.Right:
		return minA.X <= maxB.X
	case domain.Left:
		return maxA.X >= minB.X
	default:
		return false
	}
}

func (g *Game) BoxerCollision(
	playerID domain.PlayerID,
	boxer domain.BaseBoxer,
) (bool, bool) {
	other := g.Boxers[playerID.Opposite()]

	boxerMin := boxer.Area.Min.Add(boxer.Position)
	boxerMax := boxer.Area.Max.Add(boxer.Position)

	otherMin := other.Area.Min.Add(other.Position)
	otherMax := other.Area.Max.Add(other.Position)

	if horizontalCollision(boxer.Direction, boxerMin, boxerMax, otherMin, otherMax) {
		if verticalOverlap(boxerMin, boxerMax, otherMin, otherMax) {
			// override
			return false, true
		}
		return true, false
	}

	return false, false
}
