package tuff

import (
	"tuff/ds"
	"tuff/tuff/packets"
)

type Message struct {
	// Length of packet data + length of the packet ID
	ID packets.PacketID
	//Depends on the connection state and packet ID, see the sections below
	// https://minecraft.wiki/w/Protocol?oldid=2772385#Packet_format
	Data []byte
}

func ReadMessage(b []byte) (p Message, err error) {
	length, n, err := ds.ReadVarInt(b)
	if err != nil {
		return
	}
	b = b[n:] // skip length field
	packetId, n, err := ds.ReadVarInt(b)
	if err != nil {
		return
	}

	return Message{
		ID:   packets.PacketID(packetId),
		Data: b[n:length],
	}, nil
}
