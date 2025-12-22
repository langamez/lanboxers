package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

func PrintOn(position Position, color, char string) {
	var (
		reset = "\033[0m"
	)
	fmt.Printf("\033[%d;%dH%s%s%s", position.Y, position.X, color, char, reset)
}

func hitCheck(subBoxer Boxer, hitPoint Position) (bool, BodyPart, VerDirection) {
	var (
		charLen int
		part    BodyPart
		vDir    VerDirection
	)
	// Check for hit
	for part = range Positions[Idle] {
		for vDir = range Positions[Idle][part] {
			partPos := Positions[Idle][part][vDir]
			if part == Shoulder {
				partPos.X += 2
			}
			charLen = utf8.RuneCountInString(BodyChars[subBoxer.Situation][part][subBoxer.Direction][vDir])
			partPos.CalPartPos(subBoxer, charLen)
			if subBoxer.Direction == Left {
				partPos.X += charLen
			}
			if hitPoint == partPos {
				return true, part, vDir
			}
		}
	}
	return false, part, vDir
}

func RemoveBoxer(boxer Boxer, parts map[BodyPart][]VerDirection) {
	var (
		char     string
		charLen  int
		position Position
	)

	for part := range parts {
		for _, pDir := range parts[part] {
			position = Positions[boxer.Situation][part][pDir]
			charLen = utf8.RuneCountInString(BodyChars[boxer.Situation][part][boxer.Direction][pDir])
			char = strings.Repeat(" ", charLen)
			position.CalPartPos(boxer, charLen)
			PrintOn(position, "", char)
		}
	}
}

func PrintBoxer(boxer Boxer, parts map[BodyPart][]VerDirection, color string) {
	var (
		char     string
		position Position
	)
	for part := range parts {
		for _, pDir := range parts[part] {
			position = Positions[boxer.Situation][part][pDir]
			char = BodyChars[boxer.Situation][part][boxer.Direction][pDir]
			position.CalPartPos(boxer, utf8.RuneCountInString(char))
			PrintOn(position, color, char)
		}
	}
}

func HitEffect(boxer *Boxer, bodyPart BodyPart, direction VerDirection) map[BodyPart][]VerDirection {

	var (
		copyBoxer   = *boxer
		affectParts map[BodyPart][]VerDirection
	)

	switch bodyPart {
	case Head:
		//Init effect
		copyBoxer.Situation = HeadInitHit
		affectParts = map[BodyPart][]VerDirection{Head: {Up}, Arm: {Up, Down}, Shoulder: {Up, Down}}
		boxer.boxerFrame(copyBoxer, affectParts)
		//Sleep to avoid blank term
		frame(1)
		//Main effect
		copyBoxer.Situation = HeadHit
		boxer.boxerFrame(copyBoxer, AllBodyParts)
	case Shoulder:
		//Init effect
		copyBoxer.Situation = ShoulderInitHit
		affectParts = map[BodyPart][]VerDirection{Shoulder: {direction}}
		boxer.boxerFrame(copyBoxer, affectParts)
		//Sleep to avoid blank term
		frame(1)
		//Main effect
		copyBoxer.Situation = ShoulderHit
		affectParts = map[BodyPart][]VerDirection{Head: {Up}, Arm: {direction}, Shoulder: {direction}}
		boxer.boxerFrame(copyBoxer, affectParts)
	}
	return affectParts
}

func PunchEffect(boxer *Boxer, direction VerDirection) (Position, map[BodyPart][]VerDirection) {
	var (
		copyBoxer   = *boxer
		hitPoint    = Positions[Punch][Shoulder][direction]
		affectParts = map[BodyPart][]VerDirection{Shoulder: {direction}, Arm: {direction}}
		charLen     = utf8.RuneCountInString(BodyChars[Punch][Shoulder][boxer.Direction][direction])
	)

	//Init effect
	copyBoxer.Situation = PunchInit
	boxer.boxerFrame(copyBoxer, affectParts)
	frame(1)
	//Sleep to avoid blank term
	//Main effect
	copyBoxer.Situation = Punch
	boxer.boxerFrame(copyBoxer, affectParts)
	//Find hit position
	hitPoint.CalPartPos(*boxer, charLen)
	if boxer.Direction == Left {
		hitPoint.X += charLen
	}
	return hitPoint, affectParts
}

func (p *Position) CalPartPos(boxer Boxer, charLen int) {
	switch boxer.Direction {
	case Left:
		p.Y = boxer.Position.Y + p.Y
		p.X = boxer.Position.X + p.X
	case Right:
		p.Y = boxer.Position.Y + p.Y
		p.X = (boxer.Position.X - p.X) - charLen
	}
}

