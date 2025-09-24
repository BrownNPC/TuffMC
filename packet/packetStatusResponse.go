package packet

import (
	"fmt"
	"tuff/ds"
)
const StatusResponsePacketId PacketId = 0x00

// https://minecraft.wiki/w/Protocol?oldid=2772385#Status
// S -> C
type StatusResponsePacketConfig struct {
	PlayerCount int
	Description string
	// base64 encoded image
	Favicon string
}

// https://minecraft.wiki/w/Protocol?oldid=2772385#Status
// S -> C
func EncodeStatusResponsePacket(cfg StatusResponsePacketConfig) Message {
	status := fmt.Sprintf(statusJson, cfg.PlayerCount, cfg.Description, cfg.Favicon)
	statusEncoded := ds.EncodeString(status)
	return Message{PacketId: StatusResponsePacketId, Data: statusEncoded}
}

const statusJson string = `{
    "version": {
        "name": "1.12.2",
        "protocol": 340
    },
    "players": {
        "max": 67,
        "online": %d
    },
    "description": {
        "text": "%s"
    },
    "favicon": "%s"
}`
