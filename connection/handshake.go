package connection

import (
	"fmt"
	"time"
	"tuff/packet"
)

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login
func (conn *Connection) HandleHandshake(cfg packet.StatusResponseConfig) (bool, error) {
	m, err := conn.ReadMsg(time.Second * 10)
	if err != nil {
		return false, fmt.Errorf("failed to read handshake message: %w", err)
	}
	// C->S handshake
	handshake, err := packet.DecodeHandshake(m.Data)
	if err != nil {
		return false, fmt.Errorf("failed to decode handshake data: %w", err)
	}
	// Login or Status handshake?
	switch handshake.NextState {
	// If status request, return status and close the connection.
	case packet.StateStatus:
		statusMsg := packet.EncodeStatusResponse(cfg)
		err = conn.WriteMessage(statusMsg)
		if err != nil {
			return false, fmt.Errorf("failed to write status response: %w", err)
		}
		// Ping packet
		m, err = conn.ReadMsg(time.Second * 10)
		if err != nil {
			return false, fmt.Errorf("failed to read ping message: %w", err)
		}
		// ping pong requires us to resend their packet
		err = conn.WriteMessage(m)
		if err != nil {
			return false, fmt.Errorf("failed to respond with pong message: %w", err)
		}
		// status handshake complete, close connection.
		conn.conn.Close()
		return false, nil
	case packet.StateLogin:
	}
	return false, fmt.Errorf("unhandled packet state %v", handshake.NextState)
}
