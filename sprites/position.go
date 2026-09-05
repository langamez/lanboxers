package sprites

import (
	"unicode/utf8"
	"github.com/langamez/lanboxers/domain"
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

func CalculateHitPoint(
	punch domain.PunchDetail,
) domain.Position {
	var (
		hitPoint domain.Position
	)
	charLen := utf8.RuneCountInString(
		BodyChars[domain.Punch][domain.Shoulder][punch.Attacker.Direction][punch.Direction])
	switch punch.Attacker.Direction {
	case domain.Left:
		hitPoint = GetPosition(domain.Shoulder, punch.Direction, domain.Punch)
		return domain.Position{
			X: punch.Attacker.Position.X + hitPoint.X + 2 + charLen,
			Y: punch.Attacker.Position.Y + hitPoint.Y,
		}
	case domain.Right:
		hitPoint = GetPosition(domain.Shoulder, punch.Direction.Opposite(), domain.Punch)
		return domain.Position{
			X: punch.Attacker.Position.X - hitPoint.X - charLen - 4,
			Y: punch.Attacker.Position.Y + hitPoint.Y,
		}
	}

	return hitPoint
}