func frame(n float64) {
	n += 3
	n = n * 10
	// Move cursor to end
	fmt.Printf("\033[%d;1H", 10)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func (g GameInfo) BoxerMove(direction interface{}, opponent IsOpponent) {
	var (
		boxers    = g.Boxers
		mainBoxer = boxers[opponent]
		copyBoxer = *boxers[opponent]
		parts     = map[BodyPart][]VerDirection{Head: {Up}, Shoulder: {Up, Down}, Arm: {Up, Down}}
	)

	switch direction.(type) {
	case VerDirection:
		switch direction {
		case Up:
			copyBoxer.Position.Y--
		case Down:
			copyBoxer.Position.Y++
		}
	case HorDirection:
		switch direction {
		case Left:
			copyBoxer.Position.X--
		case Right:
			copyBoxer.Position.X++
		}
	}
	if !g.cageCollide(copyBoxer) {
		if !g.boxersCollide(opponent, &copyBoxer, parts) {
			mainBoxer.boxerFrame(copyBoxer, parts)
			// Move cursor to end
			fmt.Printf("\033[%d;1H", 10)
		}
	}
}

func (g GameInfo) BoxerPunch(direction VerDirection, opponent IsOpponent) {
	var (
		gotHit         bool
		part           BodyPart
		sAffectedParts map[BodyPart][]VerDirection
		mBoxer         = g.Boxers[opponent]
		sBoxer         = g.Boxers[!opponent]
		mCopyBoxer     = *g.Boxers[opponent]
		sCopyBoxer     Boxer
	)
	// Punch effect and get the hit point
	hitPoint, mAffectedParts := PunchEffect(mBoxer, direction)
	// Check for subject boxer hit
	gotHit, part, direction = hitCheck(*sBoxer, hitPoint)
	if gotHit {
		// Hit effect
		sAffectedParts = HitEffect(sBoxer, part, direction)
		// Lower hp
		sBoxer.TakeDamage(part)
		// Print hp
		g.PrintHealth(!opponent)
		// Set last hit timestamp
		sBoxer.LastHit = time.Now().Unix()
	}
	//Sleep to avoid blank term
	frame(1)
	// Form back to Idle
	mBoxer.boxerFrame(mCopyBoxer, mAffectedParts)

	if gotHit {
		sCopyBoxer = *g.Boxers[!opponent]
		sCopyBoxer.Situation = Idle
		sBoxer.boxerFrame(sCopyBoxer, sAffectedParts)
	}
	// Move cursor to end
	fmt.Printf("\033[%d;1H", 10)
}

func (g GameInfo) cageCollide(boxer Boxer) bool {
	// Get boxer related areas
	bMin := *boxer.Area[Min]
	bMax := *boxer.Area[Max]
	// Get cage areas
	cMin := g.Cage.Area[Min]
	cMax := g.Cage.Area[Max]
	// Calculate boxer real areas
	bMin.AddPosition(boxer.Position)
	bMax.AddPosition(boxer.Position)
	// Check for Collide
	if bMin.X <= cMin.X || bMax.X >= cMax.X ||
		bMin.Y <= cMin.Y || bMax.Y >= cMax.Y {
		return true
	}
	return false
}

func (g GameInfo) boxersCollide(opponent IsOpponent, mBoxer *Boxer, parts map[BodyPart][]VerDirection) bool {
	var (
		collide = false
		sBoxer  = g.Boxers[!opponent]
		// Get boxers related areas
		// Main boxer
		mMin = *mBoxer.Area[Min]
		mMax = *mBoxer.Area[Max]
		// Subject boxer
		sMin = *sBoxer.Area[Min]
		sMax = *sBoxer.Area[Max]
	)
	// Calculate boxer real areas
	mMin.AddPosition(mBoxer.Position)
	mMax.AddPosition(mBoxer.Position)
	sMin.AddPosition(sBoxer.Position)
	sMax.AddPosition(sBoxer.Position)
	// Check for collide
	if mBoxer.Direction == Left {
		if mMax.X >= sMin.X {
			if mMin.Y >= sMin.Y && mMin.Y <= sMax.Y ||
				mMax.Y >= sMin.Y && mMax.Y <= sMax.Y {
				collide = true
			}
			// Check for override
			if mMin.X >= sMax.X {
				if mMin.Y > sMax.Y || mMax.Y < sMin.Y {
					// Override
					collide = false
					// Make a copy to change
					sCopyBoxer := *g.Boxers[!opponent]
					// Set new position based on direction
					mBoxer.Position.X += IdleBoxerForwardLength + 1
					sCopyBoxer.Position.X -= IdleBoxerForwardLength + 1
					// Set direction
					mBoxer.Direction, sCopyBoxer.Direction = sBoxer.Direction, mBoxer.Direction
					// Print boxer
					sBoxer.boxerFrame(sCopyBoxer, parts)
					// Calculate boxer areas again
					g.CalArea()
				}
			}
		}
	} else if mBoxer.Direction == Right {
		if mMin.X <= sMax.X {
			if mMin.Y >= sMin.Y && mMin.Y <= sMax.Y ||
				mMax.Y >= sMin.Y && mMax.Y <= sMax.Y {
				// Collide
				collide = true
			}
			// Check for override
			if mMin.X <= sMin.X {
				if mMin.Y > sMax.Y || mMax.Y < sMin.Y {
					// Override
					collide = false
					// Make a copy to change
					sCopyBoxer := *g.Boxers[!opponent]
					// Set new position based on direction
					mBoxer.Position.X -= IdleBoxerForwardLength
					sCopyBoxer.Position.X += IdleBoxerForwardLength
					// Set direction
					mBoxer.Direction, sCopyBoxer.Direction = sBoxer.Direction, mBoxer.Direction
					// Print boxer
					sBoxer.boxerFrame(sCopyBoxer, parts)
					// Calculate boxer areas again
					g.CalArea()
				}
			}
		}
	}

	return collide
}

// CalArea : Calculate the Area that each boxer will occupy according to it's position
func (g GameInfo) CalArea() {
	for _, boxer := range g.Boxers {
		switch boxer.Direction {
		case Left:
			boxer.Area[Max] = &Position{IdleBoxerForwardLength, IdleBoxerWidth}
			boxer.Area[Min] = &Position{IdleBoxerBehindLength, -IdleBoxerWidth}
		case Right:
			boxer.Area[Max] = &Position{IdleBoxerBehindLength - 1, IdleBoxerWidth}
			boxer.Area[Min] = &Position{-IdleBoxerForwardLength + 1, -IdleBoxerWidth}
		}
	}
}

func (b *Boxer) boxerFrame(copyBoxer Boxer, parts map[BodyPart][]VerDirection) {
	if copyBoxer.Situation < 4 {
		// Not getting hit
		RemoveBoxer(*b, parts)
		*b = copyBoxer
		PrintBoxer(*b, parts, b.Color)
	} else {
		// Gotten hit
		// Hit effect doesn't affect on base parts
		*b = copyBoxer
		PrintBoxer(*b, parts, "\033[31m")
	}
}

func (p *Position) AddPosition(a Position) {
	p.X += a.X
	p.Y += a.Y
}

func (b *Boxer) TakeDamage(bPart BodyPart) {
	var amount int
	switch bPart {
	case Head:
		amount = HeadHitHp
	case Shoulder:
		amount = ShoulderHitHp
	}
	PrintOn(Position{20, 21}, b.Color, fmt.Sprintf("b.Health: %d", b.Health))
	b.Health -= amount
	if b.Health < 0 {
		// Change boxer Situation
		clearScreen()
		b.Health = 0
	}
}

func (b *Boxer) Heal(amount int) {
	b.Health += amount
	if b.Health > MaxHealthPoint {
		b.Health = MaxHealthPoint
	}
}

func (g GameInfo) PrintHud() {
	//Print name
	g.PrintName(Main)
	g.PrintName(Opponent)
	//Print hp
	g.PrintHealth(Main)
	g.PrintHealth(Opponent)
}

func (g GameInfo) PrintHealth(opponent IsOpponent) {
	var (
		hpChar string
		hpPos  = Position{0, g.Cage.Area[Min].Y + 2}
	)
	//Get cage length
	cageLength := g.Cage.Area[Max].X - g.Cage.Area[Min].X
	//Calculate each player health bar length
	hpLength := cageLength / 3
	//Calculate space between health bar
	indent := hpLength - (2 * WallLength)
	//Convert to character
	filled := g.Boxers[opponent].Health * hpLength / 100
	hpChar = strings.Repeat("█", filled) +
		strings.Repeat("░", hpLength-filled)
	//Calculate horizontal health bar position
	if !opponent {
		hpPos.X = g.Cage.Area[Min].X + WallLength
	} else {
		hpPos.X = g.Cage.Area[Min].X + WallLength + hpLength + indent
	}
	//Print
	PrintOn(hpPos, g.Boxers[opponent].Color, hpChar)
}

func (g GameInfo) PrintName(opponent IsOpponent) {
	var namePos = Position{0, g.Cage.Area[Min].Y + 1}
	//Get cage length
	cageLength := g.Cage.Area[Max].X - g.Cage.Area[Min].X
	//Calculate each player health bar length
	hpLength := cageLength / 3
	//Calculate space between health bar
	indent := hpLength - (2 * WallLength)
	//Calculate horizontal health bar position
	if !opponent {
		namePos.X = g.Cage.Area[Min].X + WallLength
	} else {
		namePos.X = g.Cage.Area[Min].X + WallLength + hpLength + indent
	}
	//Print
	PrintOn(namePos, g.Boxers[opponent].Color, g.Boxers[opponent].Name)
}
