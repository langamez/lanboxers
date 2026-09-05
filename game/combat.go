package game

import (
	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/render"
	"github.com/langamez/lanboxers/sprites"
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
		charLen int
		part    domain.BodyPart
		dir     domain.Direction
	)

	// calculate punch position
	hitPoint := sprites.CalculateHitPoint(punch)
	//render.PrintOn(domain.Position{10, 10}, punch.Attacker.Color, fmt.Sprintf("hitpoint = %d:%d", hitPoint.X, hitPoint.Y))
	//render.PrintOn(domain.Position{punch.Attacker.Position.X + 5, punch.Attacker.Position.Y}, punch.Attacker.Color, fmt.Sprintf("boxer position = %d:%d", punch.Attacker.Position.X, punch.Attacker.Position.Y))
	//char := fmt.Sprintf("suboxer position = %d:%d", defender.Position.X, defender.Position.Y)
	//render.PrintOn(domain.Position{defender.Position.X - utf8.RuneCountInString(char) - 5, defender.Position.Y}, defender.Color, fmt.Sprintf("suboxer position = %d:%d", defender.Position.X, defender.Position.Y))

	// Compare with defender body parts
	for part = range sprites.SpriteLayouts[domain.Idle].Parts {
		for dir = range sprites.SpriteLayouts[domain.Idle].Parts[part] {
			charLen = utf8.RuneCountInString(sprites.GetBodyChar(part, dir, defender.Situation, defender.Direction))
			if part != domain.Head &&
				punch.Attacker.Direction == domain.Right {
				//render.PrintOn(domain.Position{10, 13}, defender.Color, fmt.Sprintf("changed from %d dir to: %d", dir, dir.Opposite()))
				dir = dir.Opposite()
			}
			partPos := sprites.GetPosition(part, dir, domain.Idle)
			partPos = sprites.CalculatePartPosition(charLen, defender.BaseBoxer, partPos)

			if part == domain.Shoulder &&
				dir == domain.Right {
				//render.PrintOn(domain.Position{10, 14}, defender.Color, fmt.Sprintf("partpose = %d:%d", partPos.X, partPos.Y))
			}

			//if defender.Direction == domain.Left { // == Left
			//	partPos.X += charLen
			//}

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
	if g.Boxers[playerID].Health == 0 {
		// todo: move these to render
		// ex: render.LoseGame()
		//render.PrintOn(domain.Position{X: 20, Y: 24}, g.Boxers[playerID].Color, "Loser")
		render.Frame(100)
		g.CloseGame()
		return
	}

	// Set last hit timestamp
	g.Boxers[playerID].LastHit = time.Now().Unix()
	// Print hp
	render.DrawHealth(g.RenderConverter(), playerID)
}
