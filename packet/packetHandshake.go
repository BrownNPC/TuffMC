package packet

import (
	"encoding/binary"
	"tuff/ds"
)

type ConnectionState = int

const (
	StateStatus ConnectionState = 1
	StateLogin  ConnectionState = 2
	StatePlay   ConnectionState = 3
)

const HandshakePacketID PacketId = 0x00

// https://minecraft.wiki/w/Java_Edition_protocol/Packets#Handshake
// https://minecraft.wiki/w/Protocol?oldid=2772385#Handshake
// C->S
type HandshakePacket struct {
	//The version that the client plans on using to connect to the server
	ProtocolVersion int
	// Hostname or IP, e.g. localhost or 127.0.0.1, that was used to connect.
	ServerAddress string
	// Default is 25565.
	ServerPort uint16
	// 1 for status, 2 for login
	NextState ConnectionState
}

func DecodeHandshakePacket(data []byte) (p HandshakePacket, err error) {
	var n int
	p.ProtocolVersion, n, err = ds.DecodeVarInt(data)
	if err != nil {
		return
	}
	data = data[n:]
	p.ServerAddress, n, err = ds.DecodeString(data)
	if err != nil {
		return
	}
	data = data[n:]

	p.ServerPort = binary.BigEndian.Uint16(data)
	// uint16 = 2 bytes
	data = data[2:]

	p.NextState, _, err = ds.DecodeVarInt(data)
	return
}
