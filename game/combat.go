package game

import (
	"lanBox/domain"
	"lanBox/render"
	"lanBox/sprites"
	"time"
	"unicode/utf8"
)

func (g *Game) BoxerPunch(
	attackerID domain.PlayerID,
	direction domain.Direction,
) {
	attacker := g.Boxers[attackerID]
	defender := g.Boxers[attackerID.Opposite()]

	punch := NewPunch(attacker, direction)

	render.PlayPunch(punch)

	hit, part, dir := CheckHit(punch, defender)

	if hit {
		g.HitLogic(part, attackerID.Opposite())
		render.HitEffect(defender, part, dir)
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
		charLen   int
		part      domain.BodyPart
		dir       domain.Direction
		shCharLen = utf8.RuneCountInString(
			sprites.BodyChars[domain.Punch][domain.Shoulder][punch.Attacker.Direction][punch.Direction])
	)

	// calculate punch position
	//hitPoint := sprites.SpriteLayouts[domain.Punch].Parts[domain.Shoulder][punch.Direction]
	hitPoint := sprites.GetPosition(domain.Shoulder, punch.Direction, domain.Punch)
	hitPoint = sprites.CalculatePartPosition(shCharLen, punch.Attacker.BaseBoxer, hitPoint)
	if punch.Attacker.Direction == domain.Left {
		hitPoint.X += shCharLen
	}

	// Compare with defender body parts
	for part = range sprites.SpriteLayouts[domain.Idle].Parts {
		for dir = range sprites.SpriteLayouts[domain.Idle].Parts[part] {
			charLen = utf8.RuneCountInString(sprites.BodyChars[defender.Situation][part][defender.Direction][dir])
			if part != domain.Head &&
				defender.Direction == domain.Right {
				dir = dir.Opposite()
			}
			partPos := sprites.GetPosition(part, dir, domain.Idle)
			switch part {
			case domain.Arm:
				partPos.X -= 1
			case domain.Shoulder:
				partPos.X += 2
			}
			charLen = utf8.RuneCountInString(sprites.GetBodyChar(part, dir, defender.Situation, defender.Direction))
			partPos = sprites.CalculatePartPosition(charLen, defender.BaseBoxer, partPos)
			if defender.Direction == domain.Left { // == Left
				partPos.X += charLen
			}

			// todo add arm hit
			//if part == Arm {
			//	PrintOn(Position{10, 10}, defender.Color, "hitpoint = "+strconv.Itoa(hitPoint.X)+":"+strconv.Itoa(hitPoint.Y))
			//	PrintOn(Position{10, 11}, defender.Color, "partpose = "+strconv.Itoa(partPos.X)+":"+strconv.Itoa(partPos.Y))
			//}

			if hitPoint == partPos {
				if part != domain.Head &&
					defender.Direction == domain.Right {
					dir = dir.Opposite()
				}
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
	if g.Boxers[playerID].Health == 0 {
		// todo: move these to render
		// ex: render.LoseGame()
		render.PrintOn(domain.Position{X: 20, Y: 24}, g.Boxers[playerID].Color, "Loser")
		render.Frame(100)
		g.CloseGame()
		return
	}

	// Set last hit timestamp
	g.Boxers[playerID].LastHit = time.Now().Unix()
	// Print hp
	render.DrawHealth(g.RenderConverter(), playerID)
}
