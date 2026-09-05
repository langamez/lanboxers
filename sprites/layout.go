package sprites

import (
	"github.com/langamez/lanboxers/domain"
)

type SpritePosition struct {
	Parts map[domain.BodyPart]map[domain.Direction]domain.Position
}

var SpriteLayouts = map[domain.SituationType]SpritePosition{
	domain.Idle: {
		Parts: map[domain.BodyPart]map[domain.Direction]domain.Position{
			domain.Head:     {domain.Left: {X: 0, Y: 0}},
			domain.Arm:      {domain.Right: {X: 0, Y: +2}, domain.Left: {X: 0, Y: -2}},
			domain.Shoulder: {domain.Right: {X: 0, Y: +1}, domain.Left: {X: 0, Y: -1}},
		},
	},
	domain.Punch: {
		Parts: map[domain.BodyPart]map[domain.Direction]domain.Position{
			domain.Shoulder: {domain.Right: {X: +1, Y: +1}, domain.Left: {X: +1, Y: -1}},
		},
	},
	domain.PunchInit: {
		Parts: map[domain.BodyPart]map[domain.Direction]domain.Position{
			domain.Shoulder: {domain.Right: {X: 0, Y: +1}, domain.Left: {X: 0, Y: -1}},
		},
	},
	domain.HeadHit: {
		Parts: map[domain.BodyPart]map[domain.Direction]domain.Position{
			domain.Head:     {domain.Left: {X: -2, Y: 0}},
			domain.Arm:      {domain.Right: {X: -3, Y: +2}, domain.Left: {X: -3, Y: -2}},
			domain.Shoulder: {domain.Right: {X: -2, Y: +1}, domain.Left: {X: -2, Y: -1}},
		},
	},
	domain.HeadInitHit: {
		Parts: map[domain.BodyPart]map[domain.Direction]domain.Position{
			domain.Head:     {domain.Left: {X: -1, Y: 0}},
			domain.Shoulder: {domain.Right: {X: -2, Y: +1}, domain.Left: {X: -2, Y: -1}},
		},
	},
	domain.ShoulderHit: {
		Parts: map[domain.BodyPart]map[domain.Direction]domain.Position{
			domain.Head:     {domain.Left: {X: -2, Y: 0}},
			domain.Arm:      {domain.Right: {X: -2, Y: +2}, domain.Left: {X: -2, Y: -2}},
			domain.Shoulder: {domain.Right: {X: -2, Y: +1}, domain.Left: {X: -2, Y: -1}},
		},
	},
	domain.ShoulderInitHit: {
		Parts: map[domain.BodyPart]map[domain.Direction]domain.Position{
			domain.Shoulder: {domain.Right: {X: -1, Y: +1}, domain.Left: {X: -1, Y: -1}},
		},
	},
}
