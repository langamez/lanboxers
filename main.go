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
	for x := 1; x <= cage.Limit.X; x++ {
		fmt.Printf("\033[%d;%dH%s%s%s", 1, x, cage.Color, horizontalWall, reset)
		fmt.Printf("\033[%d;%dH%s%s%s", cage.Limit.Y, x, cage.Color, horizontalWall, reset)
	}

	for y := 1; y <= cage.Limit.Y; y++ {
		fmt.Printf("\033[%d;%dH%s%s%s", y, 1, cage.Color, verticalWall, reset)
		fmt.Printf("\033[%d;%dH%s%s%s", y, cage.Limit.X, cage.Color, verticalWall, reset)
	}
	fmt.Printf("\033[%d;1H", cage.Limit.Y+10)
}

func spawnBoxers(cage Cage, boxerOne Boxer, boxerTwo Boxer, Hit string) {
	//Initialize chars of body
	var (
		boxerOneHead = ":)"
		boxerTwoHead = "(:"
		reset        = "\033[0m"

		// player 1
		boxerOneLeftArm, boxerOneRightArm           string
		boxerOneLeftShoulder, boxerOneRightShoulder string
		// player 2
		boxerTwoLeftArm, boxerTwoRightArm           string
		boxerTwoLeftShoulder, boxerTwoRightShoulder string
	)

	// Initialize body part poses
	// player 1 body parts position
	boxerOneLeftArmPose := Position{X: boxerOne.Pose.X + 1, Y: boxerOne.Pose.Y - 2}
	boxerOneRightArmPose := Position{X: boxerOne.Pose.X + 1, Y: boxerOne.Pose.Y + 2}

	boxerOneLeftShoulderPose := Position{X: boxerOneLeftArmPose.X, Y: boxerOne.Pose.Y - 1}
	boxerOneRightShoulderPose := Position{X: boxerOneRightArmPose.X, Y: boxerOne.Pose.Y + 1}

	// player 2 body parts position
	boxerTwoLeftArmPose := Position{X: boxerTwo.Pose.X - 4, Y: boxerTwo.Pose.Y + 2}
	boxerTwoRightArmPose := Position{X: boxerTwo.Pose.X - 4, Y: boxerTwo.Pose.Y - 2}

	boxerTwoLeftShoulderPose := Position{X: boxerTwoLeftArmPose.X + 3, Y: boxerTwo.Pose.Y + 1}
	boxerTwoRightShoulderPose := Position{X: boxerTwoRightArmPose.X + 3, Y: boxerTwo.Pose.Y - 1}

	if Hit != "" {
		switch Hit {
		case "upInit":
			// player 1
			boxerOneLeftArm = ""
			boxerOneRightArm = "╰╰══O"
			boxerOneRightShoulder = "()"
			boxerOneLeftShoulder = "()══O"
			boxerOneLeftShoulderPose = Position{X: boxerOneLeftShoulderPose.X + 1, Y: boxerOneLeftShoulderPose.Y}

			// player 2
			boxerTwoLeftArm = "O══╯╯"
			boxerTwoRightArm = "O══╮╮"
			boxerTwoLeftShoulder = "()"
			boxerTwoRightShoulder = "()"
		case "up":
			// player 1
			boxerOneLeftArm = ""
			boxerOneRightArm = "╰╰══O"
			boxerOneRightShoulder = "()"
			boxerOneLeftShoulder = "()═══O"
			boxerOneLeftShoulderPose = Position{X: boxerOneLeftShoulderPose.X + 1, Y: boxerOneLeftShoulderPose.Y}

			// player 2
			boxerTwoHead = "),: **"
			boxerTwoLeftArm = "O══╯╯"
			boxerTwoRightArm = "O══╮╮"
			boxerTwoLeftShoulder = "()"
			boxerTwoRightShoulder = "()"
		case "upHit":
			// player 1
			boxerOneLeftArm = ""
			boxerOneRightArm = "╰╰══O"
			boxerOneRightShoulder = "()"
			boxerOneLeftShoulder = "()═══O"
			boxerOneLeftShoulderPose = Position{X: boxerOneLeftShoulderPose.X + 1, Y: boxerOneLeftShoulderPose.Y}

			// player 2
			boxerTwoHead = "),: **"
			boxerTwoLeftArm = "O══╯╯"
			boxerTwoRightArm = "O══╮╮"
			boxerTwoLeftShoulder = "()"
			boxerTwoRightShoulder = "()"
		case "downInit":
			// player 1
			boxerOneRightArm = ""
			boxerOneLeftArm = "╭╭══O"
			boxerOneLeftShoulder = "()"
			boxerOneRightShoulder = "()══O"
			boxerOneRightShoulderPose = Position{X: boxerOneRightShoulderPose.X + 1, Y: boxerOneRightShoulderPose.Y}

			// player 2
			boxerTwoLeftArm = "O══╯╯"
			boxerTwoRightArm = "O══╮╮"
			boxerTwoLeftShoulder = "()"
			boxerTwoRightShoulder = "()"
		case "down":
			// player 1
			boxerOneRightArm = ""
			boxerOneLeftArm = "╭╭══O"
			boxerOneLeftShoulder = "()"
			boxerOneRightShoulder = "()═══O"
			boxerOneRightShoulderPose = Position{X: boxerOneRightShoulderPose.X + 1, Y: boxerOneRightShoulderPose.Y}

			// player 2
			boxerTwoHead = "),: **"
			boxerTwoLeftArm = "O══╯╯"
			boxerTwoRightArm = "O══╮╮"
			boxerTwoLeftShoulder = "()"
			boxerTwoRightShoulder = "()"
		case "downHit":
			switch "actor" {
			case "one":
				// player 1
				boxerOneRightArm = ""
				boxerOneLeftArm = "╭╭══O"
				boxerOneLeftShoulder = "()"
				boxerOneRightShoulder = "()═══O"
				boxerOneRightShoulderPose = Position{X: boxerOneRightShoulderPose.X + 1, Y: boxerOneRightShoulderPose.Y}

				// player 2
				boxerTwoHead = "),: **"
				boxerTwoLeftArm = "O══╯╯"
				boxerTwoRightArm = "O══╮╮"
				boxerTwoLeftShoulder = "()"
				boxerTwoRightShoulder = "()"
			case "two":

			}
		}
	} else {
		// player 1
		boxerOneLeftArm = "╭╭══O"
		boxerOneRightArm = "╰╰══O"
		boxerOneLeftShoulder = "()"
		boxerOneRightShoulder = "()"

		// player 2
		boxerTwoLeftArm = "O══╯╯"
		boxerTwoRightArm = "O══╮╮"
		boxerTwoLeftShoulder = "()"
		boxerTwoRightShoulder = "()"
	}

	//spawn cage
	spawnCage(cage)

	// player 1
	// Shoulders
	if boxerOneLeftShoulder != "" {
		fmt.Printf(
			"\033[%d;%dH%s%s%s", boxerOneLeftShoulderPose.Y, boxerOneLeftShoulderPose.X, boxerOne.Color, boxerOneLeftShoulder, reset)
	}
	if boxerOneRightShoulder != "" {
		fmt.Printf(
			"\033[%d;%dH%s%s%s", boxerOneRightShoulderPose.Y, boxerOneRightShoulderPose.X, boxerOne.Color, boxerOneRightShoulder, reset)
	}

	// Arms
	if boxerOneLeftArm != "" {
		fmt.Printf("\033[%d;%dH%s%s%s", boxerOneLeftArmPose.Y, boxerOneLeftArmPose.X, boxerOne.Color, boxerOneLeftArm, reset)
	}
	if boxerOneRightArm != "" {
		fmt.Printf("\033[%d;%dH%s%s%s", boxerOneRightArmPose.Y, boxerOneRightArmPose.X, boxerOne.Color, boxerOneRightArm, reset)
	}

	// Head
	fmt.Printf("\033[%d;%dH%s%s%s", boxerOne.Pose.Y, boxerOne.Pose.X, boxerOne.Color, boxerOneHead, reset)

	// player 2
	// Shoulders
	if boxerTwoLeftShoulder != "" {
		fmt.Printf(
			"\033[%d;%dH%s%s%s", boxerTwoLeftShoulderPose.Y, boxerTwoLeftShoulderPose.X, boxerTwo.Color, boxerTwoLeftShoulder, reset)
	}
	if boxerTwoRightShoulder != "" {
		fmt.Printf(
			"\033[%d;%dH%s%s%s", boxerTwoRightShoulderPose.Y, boxerTwoRightShoulderPose.X, boxerTwo.Color, boxerTwoRightShoulder, reset)
	}

	// Arms
	if boxerTwoLeftArm != "" {
		fmt.Printf("\033[%d;%dH%s%s%s", boxerTwoLeftArmPose.Y, boxerTwoLeftArmPose.X, boxerTwo.Color, boxerTwoLeftArm, reset)
	}
	if boxerTwoRightArm != "" {
		fmt.Printf("\033[%d;%dH%s%s%s", boxerTwoRightArmPose.Y, boxerTwoRightArmPose.X, boxerTwo.Color, boxerTwoRightArm, reset)
	}

	// Head
	fmt.Printf("\033[%d;%dH%s%s%s", boxerTwo.Pose.Y, boxerTwo.Pose.X, boxerTwo.Color, boxerTwoHead, reset)

	// Move cursor below art
	fmt.Printf("\033[%d;1H", cage.Limit.Y+10)
}

