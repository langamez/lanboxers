package terminal

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

type Session struct {
	state *term.State
}

func Start(onExit func()) (*Session, error) {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}

	// todo: check why we needed this and now works fine
	//defer term.Restore(int(os.Stdin.Fd()), state)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		onExit()
	}()

	return &Session{
		state: state,
	}, nil
}

func (s *Session) Close() error {
	return term.Restore(int(os.Stdin.Fd()), s.state)
}
