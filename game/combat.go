package game

import (
	"time"
	"unicode/utf8"

	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/render"
	"github.com/langamez/lanboxers/sprites"
)

func (g *Game) BoxerPunch(
	attackerID domain.PlayerID,
	direction domain.Direction,
) {
	attacker := g.Boxers[attackerID]
	defender := g.Boxers[attackerID.Opposite()]

	mCopyBoxer := Snapshot(attacker.BaseBoxer)
	affectedParts := map[domain.BodyPart][]domain.Direction{domain.Shoulder: {direction}, domain.Arm: {direction}}

	// Init
	mCopyBoxer.Situation = domain.Situation{SituationType: domain.PunchInit, Direction: direction}
	g.UpdateBoxer(attackerID, mCopyBoxer, affectedParts)
	// Main effect
	mCopyBoxer.Situation.SituationType = domain.Punch
	g.UpdateBoxer(attackerID, mCopyBoxer, affectedParts)
	// Form back to Idle
	mCopyBoxer.Situation.SituationType = domain.Idle
	g.UpdateBoxer(attackerID, mCopyBoxer, affectedParts)

	punch := NewPunch(attacker, direction)
	hit, part, dir := CheckHit(punch, defender)

	if hit {
		g.BoxerHit(part, attackerID.Opposite(), dir)
	}
}

func NewPunch(
	boxer *domain.Boxer,
	direction domain.Direction,
) domain.PunchDetail {
	return domain.PunchDetail{
		Attacker:  boxer,
		Direction: direction,
	}
}

func (g *Game) HitLogic(
	part domain.BodyPart,
	playerID domain.PlayerID,
) {
	// todo:
	// Hit effect
	// Power hit
	//if pw {
	//}

	if part != domain.Arm {
		// Lower hp
		g.TakeDamage(playerID, part)
	}
}

func CheckHit(
	punch domain.PunchDetail,
	defender *domain.Boxer,
) (bool, domain.BodyPart, domain.Direction) {
	var (
		charLen int
		part    domain.BodyPart
		dir     domain.Direction
	)

	// calculate punch position
	hitPoint := sprites.CalculateHitPoint(punch)
	// Compare with defender body parts
	for part = range sprites.SpriteLayouts[domain.Idle].Parts {
		for dir = range sprites.SpriteLayouts[domain.Idle].Parts[part] {
			charLen = utf8.RuneCountInString(sprites.GetBodyChar(part, dir, defender.Situation.SituationType, defender.Direction))
			partPos := sprites.GetPosition(part, dir, domain.Idle)
			partPos = sprites.CalculatePartPosition(charLen, defender.BaseBoxer, partPos)

			if defender.Direction == domain.Right {
				dir = dir.Opposite()
			}

			// todo add arm hit
			//if part == Arm {
			//	PrintOn(Position{10, 10}, defender.Color, "hitpoint = "+strconv.Itoa(hitPoint.X)+":"+strconv.Itoa(hitPoint.Y))
			//	PrintOn(Position{10, 11}, defender.Color, "partpose = "+strconv.Itoa(partPos.X)+":"+strconv.Itoa(partPos.Y))
			//}

			if hitPoint == partPos {
				//if part != domain.Head &&
				//	defender.Direction == domain.Right {
				//	dir = dir.Opposite()
				//}
				//render.PrintOn(domain.Position{10, 14}, defender.Color, fmt.Sprintf("direction = %d", dir))
				return true, part, dir
			}
		}
	}
	return false, part, dir
}

func (g *Game) TakeDamage(
	playerID domain.PlayerID,
	bPart domain.BodyPart,
) {
	var (
		amount int
	)
	switch bPart {
	case domain.Head:
		amount = domain.HeadHitHP
	case domain.Shoulder:
		amount = domain.ShoulderHitHP
	}
	g.Boxers[playerID].AddHealth(-amount)
	if g.Boxers[playerID].Health <= 0 {
		g.LoseGame(playerID)
		return
	}
	// Set last hit timestamp
	g.Boxers[playerID].LastHit = time.Now().Unix()
	// Print hp
	render.DrawHealth(g.RenderConverter(), playerID)
}

func (g *Game) BoxerHit(
	part domain.BodyPart,
	id domain.PlayerID,
	direction domain.Direction,
) {
	var (
		boxer = g.Boxers[id]
		//baseSit     = boxer.Situation
		copyBoxer     = Snapshot(boxer.BaseBoxer)
		affectedParts map[domain.BodyPart][]domain.Direction
	)

	g.HitLogic(part, id)
	switch part {
	case domain.Head:
		// Init effect
		copyBoxer.Situation = domain.Situation{SituationType: domain.HeadInitHit}
		affectedParts = map[domain.BodyPart][]domain.Direction{
			domain.Head: {domain.Left}, domain.Shoulder: {direction, direction.Opposite()},
		}
		g.UpdateBoxer(id, copyBoxer, affectedParts)
		// Main effect
		copyBoxer.Situation.SituationType = domain.HeadHit
		affectedParts = sprites.AllBodyParts
		g.UpdateBoxer(id, copyBoxer, affectedParts)
		// Reform Idle
		//copyBoxer.Situation = baseSit
		copyBoxer.Situation.SituationType = domain.Idle
		g.UpdateBoxer(id, copyBoxer, affectedParts)
	case domain.Shoulder:
		// Init effect
		copyBoxer.Situation = domain.Situation{SituationType: domain.ShoulderInitHit, Direction: direction}
		affectedParts = map[domain.BodyPart][]domain.Direction{domain.Shoulder: {direction}}
		g.UpdateBoxer(id, copyBoxer, affectedParts)
		// Main effect
		copyBoxer.Situation.SituationType = domain.ShoulderHit
		affectedParts = map[domain.BodyPart][]domain.Direction{
			domain.Head: {domain.Left}, domain.Arm: {direction}, domain.Shoulder: {direction},
		}
		g.UpdateBoxer(id, copyBoxer, affectedParts)
		// Reform Idle
		//copyBoxer.Situation = baseSit
		copyBoxer.Situation.SituationType = domain.Idle
		g.UpdateBoxer(id, copyBoxer, affectedParts)
	case domain.Arm:
	}
}
