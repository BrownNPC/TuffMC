package connection

import (
	"fmt"
	"time"
	"tuff/packet"
)

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login
func (conn *Connection) HandleHandshake(cfg packet.StatusResponsePacketConfig) error {
	pkt_handshake, err := conn.ReadMsg(time.Second * 10)
	if err != nil {
		return fmt.Errorf("failed to read handshake message: %w", err)
	}
	// C->S handshake
	handshake, err := packet.DecodeHandshakePacket(pkt_handshake)
	if err != nil {
		return fmt.Errorf("failed to decode handshake data: %w", err)
	}
	fmt.Printf("Handshake Received: %+v\n", handshake)
	// Login or Status handshake?
	switch handshake.NextState {
	// If status request, return status and close the connection.
	case packet.StateStatus:
		err = conn.WriteMessage(packet.EncodeStatusResponsePacket(cfg))
		if err != nil {
			return fmt.Errorf("failed to write status response: %w", err)
		}
		// Ping packet
		pkt_ping, err := conn.ReadMsg(time.Second * 10)
		if err != nil {
			return fmt.Errorf("failed to read ping message: %w", err)
		}
		// ping pong requires us to resend their packet
		err = conn.WriteMessage(pkt_ping)
		if err != nil {
			return fmt.Errorf("failed to respond with pong message: %w", err)
		}
		// status handshake complete, close connection.
		conn.conn.Close()
		return nil
		// C→S: Handshake with Next State set to 2 (login)
	case packet.StateLogin:
		// C→S: Login Start
		login, err := conn.ReadMsg(time.Second * 10)
		if err != nil {
			return fmt.Errorf("failed to recive login response: %w", err)
		}
		loginStartPacket, err := packet.DecodeLoginStartPacket(login)
		if err != nil {
			return fmt.Errorf("failed to decode login start packet: %w", err)
		}
		// S→C: Set Compression (optional)
		// TODO
		// S→C: Login Success
		loginSuccessResponseMsg := packet.EncodeLoginSuccessPacket(loginStartPacket.PlayerUsername)
		err = conn.WriteMessage(loginSuccessResponseMsg)
		if err != nil {
			return fmt.Errorf("failed to write login success: %w", err)
		}

		conn.State = packet.StatePlay
		conn.isLoggedIn.Store(true)
		return err
	}
	return fmt.Errorf("unhandled packet state %v", handshake.NextState)
}
