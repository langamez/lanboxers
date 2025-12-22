package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

func clearScreen() {
	fmt.Print("\033[2J")
}

func spawnCage(cage Cage) {
	var (
		verticalWall   = "▌"
		horizontalWall = "─"

		reset = "\033[0m"
	)

	// Clear the screen
	clearScreen()

	// Spawn horizontal wall
	// todo get from env
	for y := cage.Area[Min].Y - WallLength; y <= WallLength; y++ {
		for x := cage.Area[Min].X - WallLength; x <= cage.Area[Max].X; x++ {
			fmt.Printf("\033[%d;%dH%s%s%s", y, x, cage.Color, horizontalWall, reset)
			fmt.Printf("\033[%d;%dH%s%s%s", cage.Area[Max].Y-(y-1), x, cage.Color, horizontalWall, reset)
		}
	}
	// Spawn vertical wall
	for y := cage.Area[Min].Y - WallLength; y <= cage.Area[Max].Y; y++ {
		for x := cage.Area[Min].X - WallLength; x <= WallLength; x++ {
			fmt.Printf("\033[%d;%dH%s%s%s", y, x, cage.Color, verticalWall, reset)
			fmt.Printf("\033[%d;%dH%s%s%s", y, cage.Area[Max].X-(x-1), cage.Color, verticalWall, reset)
		}
	}
	//Move cursor below term
	//fmt.Printf("\033[%d;1H", cage.Area[Max].Y+10)
}

func spawnGame(gameInfo GameInfo) {
	//Spawn cage
	spawnCage(gameInfo.Cage)
	//Spawn players health bar
	gameInfo.PrintHud()
	//Spawn player 1
	PrintBoxer(*gameInfo.Boxers[Main], AllBodyParts, gameInfo.Boxers[Main].Color)
	//Spawn player 2
	PrintBoxer(*gameInfo.Boxers[Opponent], AllBodyParts, gameInfo.Boxers[Opponent].Color)
	// Move cursor below term
	fmt.Printf("\033[%d;1H", gameInfo.Cage.Area[Max].Y+10)
}

func main() {
	var gameInfo GameInfo

	//todo get from env
	cageLimit := Position{X: 120, Y: 50}

	// Init cage, boxers (gameInfo)
	gameInfo = GameInfo{
		Cage: Cage{
			Color: "\033[44m",
			Area: map[Bounds]*Position{
				Min: {X: WallLength, Y: WallLength},
				Max: {X: cageLimit.X - WallLength, Y: cageLimit.Y - WallLength}}},
		Boxers: map[IsOpponent]*Boxer{
			Main: {
				Situation: Idle,
				Direction: Left,
				Name:      "User1",
				Color:     "\033[31m",
				Health:    MaxHealthPoint,
				Position:  Position{X: 10, Y: 10},
				Area:      map[Bounds]*Position{Min: {}, Max: {}}},
			Opponent: {
				Situation: Idle,
				Direction: Right,
				Name:      "User2",
				Color:     "\033[32m",
				Health:    MaxHealthPoint,
				Position:  Position{X: 30, Y: 10},
				Area:      map[Bounds]*Position{Min: {}, Max: {}}}}}

	// put terminal in raw mode (no buffering, instant keypresses)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		clearScreen()
		// todo reset buffer for holding a button
		term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(0)
	}()

	//Initialize occupied area
	gameInfo.CalArea()
	//Init spawn
	spawnGame(gameInfo)

	buf := make([]byte, 3) // arrow keys send 3 bytes like ESC [ A
	for {
		n, _ := os.Stdin.Read(buf)
		if n == 0 {
			continue
		}

		switch string(buf[:n]) {
		// Boxer 1
		// Move
		case "w": // up
			gameInfo.BoxerMove(Up, Main)
		case "s": // down
			gameInfo.BoxerMove(Down, Main)
		case "d": // right
			gameInfo.BoxerMove(Right, Main)
		case "a": // left
			gameInfo.BoxerMove(Left, Main)
		// Punch
		case "z": // Upper
			gameInfo.BoxerPunch(Up, Main)
		case "x": // Lower
			gameInfo.BoxerPunch(Down, Main)
		// Boxer 2
		// Move
		case "\033[A": // up
			gameInfo.BoxerMove(Up, Opponent)
		case "\033[B": // down
			gameInfo.BoxerMove(Down, Opponent)
		case "\033[C": // right
			gameInfo.BoxerMove(Right, Opponent)
		case "\033[D": // left
			gameInfo.BoxerMove(Left, Opponent)
		// Punch
		case "n": // Upper
			gameInfo.BoxerPunch(Up, Opponent)
		case "m": // Lower
			gameInfo.BoxerPunch(Down, Opponent)
		case "q", "Q": // quit
			clearScreen()
			return
		}

		time.Sleep(30 * time.Millisecond)
	}

}

//fmt.Println(boxer.Color + " ╭╭══O" + reset)
//fmt.Println(boxer.Color + " ()  " + reset)
//fmt.Println(boxer.Color + ":)   " + reset)
//fmt.Println(boxer.Color + " ()   " + reset)
//fmt.Println(boxer.Color + " ╰╰══O" + reset)
//}
