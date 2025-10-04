package packet

import (
	"slices"
	"tuff/ds"

	"github.com/google/uuid"
)

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login_Success
const LoginSuccessPacketID PacketId = 0x02

// used to calculate player UUID.
//
//	playerUUID := uuid.NewMD5(Namespace, []byte(Username))
var Namespace = uuid.MustParse("00000000-0000-0000-0000-000000000001")

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login_Success
func EncodeLoginSuccessPacket(Username string) Message {
	playerUUID := uuid.NewMD5(Namespace, []byte(Username)).
		String()
	return Message{
		PacketId: LoginSuccessPacketID,
		Data: slices.Concat(
			ds.EncodeString(playerUUID),
			ds.EncodeString(Username),
		)}
}
