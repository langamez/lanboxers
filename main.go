package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	PrintBoxer(gameInfo.Boxers[Main].Snapshot(), AllBodyParts, gameInfo.Boxers[Main].Color)
	//Spawn player 2
	PrintBoxer(gameInfo.Boxers[Opponent].Snapshot(), AllBodyParts, gameInfo.Boxers[Opponent].Color)
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
		Boxers: Boxers{
			Main: {
				Lock: sync.Mutex{},
				BaseBoxer: &BaseBoxer{
					Opponent:  Main,
					Situation: Idle,
					Direction: Left,
					Name:      "User1",
					Color:     "\033[31m",
					Health:    MaxHealthPoint,
					Position:  Position{X: 10, Y: 10},
					Area:      map[Bounds]*Position{Min: {}, Max: {}}}},
			Opponent: {
				Lock: sync.Mutex{},
				BaseBoxer: &BaseBoxer{
					Opponent:  Opponent,
					Situation: Idle,
					Direction: Right,
					Name:      "User2",
					Color:     "\033[32m",
					Health:    MaxHealthPoint,
					Position:  Position{X: 30, Y: 10},
					Area:      map[Bounds]*Position{Min: {}, Max: {}}}}}}

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

	input := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())
	go ReadKeyboard(cancel, input)
	go EventDispatcher(input, Router)

	// Main boxer
	go func() {
		for event := range MainChan {
			switch event.Act {
			case DoPunch:
				gameInfo.BoxerPunch(Down, Main)
			case -DoPunch:
				gameInfo.BoxerPunch(Up, Main)

			case DoMoveDown:
				gameInfo.BoxerMove(Down, Main)
			case -DoMoveDown:
				gameInfo.BoxerMove(Up, Main)

			case DoMoveLeft:
				gameInfo.BoxerMove(Left, Main)
			case -DoMoveLeft:
				gameInfo.BoxerMove(Right, Main)
			default:
				panic("unhandled default case")
			}
		}
	}()

	// Opponent boxer
	go func() {
		for event := range OppChan {
			switch event.Act {
			case DoPunch:
				gameInfo.BoxerPunch(Down, Opponent)
			case -DoPunch:
				gameInfo.BoxerPunch(Up, Opponent)

			case DoMoveDown:
				gameInfo.BoxerMove(Down, Opponent)
			case -DoMoveDown:
				gameInfo.BoxerMove(Up, Opponent)

			case DoMoveLeft:
				gameInfo.BoxerMove(Left, Opponent)
			case -DoMoveLeft:
				gameInfo.BoxerMove(Right, Opponent)
			default:
				panic("unhandled default case")
			}
		}
	}()

	select {
	case <-ctx.Done():
		return
	}
}

//fmt.Println(boxer.Color + " ╭╭══O" + reset)
//fmt.Println(boxer.Color + " ()  " + reset)
//fmt.Println(boxer.Color + ":)   " + reset)
//fmt.Println(boxer.Color + " ()   " + reset)
//fmt.Println(boxer.Color + " ╰╰══O" + reset)
//}
