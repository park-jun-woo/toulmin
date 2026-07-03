//ff:func feature=tangl type=util control=sequence
//ff:what TestTextPacketIsPacket — TextPacket satisfies the sealed Packet interface
package packet

import "testing"

// TestTextPacketIsPacket checks TextPacket implements Packet via its
// isPacket marker method.
func TestTextPacketIsPacket(t *testing.T) {
	var p Packet = TextPacket{}
	p.isPacket()
}
