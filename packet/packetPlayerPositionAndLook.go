package packet

import (
	"slices"
	"tuff/ds"
)

// https://minecraft.wiki/w/Protocol?oldid=2772385#Player_Position_And_Look_(clientbound)
type PlayerPositionAndLookConfig struct {
	// Absolute or relative position, depending on Flags
	X, Y, Z float64
	// Absolute or relative rotation on the X axis, in degrees
	Pitch float32
	// Absolute or relative rotation on the Y axis, in degrees
	Yaw float32
	// Bit field, see below
	//
	// X	    0x01
	// Y	    0x02
	// Z	    0x04
	// Y_ROT    0x08
	// X_ROT	0x10
	//
	//
	// It's a bitfield, X/Y/Z/Y_ROT/X_ROT. If X is set, the x value is relative and not absolute.
	Flags byte
	// Client should confirm this packet with Teleport Confirm containing the same Teleport ID
	TeleportId uint
}

// https://minecraft.wiki/w/Protocol?oldid=2772385#Player_Position_And_Look_(clientbound)
const PlayerPositionAndLookPacketId PacketId = 0x2F

// https://minecraft.wiki/w/Protocol?oldid=2772385#Player_Position_And_Look_(clientbound)
// Updates the player's position on the server. This packet will also close the “Downloading Terrain” screen when joining/respawning.
//
// If the distance between the last known position of the player on the server and the new position set by this packet is greater than 100 meters, the client will be kicked for “You moved too quickly :( (Hacking?)”.
//
// Also if the fixed-point number of X or Z is set greater than 3.2E7D the client will be kicked for “Illegal position”.
//
// Yaw is measured in degrees, and does not follow classical trigonometry rules. The unit circle of yaw on the XZ-plane starts at (0, 1) and turns counterclockwise, with 90 at (-1, 0), 180 at (0, -1) and 270 at (1, 0). Additionally, yaw is not clamped to between 0 and 360 degrees; any number is valid, including negative numbers and numbers greater than 360.
//
// Pitch is measured in degrees, where 0 is looking straight ahead, -90 is looking straight up, and 90 is looking straight down.
func EncodePlayerPositionAndLookPacket(cfg PlayerPositionAndLookConfig) Message {
	return Message{
		PacketId: PlayerPositionAndLookPacketId,
		Data: slices.Concat(
			ds.EncodeDouble(cfg.X),
			ds.EncodeDouble(cfg.Y),
			ds.EncodeDouble(cfg.Z),
			ds.EncodeFloat(cfg.Yaw),
			ds.EncodeFloat(cfg.Pitch),
			[]byte{cfg.Flags},
			ds.EncodeVarInt(cfg.TeleportId),
		),
	}
}
