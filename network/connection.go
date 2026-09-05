package network

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"

	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/render"
)

type Connection struct {
	conn net.Conn
	mu   sync.Mutex
}

func newConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

// send and receive
func (c *Connection) send(data string) error {
	var byteData = []byte(data)
	c.mu.Lock()
	defer c.mu.Unlock()

	length := uint32(len(byteData))

	if err := binary.Write(c.conn, binary.BigEndian, length); err != nil {
		return err
	}

	_, err := c.conn.Write(byteData)
	return err
}

func (c *Connection) SendAct(act domain.Action) error {
	var actStr = strconv.Itoa(int(act))
	if err := c.send(actStr); err != nil {
		return err
	}
	return nil
}

func (c *Connection) receive() (string, error) {
	var length uint32

	if err := binary.Read(c.conn, binary.BigEndian, &length); err != nil {
		return "", err
	}

	buffer := make([]byte, length)

	_, err := io.ReadFull(c.conn, buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func (c *Connection) Listen(eventChan chan<- domain.Action) {
	for {
		msg, err := c.receive()
		if err != nil {
			// todo: move these to render
			fmt.Printf("error %s", err)
			render.Frame(100)
		}

		// Convert to Act
		msgInt, err := strconv.Atoi(msg)
		if err != nil {
		}

		eventChan <- domain.Action(msgInt)
	}
}
