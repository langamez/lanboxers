package network

import (
	"fmt"
	"net"

	"github.com/langamez/lanboxers/render"
)

func StartServer(port int) (*Connection, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	fmt.Println("Listening...")

	netConn, err := ln.Accept() // Blocks until one client connects.
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected")

	conn := newConnection(netConn)

	msg, err := conn.receive()
	if err != nil {
		return nil, err
	}

	if msg != "CONNECTED" {
		conn.Close()
		return nil, fmt.Errorf("invalid handshake")
	}

	if err := conn.send("OK"); err != nil {
		return nil, err
	}

	// todo: move these to render
	fmt.Println("Client connected!")
	fmt.Println("Starting game!")
	render.Frame(20)

	// Start normal message processing
	return conn, nil
}
