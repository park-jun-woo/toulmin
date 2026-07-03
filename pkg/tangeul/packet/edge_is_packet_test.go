//ff:func feature=tangl type=util control=sequence
//ff:what TestEdgePacketIsPacket — EdgePacket satisfies the sealed Packet interface
package packet

import "testing"

// TestEdgePacketIsPacket checks EdgePacket implements Packet via its
// isPacket marker method.
func TestEdgePacketIsPacket(t *testing.T) {
	var p Packet = EdgePacket{}
	p.isPacket()
}
