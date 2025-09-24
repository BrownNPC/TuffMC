package connection

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
	"tuff/packet"
)

type Connection struct {
	State      packet.ConnectionState
	isLoggedIn atomic.Bool
	// underlying socket
	conn net.Conn
	buf  []byte
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
		// 2 MiB buffer
		buf: make([]byte, 2*1024*1024),
	}
}

// ReadMsg waits for a new message with timeout
func (c *Connection) ReadMsg(timeout time.Duration) (m packet.Message, err error) {
	c.conn.SetReadDeadline(time.Now().Add(timeout))
	// read packet from socket
	n, err := c.conn.Read(c.buf)
	if err != nil {
		return m, fmt.Errorf("failed to read socket: %w", err)
	}
	// decode packet
	m, err = packet.DecodeMessage(c.buf[:n])
	if err != nil {
		return m, fmt.Errorf("failed to decode message: %w", err)
	}
	return
}

// Encode and WriteMessage to the socket
func (c *Connection) WriteMessage(m packet.Message) (err error) {
	_, err = c.conn.Write(m.Encode())
	if err != nil {
		return fmt.Errorf("failed to write message to socket: %w", err)
	}
	return
}
func (c *Connection) Close() error {
	return c.conn.Close()
}
func (c *Connection) IsLoggedIn() bool {
	return c.isLoggedIn.Load()
}
