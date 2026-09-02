package sprites

import (
	"lanBox/domain"
)

func CalculatePartPosition(
	charLen int,
	boxer domain.BaseBoxer,
	partPos domain.Position,
) domain.Position {
	switch boxer.Direction {
	case domain.Left:
		return domain.Position{
			X: boxer.Position.X + partPos.X,
			Y: boxer.Position.Y + partPos.Y,
		}
	case domain.Right:
		return domain.Position{
			X: (boxer.Position.X - partPos.X) - charLen,
			Y: boxer.Position.Y + partPos.Y,
		}
	}

	return partPos
}
