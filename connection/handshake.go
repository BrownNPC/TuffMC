package connection

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"time"
	"tuff/packet"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

func timeout(t time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), t)
	return ctx
}

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login
func (conn *Connection) HandleHandshake(cfg packet.StatusResponsePacketConfig) error {
	if conn.isEagler {
		return conn.eaglerHandshake(cfg)
	}
	return conn.javaHandshake(cfg)
}

func (conn *Connection) eaglerHandshake(cfg packet.StatusResponsePacketConfig) error {
	typ, b, err := conn.ws.Read(timeout(time.Second * 10))
	_ = b
	if err != nil {
		return fmt.Errorf("websocket error: %w", err)
	}
	// server list ping
	if typ == websocket.MessageText {
		fmt.Println(string(b))
		const eaglerStatusJson = `
			{"name":"%s",
			"brand":"TuffMC",
			"vers":"EaglercraftXServer/1.0.7",
			"cracked":true,
			"time":%d,
			"type":"motd",
			"data":{"cache":true,
			"motd":["%s"],
			"icon":true,
			"online":%d,
			"max":67,"players":[]}}`
		status := fmt.Sprintf(eaglerStatusJson,
			cfg.Description,        // %s - name
			time.Now().UnixMilli(), // %d - time
			cfg.Description,        // %s - motd
			// cfg.Favicon,            // %s - icon
			cfg.PlayerCount, // %d - online
		)
		return conn.ws.Write(timeout(time.Second*10), websocket.MessageText, []byte(status))
	}
	if typ != websocket.MessageBinary {
		return fmt.Errorf("Expected binary message, got %s", typ.String())
	}
	//the client has sent us their handshake request. we only accept 1.12.2 and
	// eagler handshake v2 clients, so we dont care what the request says.
	// Just blindly send an acknowledgement of a v2 1.12.2 client.
	err = conn.ws.Write(timeout(time.Second*10), websocket.MessageBinary,
		packet.EncodeEaglerHandshakeAckPacket(),
	)
	if err != nil {
		return fmt.Errorf("failed to write handshake response: %w", err)
	}
	typ, b, err = conn.ws.Read(timeout(time.Second * 10))
	if err != nil {
		return fmt.Errorf("failed to read Username packet: %w", err)
	}
	// get username
	var buf = bytes.NewBuffer(b)
	buf.ReadByte() //ignore packet id
	unameLength, _ := buf.ReadByte()
	username := buf.Next(int(unameLength))
	conn.Username = string(username)

	// build packet for login ack
	buf.Reset()
	buf.WriteByte(5) //packet id
	// username
	buf.WriteByte(unameLength)
	buf.WriteString(string(username))
	//player uuid
	playerUUID := uuid.NewMD5(packet.Namespace, username)
	buf.Write(playerUUID[:])
	// send login ack
	err = conn.ws.Write(timeout(time.Second*10), websocket.MessageBinary,
		buf.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("failed to write login ack: %w", err)
	}
	//Client may send profile / skin packet packetId == 7:
	// we want to ignore it.
	typ, b, err = conn.ws.Read(timeout(time.Second * 10))
	if err != nil {
		return fmt.Errorf("failed to wait for login state packet: %w", err)
	}
	if len(b) < 1 {
		return fmt.Errorf("too small packet recieved")
	}
	// if b[0] == 7 {
	// 	typ, b, err = conn.ws.Read(timeout(time.Second * 10))
	// 	if err != nil {
	// 		return fmt.Errorf("failed to wait for login state packet")
	// 	}
	// }
	// if b[0] != 8 {
	// 	return fmt.Errorf("expected play stage packet")
	// }
	conn.conn = websocket.NetConn(context.Background(), conn.ws, websocket.MessageBinary)
	conn.reader = bufio.NewReader(conn.conn)
	conn.isLoggedIn.Store(true)
	err = conn.ws.Write(timeout(time.Second*10), websocket.MessageBinary, []byte{9})
	// _, err = conn.conn.Write([]byte{9})
	if err != nil {
		return fmt.Errorf("failed to send final login packet: %w", err)
	}
	return nil
}
func (conn *Connection) javaHandshake(cfg packet.StatusResponsePacketConfig) error {
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
