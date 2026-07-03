//ff:func feature=tangl type=util control=sequence
//ff:what EdgePacket.isPacket — seals EdgePacket into the Packet interface
package packet

// isPacket seals EdgePacket into the Packet interface (marker method).
func (EdgePacket) isPacket() {}
