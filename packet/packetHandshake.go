package packet

import (
	"encoding/binary"
	"tuff/ds"
)

type LoginState = int

const (
	StateStatus LoginState = 1
	StateLogin  LoginState = 2
)

const HandshakePacketID PacketId = 0x00

// https://minecraft.wiki/w/Java_Edition_protocol/Server_List_Ping#Handshake
// https://minecraft.wiki/w/Protocol?oldid=2772385#Handshake
type Handshake struct {
	//The version that the client plans on using to connect to the server
	ProtocolVersion int
	// Hostname or IP, e.g. localhost or 127.0.0.1, that was used to connect.
	ServerAddress string
	// Default is 25565.
	ServerPort uint16
	// 1 for status, 2 for login
	NextState LoginState
}

func DecodeHandshake(data []byte) (h Handshake, err error) {
	var n int
	h.ProtocolVersion, n, err = ds.ReadVarInt(data)
	if err != nil {
		return
	}
	data = data[n:]
	h.ServerAddress, n, err = ds.ReadString(data)
	if err != nil {
		return
	}
	data = data[n:]

	h.ServerPort = binary.BigEndian.Uint16(data)
	// uint16 = 2 bytes
	data = data[2:]

	h.NextState, _, err = ds.ReadVarInt(data)
	return
}
