package connection

import (
	"errors"
	"tuff/packet"
)

// https://minecraft.wiki/w/Protocol?oldid=2772385#Login
func HandleLogin(m packet.Message, cfg packet.StatusResponseConfig, conn Connection) (c Connection, err error) {
	handshake, err := packet.DecodeHandshake(m.Data)
	if err != nil {
		return
	}

	switch handshake.NextState {
	// If status request, return status and close the connection.
	case packet.StateStatus:
		statusMsg := packet.EncodeStatusResponse(cfg)
		_, err = conn.Write(statusMsg.Encode())
		if err != nil {
			return
		}
		return Connection{}, errors.New("status requested, cannot continue")
	case packet.StateLogin:

	}
}
