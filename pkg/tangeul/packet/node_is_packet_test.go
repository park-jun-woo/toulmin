//ff:func feature=tangl type=util control=sequence
//ff:what TestNodePacketIsPacket — NodePacket satisfies the sealed Packet interface
package packet

import "testing"

// TestNodePacketIsPacket checks NodePacket implements Packet via its
// isPacket marker method.
func TestNodePacketIsPacket(t *testing.T) {
	var p Packet = NodePacket{}
	p.isPacket()
}
