package main

import (
	"context"
	"fmt"
	"os"
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

func hitCheck(subBoxer BaseBoxer, hitPoint Position) (bool, BodyPart, HorDirection) {
	var (
		charLen int
		part    BodyPart
		dir     HorDirection
	)
	// Check for hit
	for part = range Positions[Idle] {
		for dir = range Positions[Idle][part] {
			charLen = utf8.RuneCountInString(BodyChars[subBoxer.Situation][part][subBoxer.Direction][dir])
			if part != Head &&
				subBoxer.Direction { // == right
				dir = !dir
			}
			partPos := Positions[Idle][part][dir]
			if part == Shoulder {
				partPos.X += 2
			}
			charLen = utf8.RuneCountInString(BodyChars[subBoxer.Situation][part][subBoxer.Direction][dir])
			partPos.CalPartPos(subBoxer, charLen)
			if !subBoxer.Direction { // == Left
				partPos.X += charLen
			}

			if hitPoint == partPos {
				if part != Head &&
					subBoxer.Direction { // == right
					dir = !dir
				}
				return true, part, dir
			}
		}
	}
	return false, part, dir
}

func RemoveBoxer(boxer BaseBoxer, parts map[BodyPart][]HorDirection) {
	var (
		char     string
		charLen  int
		position Position
	)

	for part := range parts {
		for _, dir := range parts[part] {
			position = Positions[boxer.Situation][part][dir]
			charLen = utf8.RuneCountInString(BodyChars[boxer.Situation][part][boxer.Direction][dir])
			char = strings.Repeat(" ", charLen)
			// Check if it's right boxer change the direction to opposite
			if part != Head &&
				boxer.Direction { // == Right
				position = Positions[boxer.Situation][part][!dir]
			}
			position.CalPartPos(boxer, charLen)
			PrintOn(position, "", char)
		}
	}
}

func PrintBoxer(boxer BaseBoxer, parts map[BodyPart][]HorDirection, color string) {
	var (
		char     string
		position Position
	)
	for part := range parts {
		for _, dir := range parts[part] {
			position = Positions[boxer.Situation][part][dir]
			char = BodyChars[boxer.Situation][part][boxer.Direction][dir]
			// Check if it's right boxer change the direction to opposite
			if part != Head &&
				boxer.Direction { // == Right
				position = Positions[boxer.Situation][part][!dir]
			}
			position.CalPartPos(boxer, utf8.RuneCountInString(char))
			PrintOn(position, color, char)
		}
	}
}

func HitEffect(boxer *Boxer, bodyPart BodyPart, direction HorDirection) {

	var (
		baseSit     = boxer.Situation
		copyBoxer   = boxer.Snapshot()
		affectParts map[BodyPart][]HorDirection
	)

	switch bodyPart {
	case Head:
		// Init effect
		// Set lock
		boxer.Lock.Lock()
		copyBoxer.Situation = HeadInitHit
		affectParts = map[BodyPart][]HorDirection{Head: {Left}, Arm: {Left, Right}, Shoulder: {Left, Right}}
		boxer.boxerFrame(copyBoxer, affectParts)
		frame(1)
		// Main effect
		copyBoxer.Situation = HeadHit
		boxer.boxerFrame(copyBoxer, AllBodyParts)
		frame(1)
		// Reform Idle
		//sBaseSit = Idle
		copyBoxer.Situation = baseSit
		boxer.boxerFrame(copyBoxer, AllBodyParts)
		// Release lock
		boxer.Lock.Unlock()
	case Shoulder:
		// Init effect
		// Set lock
		boxer.Lock.Lock()
		copyBoxer.Situation = ShoulderInitHit
		affectParts = map[BodyPart][]HorDirection{Shoulder: {direction}}
		boxer.boxerFrame(copyBoxer, affectParts)
		frame(1)
		// Main effect
		copyBoxer.Situation = ShoulderHit
		affectParts = map[BodyPart][]HorDirection{Head: {Left}, Arm: {direction}, Shoulder: {direction}}
		boxer.boxerFrame(copyBoxer, affectParts)
		frame(1)
		// Reform Idle
		//sBaseSit = Idle
		copyBoxer.Situation = baseSit
		boxer.boxerFrame(copyBoxer, affectParts)
		// Release lock
		boxer.Lock.Unlock()
	}
	frame(1)
}

func (p *Position) CalPartPos(boxer BaseBoxer, charLen int) {
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
		baseDir   = mainBoxer.Direction
		copyBoxer = boxers[opponent].Snapshot()
		parts     = map[BodyPart][]HorDirection{Head: {Left}, Shoulder: {Left, Right}, Arm: {Left, Right}}
	)

	switch direction.(type) {
	case VerDirection:
		switch direction {
		case Upper:
			copyBoxer.Position.Y--
		case Lower:
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
			if baseDir != mainBoxer.Direction {
				g.CalArea()
			} // Override: Calculate boxer areas again
			fmt.Printf("\033[%d;1H", 10) // Move cursor to end
		}
	}

	if !g.cageCollide(copyBoxer) {
		if !g.boxersCollide(opponent, &copyBoxer, parts) {
			mainBoxer.boxerFrame(copyBoxer, parts)
			if baseDir != mainBoxer.Direction {
				g.CalArea()
			} // Override: Calculate boxer areas again
			fmt.Printf("\033[%d;1H", 10) // Move cursor to end
		}
	}
}

