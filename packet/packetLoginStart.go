package packet

import (
	"fmt"
	"tuff/ds"
)

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login_Start
// C->S
type LoginStartPacket struct {
	//Player's Username
	PlayerUsername string
}

const LoginStartPacketID PacketId = 0x00

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login_Start
// C->S
func DecodeLoginStartPacket(msg Message) (p LoginStartPacket, err error) {
	if msg.PacketId != LoginStartPacketID {
		err = fmt.Errorf("LoginStartPacket: Incorrect packet id. Expected %v Got %v", LoginStartPacketID, msg.PacketId)
		return
	}
	p.PlayerUsername, _, err = ds.DecodeString(msg.Data)
	return
}
