package sprites

import (
	"github.com/langamez/lanboxers/domain"
)

func GetBodyChar(
	part domain.BodyPart,
	partDir domain.Direction,
	situation domain.SituationType,
	boxerDir domain.Direction,
) string {
	return BodyChars[situation][part][boxerDir][partDir]
}

func GetPosition(
	part domain.BodyPart,
	partDir domain.Direction,
	situation domain.SituationType,
) domain.Position {
	return SpriteLayouts[situation].Parts[part][partDir]
}
