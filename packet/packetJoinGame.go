package packet

import (
	"slices"
	"tuff/ds"
)

// https://minecraft.wiki/w/Protocol?oldid=2772385#Join_Game
// S->C
type JoinGamePacketConfig struct {
	// The player's Entity ID (EID)
	EntityID int32
	//0: Survival, 1: Creative, 2: Adventure, 3: Spectator. Bit 3 (0x8) is the hardcore flag.
	Gamemode byte
	//1: Nether, 0: Overworld, 1: End; also, note that this is not a VarInt but instead a regular int.
	Dimension int32
	//0: peaceful, 1: easy, 2: normal, 3: hard
	Difficulty byte
	//Was once used by the client to draw the player list, but now is ignored
	MaxPlayers byte
	//default, flat, largeBiomes, amplified, default_1_1
	LevelType string
	//If true, a Notchian client shows reduced information on the debug screen (F3).
	// For servers in development, this should almost always be false.
	ReducedDebugInfo bool
}

const JoinGamePacketId PacketId = 0x23

// https://minecraft.wiki/w/Protocol?oldid=2772385#Join_Game
// S->C
func EncodeJoinGamePacket(cfg JoinGamePacketConfig) Message {
	return Message{
		PacketId: JoinGamePacketId,
		Data: slices.Concat(
			ds.EncodeInt(cfg.EntityID),
			[]byte{cfg.Gamemode},
			ds.EncodeInt(cfg.Dimension),
			[]byte{cfg.Difficulty,
				cfg.MaxPlayers},
			ds.EncodeString(cfg.LevelType),
			ds.EncodeBool(cfg.ReducedDebugInfo),
		),
	}
}
