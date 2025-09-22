package packet

import (
	"fmt"
	"tuff/ds"
)

type StatusResponse struct {
	JSONResponse string
}

const StatusResponsePacketId PacketId = 0x00

func EncodeStatusResponse(playerCount int, description string, favicon string) Message {
	status := fmt.Sprintf(statusJson, playerCount, description, favicon)
	statusEncoded := ds.WriteString(status)
	return Message{ID: StatusResponsePacketId, Data: statusEncoded}
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
