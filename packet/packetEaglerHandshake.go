package packet

import (
	"bytes"
	"encoding/binary"
)

func EncodeEaglerHandshakeAckPacket() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(0x02)                              // ACK packet ID
	binary.Write(buf, binary.BigEndian, uint16(2))   // Eaglercraft version 2
	binary.Write(buf, binary.BigEndian, uint16(340)) // MC protocol 340 (1.12.2)

	brand := []byte("TuffMC")
	buf.WriteByte(byte(len(brand)))
	buf.Write(brand)

	version := []byte("1.12.2")
	buf.WriteByte(byte(len(version)))
	buf.Write(version)

	buf.WriteByte(0x00)                            // padding
	binary.Write(buf, binary.BigEndian, uint16(0)) // short padding

	return buf.Bytes()
}
