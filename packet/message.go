package packet

import (
	"slices"
	"tuff/ds"
)

type Message struct {
	// Length of packet data + length of the packet PacketId
	PacketId
	//Depends on the connection state and packet ID, see the sections below
	// https://minecraft.wiki/w/Protocol?oldid=2772385#Packet_format
	Data []byte
}

// https://minecraft.wiki/w/Protocol?oldid=2772385#Without_compression
// The base packet.
func DecodeMessage(b []byte) (_ Message, err error) {
	payloadLength, n, err := ds.DecodeVarInt(b)
	if err != nil {
		return
	}
	b = b[n:] // skip length field
	packetId, n, err := ds.DecodeVarInt(b)
	if err != nil {
		return
	}

	return Message{
		PacketId: PacketId(packetId),
		Data:     b[n:payloadLength],
	}, nil
}

// https://minecraft.wiki/w/Protocol?oldid=2772385#Without_compression
func (m Message) Encode() []byte {
	packetId := ds.EncodeVarInt(uint(m.PacketId))
	// Length of packet id+data
	length := len(packetId) + len(m.Data)
	encodedLength := ds.EncodeVarInt(uint(length))

	return slices.Concat(encodedLength, packetId, m.Data)
}
