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

func spawnBoxer(cage Cage, boxer Boxer, Hit string) {
	//Initialize chars of body
	var (
		head  = ":)"
		reset = "\033[0m"

		leftArm, rightArm           string
		leftShoulder, rightShoulder string
	)

	// Initialize body part poses
	leftArmPose := Position{X: boxer.Pose.X + 1, Y: boxer.Pose.Y - 2}
	rightArmPose := Position{X: boxer.Pose.X + 1, Y: boxer.Pose.Y + 2}

	leftShoulderPose := Position{X: leftArmPose.X, Y: boxer.Pose.Y - 1}
	rightShoulderPose := Position{X: rightArmPose.X, Y: boxer.Pose.Y + 1}

	if Hit != "" {
		switch Hit {
		case "UP":
			leftShoulder = "()═══O"
			leftShoulderPose = Position{X: leftShoulderPose.X + 1, Y: boxer.Pose.Y}
		case "DOWN":
			rightShoulder = "()═══O"
			rightShoulderPose = Position{X: rightShoulderPose.X + 1, Y: boxer.Pose.Y}
		}
	} else {
		leftShoulder = "()"
		rightShoulder = "()"
		leftArm = "╭╭══O"
		rightArm = "╰╰══O"
	}

	//spawn cage
	spawnCage(cage)

	// clear screen first
	//fmt.Print("\033[2J")

	// Shoulders
	fmt.Printf("\033[%d;%dH%s%s%s", leftShoulderPose.Y, leftShoulderPose.X, boxer.Color, leftShoulder, reset)
	fmt.Printf("\033[%d;%dH%s%s%s", rightShoulderPose.Y, rightShoulderPose.X, boxer.Color, rightShoulder, reset)

	// Arms
	fmt.Printf("\033[%d;%dH%s%s%s", leftArmPose.Y, leftArmPose.X, boxer.Color, leftArm, reset)
	fmt.Printf("\033[%d;%dH%s%s%s", rightArmPose.Y, rightArmPose.X, boxer.Color, rightArm, reset)

	// Head
	fmt.Printf("\033[%d;%dH%s%s%s", boxer.Pose.Y, boxer.Pose.X, boxer.Color, head, reset)

	// Move cursor below art
	fmt.Printf("\033[%d;1H", rightShoulderPose.Y+10)
}

func upArmHit(boxer Boxer) {

}

func main() {
	// Spawn cage
	cageLimit := Position{X: 120, Y: 50}
	cage := Cage{Limit: cageLimit, Color: "\033[44m"}

	// Spawn boxer
	headPose := Position{X: 10, Y: 10}
	boxer := Boxer{Pose: headPose, Color: "\033[31m"}
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

	spawnBoxer(cage, boxer, "")

	buf := make([]byte, 3) // arrow keys send 3 bytes like ESC [ A
	for {
		n, _ := os.Stdin.Read(buf)
		if n == 0 {
			continue
		}

		switch string(buf[:n]) {
		case "\033[A": // up
			if boxer.Pose.Y > 4 {
				boxer.Pose.Y--
				spawnBoxer(cage, boxer, "")
			}
		case "\033[B": // down
			if boxer.Pose.Y < cageLimit.Y-4 {
				boxer.Pose.Y++
				spawnBoxer(cage, boxer, "")
			}
		case "\033[C": // right
			if boxer.Pose.X < cageLimit.X-6 {
				boxer.Pose.X++
				spawnBoxer(cage, boxer, "")
			}
		case "\033[D": // left
			if boxer.Pose.X > 2 {
				boxer.Pose.X--
				spawnBoxer(cage, boxer, "")
			}
		//case "s":
		//	upArmHit(boxer)
		//case "x":
		//	downArmHit(boxer)
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
