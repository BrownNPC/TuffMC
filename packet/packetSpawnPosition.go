package packet

import "tuff/ds"

const SpawnPositionPacketId PacketId = 0x46

// https://minecraft.wiki/w/Protocol?oldid=2772385#Spawn_Position
//
// Sent by the server after login to specify the coordinates of the spawn point (the point at which players spawn at,
// and which the compass points to).
//
// It can be sent at any time to update the point compasses point at.
func EncodeSpawnPositionPacket(X, Y, Z int32) Message {
	return Message{
		PacketId: SpawnPositionPacketId,
		Data:     ds.EncodePosition(X, Y, Z),
	}
}
