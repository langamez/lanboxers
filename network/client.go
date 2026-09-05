package network

import (
	"fmt"
	"github.com/langamez/lanboxers/render"
	"net"
)

func Dial(address string, port int) (*Connection, error) {
	netConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}

	conn := newConnection(netConn)

	if err := conn.send("CONNECTED"); err != nil {
		return nil, err
	}

	msg, err := conn.receive()
	if err != nil {
		return nil, err
	}

	if msg != "OK" {
		conn.Close()
		return nil, fmt.Errorf("invalid handshake")
	}

	// todo: move these to render
	fmt.Println("Connected to host!")
	fmt.Println("Starting game!")
	render.Frame(20)

	return conn, nil
}
