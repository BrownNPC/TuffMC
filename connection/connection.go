package connection

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
	"tuff/ds"
	"tuff/packet"

	"github.com/coder/websocket"
)

type Connection struct {
	// the Username of the connection joining.
	Username   string
	State      packet.ConnectionState
	isLoggedIn atomic.Bool
	// underlying socket
	conn     net.Conn
	isEagler bool
	// used if isEagler
	ws *websocket.Conn

	reader *bufio.Reader
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
}
func NewEaglerConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		conn:     nil,
		isEagler: true,
		ws:       conn,
		reader:   nil,
	}
}

// ReadMsg waits for a new message with timeout
// https://minecraft.wiki/w/Protocol?oldid=2772385#Without_compression
func (c *Connection) ReadMsg(timeout time.Duration) (m packet.Message, err error) {
	c.conn.SetReadDeadline(time.Now().Add(timeout))
	// read packet length from socket
	length, err := ds.DecodeVarIntFromReader(c.reader)
	if err != nil {
		return m, fmt.Errorf("failed to read packet length: %w", err)
	}

	var body = make([]byte, length)
	_, err = io.ReadFull(c.reader, body)
	if err != nil {
		return m, fmt.Errorf("Failed reading for full packet body: %w", err)
	}
	packetId, n, err := ds.DecodeVarInt(body)
	if err != nil {
		return m, fmt.Errorf("failed to decode packet id: %w", err)
	}
	m.PacketId = packet.PacketId(packetId)
	m.Data = body[n:]
	return
}

// Encode and WriteMessage to the socket
func (c *Connection) WriteMessage(m packet.Message) error {
	_, err := c.conn.Write(m.Encode())
	if err != nil {
		return fmt.Errorf("failed to write message to socket: %w", err)
	}
	return nil
}
func (c *Connection) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return c.ws.Close(websocket.StatusNormalClosure, "connection closed")
}
func (c *Connection) IsLoggedIn() bool {
	return c.isLoggedIn.Load()
}
