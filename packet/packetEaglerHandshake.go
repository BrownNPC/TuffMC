package packet

type EaglerHandshakeRequestPacket struct {
	PacketId                             byte
	EaglerCraftVersion, MinecraftVersion byte
	ClientBrand                          string
	ClientVersionString                  string
}
func DecodeEagler
