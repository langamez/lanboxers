package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
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
	var mainBoxer IsOpponent

	if mainBoxer {
		//todo get from env
		cageLimit := Position{X: 80, Y: 30}

		// Init cage
		gameInfo = GameInfo{
			Cage: Cage{
				Color: "\033[44m",
				Area: map[Bounds]*Position{
					Min: {X: WallLength, Y: WallLength},
					Max: {X: cageLimit.X - WallLength, Y: cageLimit.Y - WallLength}}}}

		//PrintOn(Position{20, 25}, "\u001B[44m", strconv.Itoa(gameInfo.Cage.Area[Max].X/3))
		//PrintOn(Position{20, 26}, "\u001B[44m", strconv.Itoa((gameInfo.Cage.Area[Max].X/3)*2))
		//frame(1000)

		// Init boxers
		gameInfo.Boxers = Boxers{
			Main: {
				Lock: sync.Mutex{},
				BaseBoxer: &BaseBoxer{
					Opponent:  Main,
					Situation: Idle,
					Direction: Left,
					Name:      "User1",
					Color:     "\033[31m",
					Health:    MaxHealthPoint,
					Position:  Position{X: gameInfo.Cage.Area[Max].X / 3, Y: gameInfo.Cage.Area[Max].Y / 2},
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
					Position:  Position{X: (gameInfo.Cage.Area[Max].X / 3) * 2, Y: gameInfo.Cage.Area[Max].Y / 2},
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

		// Initialize websocket (room) on server
	} else {
	}

	// Initialize connection to server
	url := "ws://localhost:8787/ws"

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer conn.Close()

	fmt.Println("connected")

	//Init spawn
	spawnGame(gameInfo)

	input := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())
	gameInfo.Cancel = cancel

	go ReadKeyboard(gameInfo.Cancel, input)
	go EventDispatcher(input, Router)

	// Receive the "connected" message
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			fmt.Println("received:", string(msg))

			msg_dict := map[string]string{
				"Main": "main",
				"Act":  "DoPunch",
			}
			_ = msg_dict

			gameInfo.BoxerPunch(Left, IsOpponent(bool(msg_dict["Main"])))

		}
	}()

	// Send messages
	//for {
	//
	//	err := conn.WriteMessage(
	//		websocket.TextMessage,
	//		[]byte(input),
	//	)
	//
	//	if err != nil {
	//		log.Println("write error:", err)
	//		return
	//	}
	//}

	// Main boxer
	go func() {
		for event := range MainChan {
			switch event.Act {
			case DoPunch:
				msg := map[string]string{
					"Main": "main",
					"Act":  "DoPunch",
				}
				_ = msg
				gameInfo.BoxerPunch(Left, mainBoxer)
			case -DoPunch:
				gameInfo.BoxerPunch(Right, mainBoxer)

			case DoMoveDown:
				gameInfo.BoxerMove(Lower, mainBoxer)
			case -DoMoveDown:
				gameInfo.BoxerMove(Upper, mainBoxer)

			case DoMoveLeft:
				gameInfo.BoxerMove(Left, mainBoxer)
			case -DoMoveLeft:
				gameInfo.BoxerMove(Right, mainBoxer)

			default:
				panic("unhandled default case")
			}
		}
	}()

	// Opponent boxer
	//go func() {
	//	for event := range OppChan {
	//		switch event.Act {
	//		case DoPunch:
	//			gameInfo.BoxerPunch(Left, Opponent)
	//		case -DoPunch:
	//			gameInfo.BoxerPunch(Right, Opponent)
	//
	//		case DoMoveDown:
	//			gameInfo.BoxerMove(Lower, Opponent)
	//		case -DoMoveDown:
	//			gameInfo.BoxerMove(Upper, Opponent)
	//
	//		case DoMoveLeft:
	//			gameInfo.BoxerMove(Left, Opponent)
	//		case -DoMoveLeft:
	//			gameInfo.BoxerMove(Right, Opponent)
	//		default:
	//			panic("unhandled default case")
	//		}
	//	}
	//}()

	select {
	case <-ctx.Done():
		return
	}
}
