package main

import (
	"fmt"
	"github.com/langamez/lanboxers/game"
	"github.com/langamez/lanboxers/internal/terminal"
	"github.com/langamez/lanboxers/render"
)

func main() {
	g := game.NewGame(terminal.Size())

	session, err := game.NewSession(g)
	if err != nil {
		// todo: do in render
		fmt.Println(err)
		return
	}

	if err := terminal.RawMode(); err != nil {
		fmt.Println(err)
		return
	}
	defer terminal.NormalMode()

	render.DrawGame(g.RenderConverter())

	session.Start()

	<-g.Context.Done()
}
