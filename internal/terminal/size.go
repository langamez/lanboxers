package terminal

import (
	"github.com/langamez/lanboxers/config"
	"os"

	"golang.org/x/term"
)

const (
	defaultWidth  = 80
	defaultHeight = 30
)

func Size() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		return width, height
	}

	return config.Int("LANBOX_WIDTH", defaultWidth),
		config.Int("LANBOX_HEIGHT", defaultHeight)
}
