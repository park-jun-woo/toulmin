//ff:func feature=tangl type=util control=sequence
//ff:what ControlPacket.isPacket — seals ControlPacket into the Packet interface
package packet

// isPacket seals ControlPacket into the Packet interface (marker method).
func (ControlPacket) isPacket() {}