func (g GameInfo) BoxerPunch(direction HorDirection, opponent IsOpponent) {
	var (
		gotHit         bool
		pwPunch        bool
		part           BodyPart
		mBoxer         = g.Boxers[opponent]
		sBoxer         = g.Boxers[!opponent]
		mCopyBoxer     = g.Boxers[opponent].Snapshot()
		mAffectedParts = map[BodyPart][]HorDirection{Shoulder: {direction}, Arm: {direction}}
		shCharLen      = utf8.RuneCountInString(BodyChars[Punch][Shoulder][mBoxer.Direction][direction])
	)
	// Set punch direction
	if mBoxer.Direction {
		direction = !direction
	}
	mCopyBoxer.SituationDir = direction
	// Check for subject boxer situation
	if sBoxer.SituationDir != mBoxer.SituationDir {
		switch sBoxer.Situation {
		case Punch:
			// I guess sub: power shot
			// sBoxer: punch
			// mBoxer: idle to punch init
			pwPunch = true
			//mCopyBoxer.Situation = PwHit
			PrintOn(Position{20, 24}, mBoxer.Color, "got power shot")
		case PunchInit:
			// I guess: collapse
			// sBoxer: punch init
			// mBoxer: idle to punch init
			//mCopyBoxer.Situation = Collapse
			PrintOn(Position{20, 25}, sBoxer.Color, "collapse")
		default:
			mCopyBoxer.Situation = PunchInit
		}
	} else {
		mCopyBoxer.Situation = PunchInit
	}
	// Init effect
	mBoxer.boxerFrame(mCopyBoxer, mAffectedParts)
	frame(1)
	PrintOn(Position{10, 10}, mBoxer.Color, "main punch init")
	// Check for collapse
	//if mCopyBoxer.Situation != Collapse {
	// Check for subject boxer situation
	if sBoxer.Situation != Punch {
		// Main effect
		mCopyBoxer.Situation = Punch
		mBoxer.boxerFrame(mCopyBoxer, mAffectedParts)
		PrintOn(Position{10, 11}, mBoxer.Color, "main punch")
		//Find hit position
		//hitPoint := PunchEffect(mBoxer, direction)
		if sBoxer.Situation == PunchInit &&
			sBoxer.SituationDir != mBoxer.SituationDir {
			// hit with power shot
			// sBoxer: punch init
			// mBoxer: punch init to punch
			pwPunch = true
			PrintOn(Position{20, 26}, mBoxer.Color, "power punch")
		}

		// Get punch position
		hitPoint := Positions[Punch][Shoulder][direction]
		// Form based on boxer position to get real position
		hitPoint.CalPartPos(*mBoxer.BaseBoxer, shCharLen)
		if mBoxer.Direction == Left {
			hitPoint.X += shCharLen
		}
		// Check for hit
		gotHit, part, direction = hitCheck(*sBoxer.BaseBoxer, hitPoint)
		if gotHit {
			g.HitLogic(!opponent, part, direction, pwPunch)
		}
	} else {
		PrintOn(Position{10, 11}, mBoxer.Color, "s in punch")
	}
	//else {
	//	if sBoxer.VerticalSituation == mBoxer.VerticalSituation {
	//		// Collapse
	//		PrintOn(Position{20, 27}, sBoxer.Color, "got hit")
	//	}
	//}

	// Form back to Idle
	mCopyBoxer.Situation = Idle
	mBoxer.boxerFrame(mCopyBoxer, mAffectedParts)

	//}

	if gotHit {

		PrintOn(Position{20, 26}, sBoxer.Color, "to default")
	}
	// Move cursor to end
	fmt.Printf("\033[%d;1H", 10)
}

