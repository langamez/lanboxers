package render

import (
	"fmt"
	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/sprites"
)

func drawCage(cage domain.Cage) {
	var (
		verticalWall   = "▌"
		horizontalWall = "─"

		reset = "\033[0m"
	)

	// Clear the screen
	ClearScreen()

	// Spawn horizontal wall
	// todo get from env
	for y := cage.Area.Min.Y - domain.WallLength; y <= domain.WallLength; y++ {
		for x := cage.Area.Min.X - domain.WallLength; x <= cage.Area.Max.X; x++ {
			fmt.Printf("\033[%d;%dH%s%s%s", y, x, cage.Color, horizontalWall, reset)
			fmt.Printf("\033[%d;%dH%s%s%s", cage.Area.Max.Y-(y-1), x, cage.Color, horizontalWall, reset)
		}
	}
	// Spawn vertical wall
	for y := cage.Area.Min.Y - domain.WallLength; y <= cage.Area.Max.Y; y++ {
		for x := cage.Area.Min.X - domain.WallLength; x <= domain.WallLength; x++ {
			fmt.Printf("\033[%d;%dH%s%s%s", y, x, cage.Color, verticalWall, reset)
			fmt.Printf("\033[%d;%dH%s%s%s", y, cage.Area.Max.X-(x-1), cage.Color, verticalWall, reset)
		}
	}
	//Move cursor below term
	//fmt.Printf("\033[%d;1H", cage.Area[Max].Y+10)
}

func DrawGame(renderInfo domain.RenderInfo) {
	// Render cage
	drawCage(renderInfo.Cage)
	// Render players health bar
	drawHud(renderInfo)
	//Spawn player 1
	DrawBoxer(Snapshot(renderInfo.Boxers[domain.PlayerMain].BaseBoxer), sprites.AllBodyParts, renderInfo.Boxers[domain.PlayerMain].Color)
	//Spawn player 2
	DrawBoxer(Snapshot(renderInfo.Boxers[domain.PlayerOpponent].BaseBoxer), sprites.AllBodyParts, renderInfo.Boxers[domain.PlayerOpponent].Color)
	// Move cursor below term
	fmt.Printf("\033[%d;1H", renderInfo.Cage.Area.Max.Y+10)
}
