package connection

import (
	"net"
	"time"
	"tuff/packet"
)

type Connection struct {
	IsLoggedIn bool
	// underlying socket
	conn net.Conn
	buf  []byte
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		IsLoggedIn: false,
		conn:       conn,
		// 2 MiB buffer
		buf: make([]byte, 2*1024*1024),
	}
}

// ReadMsg waits for a new message with timeout
func (c *Connection) ReadMsg(timeout time.Duration) (m packet.Message, err error) {
	c.conn.SetReadDeadline(time.Now().Add(timeout))
	n, err := c.conn.Read(c.buf)
	if err != nil {
		return
	}
	m, err = packet.DecodeMessage(c.buf[:n])
	return
}
