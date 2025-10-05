package packet

import "bytes"

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
	clientVersionStrLength, _ := buf.ReadByte()
	clientVersionString := buf.Next(int(clientVersionStrLength))
}
