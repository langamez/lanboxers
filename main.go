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
		verticalWall   = "▌▌▌▌"
		horizontalWall = "─"

		reset = "\033[0m"
	)

	// Clear the screen
	clearScreen()

	// Spawn horizontal wall
	// todo get from env
	for y := 1; y < 4; y++ {
		for x := 1; x <= cage.Limit.X; x++ {
			fmt.Printf("\033[%d;%dH%s%s%s", y, x, cage.Color, horizontalWall, reset)
			fmt.Printf("\033[%d;%dH%s%s%s", cage.Limit.Y-(y-1), x, cage.Color, horizontalWall, reset)
		}
	}

	for y := 1; y <= cage.Limit.Y; y++ {
		fmt.Printf("\033[%d;%dH%s%s%s", y, 1, cage.Color, verticalWall, reset)
		fmt.Printf("\033[%d;%dH%s%s%s", y, cage.Limit.X, cage.Color, verticalWall, reset)
	}
	fmt.Printf("\033[%d;1H", cage.Limit.Y+10)
}

func spawnBoxers(cage Cage, boxerOne Boxer, boxerTwo Boxer) {
	//Spawn cage
	spawnCage(cage)
	//Spawn player 1
	PrintPlayer(boxerOne)
	//Spawn player 2
	PrintPlayer(boxerTwo)

	// Move cursor below art
	fmt.Printf("\033[%d;1H", cage.Limit.Y+10)
}

func main() {
	// Spawn cage
	cageLimit := Position{X: 120, Y: 50}
	cage := Cage{Limit: cageLimit, Color: "\033[44m"}

	// Spawn boxer
	boxerOne := Boxer{Pose: Position{X: 10, Y: 10}, Color: "\033[31m", Situation: Idle, Direction: Left}
	boxerTwo := Boxer{Pose: Position{X: 30, Y: 10}, Color: "\033[32m", Situation: Idle, Direction: Right}
	//spawnBoxer(cage, boxer)

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
		term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(0)
	}()

	spawnBoxers(cage, boxerOne, boxerTwo)

	buf := make([]byte, 3) // arrow keys send 3 bytes like ESC [ A
	for {
		n, _ := os.Stdin.Read(buf)
		if n == 0 {
			continue
		}

		switch string(buf[:n]) {
		// Boxer 1
		case "w": // up
			if boxerOne.Pose.Y > 6 {
				boxerOne.Pose.Y--
				spawnBoxers(cage, boxerOne, boxerTwo)
			}
		case "s": // down
			if boxerOne.Pose.Y < cageLimit.Y-6 {
				boxerOne.Pose.Y++
				spawnBoxers(cage, boxerOne, boxerTwo)
			}
		case "d": // right
			if boxerOne.Pose.X < cageLimit.X-6 {
				boxerOne.Pose.X++
				spawnBoxers(cage, boxerOne, boxerTwo)
			}
		case "a": // left
			if boxerOne.Pose.X > 5 {
				boxerOne.Pose.X--
				spawnBoxers(cage, boxerOne, boxerTwo)
			}
		case "z":
			// Form init hit situations
			boxerOne.Situation = UpPunchInit
			boxerTwo.Situation = hitCheck(boxerOne, boxerTwo, Up, true)
			spawnBoxers(cage, boxerOne, boxerTwo)

			//Sleep to avoid blank term
			time.Sleep(100 * time.Millisecond)
			//time.Sleep(time.Second)

			// Form hit situations
			boxerOne.Situation = UpPunch
			boxerTwo.Situation = hitCheck(boxerOne, boxerTwo, Up, false)
			spawnBoxers(cage, boxerOne, boxerTwo)

			//Sleep to avoid blank term
			time.Sleep(100 * time.Millisecond)
			//time.Sleep(time.Second)

			//Reform back to default
			boxerOne.Situation = Idle
			boxerTwo.Situation = Idle
			spawnBoxers(cage, boxerOne, boxerTwo)
		case "x":
			// Form init hit situations
			boxerOne.Situation = DownPunchInit
			boxerTwo.Situation = hitCheck(boxerOne, boxerTwo, Down, true)
			spawnBoxers(cage, boxerOne, boxerTwo)

			//Sleep to avoid blank term
			time.Sleep(100 * time.Millisecond)
			//time.Sleep(time.Second)

			// Form hit situations
			boxerOne.Situation = DownPunch
			boxerTwo.Situation = hitCheck(boxerOne, boxerTwo, Down, false)
			spawnBoxers(cage, boxerOne, boxerTwo)

			//Sleep to avoid blank term
			time.Sleep(100 * time.Millisecond)
			//time.Sleep(time.Second)

			//Reform back to default
			boxerOne.Situation = Idle
			boxerTwo.Situation = Idle
			spawnBoxers(cage, boxerOne, boxerTwo)
		// Boxer 2
		case "\033[A": // up
			if boxerTwo.Pose.Y > 6 {
				boxerTwo.Pose.Y--
				spawnBoxers(cage, boxerOne, boxerTwo)
			}
		case "\033[B": // down
			if boxerTwo.Pose.Y < cageLimit.Y-6 {
				boxerTwo.Pose.Y++
				spawnBoxers(cage, boxerOne, boxerTwo)
			}
		case "\033[C": // right
			if boxerTwo.Pose.X < cageLimit.X-6 {
				boxerTwo.Pose.X++
				spawnBoxers(cage, boxerOne, boxerTwo)
			}
		case "\033[D": // left
			if boxerTwo.Pose.X > 5 {
				boxerTwo.Pose.X--
				spawnBoxers(cage, boxerOne, boxerTwo)
			}
		case "n":
			// Form init hit situations
			boxerTwo.Situation = UpPunchInit
			boxerOne.Situation = hitCheck(boxerTwo, boxerOne, Up, true)
			spawnBoxers(cage, boxerOne, boxerTwo)

			//Sleep to avoid blank term
			//time.Sleep(time.Second)
			time.Sleep(100 * time.Millisecond)

			// Form hit situations
			boxerTwo.Situation = UpPunch
			boxerOne.Situation = hitCheck(boxerTwo, boxerOne, Up, false)
			spawnBoxers(cage, boxerOne, boxerTwo)

			//Sleep to avoid blank term
			//time.Sleep(time.Second)
			time.Sleep(100 * time.Millisecond)

			//Reform back to default
			boxerTwo.Situation = Idle
			boxerOne.Situation = Idle
			spawnBoxers(cage, boxerOne, boxerTwo)
		case "m":
			// Form init hit situations
			boxerTwo.Situation = DownPunchInit
			boxerOne.Situation = hitCheck(boxerTwo, boxerOne, Down, true)
			spawnBoxers(cage, boxerOne, boxerTwo)

			//Sleep to avoid blank term
			//time.Sleep(time.Second)
			time.Sleep(100 * time.Millisecond)

			// Form hit situations
			boxerTwo.Situation = DownPunch
			boxerOne.Situation = hitCheck(boxerTwo, boxerOne, Down, false)
			spawnBoxers(cage, boxerOne, boxerTwo)

			//Sleep to avoid blank term
			//time.Sleep(time.Second)
			time.Sleep(100 * time.Millisecond)

			//Reform back to default
			boxerTwo.Situation = Idle
			boxerOne.Situation = Idle
			spawnBoxers(cage, boxerOne, boxerTwo)
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
