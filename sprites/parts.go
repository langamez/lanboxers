package sprites

import "lanBox/domain"

var AllBodyParts = map[domain.BodyPart][]domain.Direction{
	domain.Head:     {domain.Left},
	domain.Shoulder: {domain.Left, domain.Right},
	domain.Arm:      {domain.Left, domain.Right},
}
