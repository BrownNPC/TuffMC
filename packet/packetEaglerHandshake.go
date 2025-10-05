package packet

import (
	"bytes"
	"fmt"
)

type EaglerHandshakeRequestPacket struct {
	PacketId                             byte
	EaglerCraftVersion, MinecraftVersion byte
	ClientBrand                          string
	ClientVersionString                  string
}

func DecodeEaglerHandshakeRequestPacket(b []byte) (EaglerHandshakeRequestPacket, error) {
	var buf = bytes.NewBuffer(b)
	packetId, _ := buf.ReadByte()

	eaglerCraftVersion, _ := buf.ReadByte()
	minecraftVersion, _ := buf.ReadByte()

	clientBrandStrLength, _ := buf.ReadByte()
	clientBrand := buf.Next(int(clientBrandStrLength))

	clientVersionStrLength, err := buf.ReadByte()
	clientVersionString := buf.Next(int(clientVersionStrLength))
	if err != nil {
		return EaglerHandshakeRequestPacket{}, fmt.Errorf("failed to read packet")
	}
	return EaglerHandshakeRequestPacket{
		PacketId:            packetId,
		EaglerCraftVersion:  eaglerCraftVersion,
		MinecraftVersion:    minecraftVersion,
		ClientBrand:         string(clientBrand),
		ClientVersionString: string(clientVersionString),
	}, nil
}
