package terminal

import (
	"os"

	"golang.org/x/term"
)

const (
	defaultWidth  int = 80
	defaultHeight int = 30
)

func Size() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return defaultWidth, defaultHeight
	}
	return width, height
}
