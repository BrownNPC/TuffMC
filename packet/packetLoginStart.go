package packet

import "tuff/ds"

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login_Start
// C->S
type LoginStartPacket struct {
	//Player's Username
	PlayerUsername string
}

const LoginStartPacketID PacketId = 0x00

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login_Start
// C->S
func DecodeLoginStartPacket(data []byte) (p LoginStartPacket, err error) {
	p.PlayerUsername, _, err = ds.DecodeString(data)
	return
}
