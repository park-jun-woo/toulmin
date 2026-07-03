//ff:type feature=tangl type=interface
//ff:what Packet — the sealed interface over the four TANGEUL packet classes
package packet

// Packet is the sealed interface implemented by every TANGEUL packet type
// (NodePacket, EdgePacket, ControlPacket, TextPacket and ExprPacket). The
// unexported isPacket method keeps the set closed to this package so the stream
// layer can hold a heterogeneous []Packet in declaration order (§7.1, §7.2)
// while every concrete type stays known at compile time.
type Packet interface {
	isPacket()
}
