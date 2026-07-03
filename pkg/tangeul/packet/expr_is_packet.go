//ff:func feature=tangl type=util control=sequence
//ff:what ExprPacket.isPacket — seals ExprPacket into the Packet interface
package packet

// isPacket seals ExprPacket into the Packet interface (marker method).
func (ExprPacket) isPacket() {}