func (g GameInfo) cageCollide(boxer BaseBoxer) bool {
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

func (g GameInfo) boxersCollide(opponent IsOpponent, mBoxer *BaseBoxer, parts map[BodyPart][]HorDirection) bool {
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
	if mBoxer.Direction { //  == Right
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
					sCopyBoxer := g.Boxers[!opponent].Snapshot()
					// Set new position based on direction
					mBoxer.Position.X -= IdleBoxerForwardLength
					sCopyBoxer.Position.X += IdleBoxerForwardLength
					// Set direction
					g.Boxers[opponent].Direction, sCopyBoxer.Direction = sBoxer.Direction, mBoxer.Direction
					// Print boxer
					sBoxer.boxerFrame(sCopyBoxer, parts)
					// Calculate boxer areas again
					g.CalArea()
				}
			}
		}
	} else { //  == Left
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
					sCopyBoxer := g.Boxers[!opponent].Snapshot()
					// Set new position based on direction
					mBoxer.Position.X += IdleBoxerForwardLength + 1
					sCopyBoxer.Position.X -= IdleBoxerForwardLength + 1
					// Set direction
					mBoxer.Direction, sCopyBoxer.Direction = sBoxer.Direction, mBoxer.Direction
					// Print boxer
					sBoxer.boxerFrame(sCopyBoxer, parts)
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

func (b *Boxer) boxerFrame(copyBoxer BaseBoxer, parts map[BodyPart][]HorDirection) {
	color := "\033[31m"
	// Hit effect is beside the boxer body parts
	// so it will not remove the parts if it was a hit
	if copyBoxer.Situation < HeadHit {
		//if sit < HeadHit {
		// Not getting hit
		color = b.Color
		RemoveBoxer(*b.BaseBoxer, parts)
	}

	//b.Lock.Lock()
	*b.BaseBoxer = copyBoxer
	//b.BaseBoxer.Situation = sit
	//b.Lock.Unlock()

	PrintBoxer(*b.BaseBoxer, parts, color)
}

func (p *Position) AddPosition(a Position) {
	p.X += a.X
	p.Y += a.Y
}

func (g GameInfo) TakeDamage(opp IsOpponent, bPart BodyPart) {
	var (
		amount int
		boxer  = g.Boxers[opp]
	)
	switch bPart {
	case Head:
		amount = HeadHitHp
	case Shoulder:
		amount = ShoulderHitHp
	}
	boxer.Health -= amount
	// Set last hit timestamp
	boxer.LastHit = time.Now().Unix()
	// Print hp
	g.PrintHealth(opp)
	if !(boxer.Health > 0) {
		// Change boxer Situation
		clearScreen()
		PrintOn(Position{20, 24}, boxer.Color, "Loser")
		frame(100)
		g.Cancel()
		//boxer.Health = 0
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

func EventDispatcher(in <-chan Event, router map[IsOpponent]chan Event) {
	for e := range in {
		if ch, ok := router[e.Actor]; ok {
			ch <- e
		}
	}
}

func ReadKeyboard(cancel context.CancelFunc, eventChan chan<- Event) {
	buf := make([]byte, 3)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}
		mainEvent := Event{Actor: Main}
		oppEvent := Event{Actor: Opponent}

		switch string(buf[:n]) {
		// Boxer 1
		// Move
		case "w": // up
			mainEvent.Act = -DoMoveDown
			eventChan <- mainEvent
		case "s": // down
			mainEvent.Act = DoMoveDown
			eventChan <- mainEvent
		case "d": // right
			mainEvent.Act = -DoMoveLeft
			eventChan <- mainEvent
		case "a": // left
			mainEvent.Act = DoMoveLeft
			eventChan <- mainEvent
		// Punch
		case "z": // Left
			mainEvent.Act = -DoPunch
			eventChan <- mainEvent
		case "x": // Right
			mainEvent.Act = DoPunch
			eventChan <- mainEvent
		// Boxer 2
		// Move
		case "\033[A": // up
			oppEvent.Act = -DoMoveDown
			eventChan <- oppEvent
		case "\033[B": // down
			oppEvent.Act = DoMoveDown
			eventChan <- oppEvent
		case "\033[C": // right
			oppEvent.Act = -DoMoveLeft
			eventChan <- oppEvent
		case "\033[D": // left
			oppEvent.Act = DoMoveLeft
			eventChan <- oppEvent
		// Punch
		case "n": // Left
			oppEvent.Act = -DoPunch
			eventChan <- oppEvent
		case "m": // Right
			oppEvent.Act = DoPunch
			eventChan <- oppEvent
		case "q", "Q": // quit
			clearScreen()
			cancel()
			return
		}
	}
}

func (b *Boxer) Snapshot() BaseBoxer {
	return BaseBoxer{
		Area:      b.Area,
		Name:      b.Name,
		Color:     b.Color,
		Health:    b.Health,
		LastHit:   b.LastHit,
		Position:  b.Position,
		Direction: b.Direction,
		Situation: b.Situation,
	}
}

func (g GameInfo) HitLogic(opp IsOpponent, part BodyPart, dir HorDirection, PwPunch bool) {
	var (
		boxer = g.Boxers[opp]
	)
	// Hit effect
	// Power hit
	if PwPunch {
	}

	HitEffect(boxer, part, dir)
	if part != Arm {
		// Lower hp
		g.TakeDamage(opp, part)
	}
}
