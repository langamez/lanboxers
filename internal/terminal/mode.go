package terminal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var originalState *term.State

func RawMode() error {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	originalState = state

	removeCursor()
	return nil
}

func NormalMode() error {
	if originalState != nil {
		restoreCursor()
		return term.Restore(int(os.Stdin.Fd()), originalState)
	}
	return nil
}

func removeCursor() {
	fmt.Print("\033[?25l")
}

func restoreCursor() {
	fmt.Print("\033[?25h")
}