func main() {
	// Spawn cage
	cageLimit := Position{X: 120, Y: 50}
	cage := Cage{Limit: cageLimit, Color: "\033[44m"}

	// Spawn boxer
	boxerOne := Boxer{Pose: Position{X: 10, Y: 10}, Color: "\033[31m"}
	boxerTwo := Boxer{Pose: Position{X: 30, Y: 10}, Color: "\033[32m"}
	//spawnBoxer(cage, boxer, "")

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

	spawnBoxers(cage, boxerOne, boxerTwo, "")

	buf := make([]byte, 3) // arrow keys send 3 bytes like ESC [ A
	for {
		n, _ := os.Stdin.Read(buf)
		if n == 0 {
			continue
		}

		switch string(buf[:n]) {
		case "\033[A": // up
			if boxerOne.Pose.Y > 4 {
				boxerOne.Pose.Y--
				spawnBoxers(cage, boxerOne, boxerTwo, "")
			}
		case "\033[B": // down
			if boxerOne.Pose.Y < cageLimit.Y-4 {
				boxerOne.Pose.Y++
				spawnBoxers(cage, boxerOne, boxerTwo, "")
			}
		case "\033[C": // right
			if boxerOne.Pose.X < cageLimit.X-6 {
				boxerOne.Pose.X++
				spawnBoxers(cage, boxerOne, boxerTwo, "")
			}
		case "\033[D": // left
			if boxerOne.Pose.X > 2 {
				boxerOne.Pose.X--
				spawnBoxers(cage, boxerOne, boxerTwo, "")
			}
		case "s":
			spawnBoxers(cage, boxerOne, boxerTwo, "upInit")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "up")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "upHit")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "")
		case "x":
			spawnBoxers(cage, boxerOne, boxerTwo, "downInit")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "down")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "")
		case "k":
			spawnBoxers(cage, boxerOne, boxerTwo, "upInit")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "up")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "upHit")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "")
		case "m":
			spawnBoxers(cage, boxerOne, boxerTwo, "downInit")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "down")
			time.Sleep(100 * time.Millisecond)
			spawnBoxers(cage, boxerOne, boxerTwo, "")
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
