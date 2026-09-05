package render

import (
	"fmt"
	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/sprites"
)

func Snapshot(boxer domain.BaseBoxer) domain.BaseBoxer {
	return domain.BaseBoxer{
		Area:      boxer.Area,
		Name:      boxer.Name,
		Color:     boxer.Color,
		Health:    boxer.Health,
		LastHit:   boxer.LastHit,
		Position:  boxer.Position,
		Direction: boxer.Direction,
		Situation: boxer.Situation,
	}
}

func BoxerFrame(boxer *domain.Boxer, copyBoxer domain.BaseBoxer, parts map[domain.BodyPart][]domain.Direction) {
	color := "\033[31m"
	// Hit effect is beside the boxer body parts
	// so it will not remove the parts if it was a hit
	if copyBoxer.Situation < domain.HeadHit {
		//if sit < HeadHit {
		// Not getting hit
		color = boxer.Color
		ClearBoxer(boxer.BaseBoxer, parts)
	}

	//b.Lock.Lock()
	boxer.BaseBoxer = copyBoxer
	//b.BaseBoxer.Situation = sit
	//b.Lock.Unlock()

	DrawBoxer(boxer.BaseBoxer, parts, color)

	// Move cursor to end
	fmt.Printf("\033[%d;1H", 10)
}

func HitEffect(boxer *domain.Boxer, bodyPart domain.BodyPart, direction domain.Direction) {

	var (
		baseSit     = boxer.Situation
		copyBoxer   = Snapshot(boxer.BaseBoxer)
		affectParts map[domain.BodyPart][]domain.Direction
	)
	// Set lock
	boxer.Lock.Lock()
	switch bodyPart {
	case domain.Head:
		// Init effect
		copyBoxer.Situation = domain.HeadInitHit
		BoxerFrame(boxer, copyBoxer, affectParts)
		Frame(1)
		// Main effect
		copyBoxer.Situation = domain.HeadHit
		BoxerFrame(boxer, copyBoxer, sprites.AllBodyParts)
		Frame(1)
		// Reform Idle
		//sBaseSit = Idle
		copyBoxer.Situation = baseSit
		BoxerFrame(boxer, copyBoxer, sprites.AllBodyParts)
	case domain.Shoulder:
		// Init effect
		copyBoxer.Situation = domain.ShoulderInitHit
		affectParts = map[domain.BodyPart][]domain.Direction{domain.Shoulder: {direction}}
		BoxerFrame(boxer, copyBoxer, affectParts)
		Frame(1)
		// Main effect
		copyBoxer.Situation = domain.ShoulderHit
		affectParts = map[domain.BodyPart][]domain.Direction{domain.Head: {domain.Left}, domain.Arm: {direction}, domain.Shoulder: {direction}}
		BoxerFrame(boxer, copyBoxer, affectParts)
		Frame(1)
		// Reform Idle
		//sBaseSit = Idle
		copyBoxer.Situation = baseSit
		BoxerFrame(boxer, copyBoxer, affectParts)
	case domain.Arm:
	}
	Frame(1)
	// Release lock
	boxer.Lock.Unlock()
}
