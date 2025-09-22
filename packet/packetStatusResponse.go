package packet

import (
	"fmt"
	"tuff/ds"
)

type StatusResponse struct {
	JSONResponse string
}

const StatusResponsePacketId PacketId = 0x00

type StatusResponseConfig struct {
	PlayerCount int
	Description string
	// base64 encoded image
	Favicon string
}

func EncodeStatusResponse(cfg StatusResponseConfig) Message {
	status := fmt.Sprintf(statusJson, cfg.PlayerCount, cfg.Description, cfg.Favicon)
	statusEncoded := ds.WriteString(status)
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
