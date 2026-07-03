//ff:func feature=tangl type=util control=sequence
//ff:what NodePacket.isPacket — seals NodePacket into the Packet interface
package packet

// isPacket seals NodePacket into the Packet interface (marker method).
func (NodePacket) isPacket() {}
