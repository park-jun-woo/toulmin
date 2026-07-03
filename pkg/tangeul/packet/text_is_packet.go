//ff:func feature=tangl type=util control=sequence
//ff:what TextPacket.isPacket — seals TextPacket into the Packet interface
package packet

// isPacket seals TextPacket into the Packet interface (marker method).
func (TextPacket) isPacket() {}
