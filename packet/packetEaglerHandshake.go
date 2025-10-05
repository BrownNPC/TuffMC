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

type EaglerHandshakeResponsePacket struct {
	PacketId            byte
	EaglerCraftVersion  byte
	ClientBrand         string
	ClientVersionString string
	// 3 bytes of padding required..
	_ [3]byte
}

func EncodeEaglerHandshakeResponsePacket() []byte {
	var b bytes.Buffer
	// packet id
	b.WriteByte(1)
	// eaglercraft version
	b.WriteByte(1)
	//Brand
	b.WriteByte(byte(len("TuffMC")))
	b.WriteString("TuffMC")

	// Version
	b.WriteByte(1)
	b.WriteRune('1')

	// 3 bytes padding required..
	b.Write([]byte{0, 0, 0})
	return b.Bytes()
}
