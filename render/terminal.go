package render

import (
	"fmt"
	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/sprites"
	"strings"
	"time"
	"unicode/utf8"
)

func PrintOn(position domain.Position, color, char string) {
	var (
		reset = "\033[0m"
	)
	fmt.Printf("\033[%d;%dH%s%s%s", position.Y, position.X, color, char, reset)
}

func Frame(n float64) {
	n += 3
	n = n * 10
	// Move cursor to end
	fmt.Printf("\033[%d;1H", 10)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func ClearScreen() {
	fmt.Print("\033[2J")
}

func ClearBoxer(boxer domain.BaseBoxer, parts map[domain.BodyPart][]domain.Direction) {
	var (
		char     string
		charLen  int
		position domain.Position
	)

	for part := range parts {
		for _, dir := range parts[part] {
			position = sprites.GetPosition(part, dir, boxer.Situation)
			charLen = utf8.RuneCountInString(sprites.GetBodyChar(part, dir, boxer.Situation, boxer.Direction))
			char = strings.Repeat(" ", charLen)
			// Check if it's right boxer change the direction to opposite
			if part != domain.Head &&
				boxer.Direction == domain.Right {
				position = sprites.GetPosition(part, dir.Opposite(), boxer.Situation)
			}
			position = sprites.CalculatePartPosition(charLen, boxer, position)
			PrintOn(position, "", char)
		}
	}
}

func DrawBoxer(boxer domain.BaseBoxer, parts map[domain.BodyPart][]domain.Direction, color string) {
	var (
		char     string
		position domain.Position
	)
	for part := range parts {
		for _, dir := range parts[part] {
			position = sprites.GetPosition(part, dir, boxer.Situation)
			char = sprites.GetBodyChar(part, dir, boxer.Situation, boxer.Direction)
			// Check if it's right boxer change the direction to opposite
			if part != domain.Head &&
				boxer.Direction == domain.Right {
				position = sprites.GetPosition(part, dir.Opposite(), boxer.Situation)
			}
			position = sprites.CalculatePartPosition(utf8.RuneCountInString(char), boxer, position)
			PrintOn(position, color, char)
		}
	}
}
