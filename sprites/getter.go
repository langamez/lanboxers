package sprites

import (
	"github.com/langamez/lanboxers/domain"
)

func GetBodyChar(
	part domain.BodyPart,
	partDir domain.Direction,
	situation domain.Situation,
	boxerDir domain.Direction,
) string {
	return BodyChars[situation][part][boxerDir][partDir]
}

func GetPosition(
	part domain.BodyPart,
	partDir domain.Direction,
	situation domain.Situation,
) domain.Position {
	return SpriteLayouts[situation].Parts[part][partDir]
}
