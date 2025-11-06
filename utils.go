package main

import (
	"fmt"
)

func PrintOn(boxer Boxer, position Position, char string) {
	var (
		reset = "\033[0m"
	)
	fmt.Printf("\033[%d;%dH%s%s%s", position.Y, position.X, boxer.Color, char, reset)
}

func PrintPlayer(boxer Boxer) {
	var (
		reset = "\033[0m"
	)
	//Head
	fmt.Printf("\033[%d;%dH%s%s%s",
		boxer.Pose.Y+Situations[boxer.Situation].headPose[boxer.Direction].Y,
		boxer.Pose.X+Situations[boxer.Situation].headPose[boxer.Direction].X,
		boxer.Color,
		Situations[boxer.Situation].head[boxer.Direction],
		reset)
	//Shoulder
	//Upper
	fmt.Printf("\033[%d;%dH%s%s%s",
		boxer.Pose.Y+Situations[boxer.Situation].shoulderPose[boxer.Direction][Up].Y,
		boxer.Pose.X+Situations[boxer.Situation].shoulderPose[boxer.Direction][Up].X,
		boxer.Color,
		Situations[boxer.Situation].shoulder[boxer.Direction][Up],
		reset)
	//Down
	fmt.Printf("\033[%d;%dH%s%s%s",
		boxer.Pose.Y+Situations[boxer.Situation].shoulderPose[boxer.Direction][Down].Y,
		boxer.Pose.X+Situations[boxer.Situation].shoulderPose[boxer.Direction][Down].X,
		boxer.Color,
		Situations[boxer.Situation].shoulder[boxer.Direction][Down],
		reset)
	//Arm
	//Upper
	fmt.Printf("\033[%d;%dH%s%s%s",
		boxer.Pose.Y+Situations[boxer.Situation].armPose[boxer.Direction][Up].Y,
		boxer.Pose.X+Situations[boxer.Situation].armPose[boxer.Direction][Up].X,
		boxer.Color,
		Situations[boxer.Situation].arm[boxer.Direction][Up],
		reset)
	//Down
	fmt.Printf("\033[%d;%dH%s%s%s",
		boxer.Pose.Y+Situations[boxer.Situation].armPose[boxer.Direction][Down].Y,
		boxer.Pose.X+Situations[boxer.Situation].armPose[boxer.Direction][Down].X,
		boxer.Color,
		Situations[boxer.Situation].arm[boxer.Direction][Down],
		reset)
}

func hitCheck(boxer Boxer, subBoxer Boxer, direction Direction, init bool) SituationType {
	var hitDiffPose Position
	switch boxer.Direction {
	case Right:
		switch direction {
		case Up:
			hitDiffPose = Position{X: -7, Y: -1}
		case Down:
			hitDiffPose = Position{X: -7, Y: 1}
		}
	case Left:
		switch direction {
		case Up:
			// Find hit position
			hitDiffPose = Position{X: 7, Y: -1}
		case Down:
			// Find hit position
			hitDiffPose = Position{X: 7, Y: 1}
		}
	}
	// Find hit position
	hitPose := Position{
		X: boxer.Pose.X + hitDiffPose.X,
		Y: boxer.Pose.Y + hitDiffPose.Y,
	}
	// Check for hit
	switch hitPose {
	case subBoxer.Pose:
		if init {
			// Head got hit init
			return HeadHitInit
		}
		// Head got hit
		return HeadHit
	case Position{X: subBoxer.Pose.X, Y: subBoxer.Pose.Y - 1}:
		if init {
			// Upper shoulder got hit init
			return ShoulderHitInit
		}
		// Upper shoulder got hit
		return ShoulderHit
	case Position{X: subBoxer.Pose.X, Y: subBoxer.Pose.Y + 1}:
		if init {
			// Lower shoulder got hit init
			return ShoulderHitInit
		}
		// Lower shoulder got hit
		return ShoulderHit
	default:
		return Idle
	}
}
