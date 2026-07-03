//ff:func feature=tangl type=util control=sequence
//ff:what TestControlPacketIsPacket — ControlPacket satisfies the sealed Packet interface
package packet

import "testing"

// TestControlPacketIsPacket checks ControlPacket implements Packet via its
// isPacket marker method.
func TestControlPacketIsPacket(t *testing.T) {
	var p Packet = ControlPacket{}
	p.isPacket()
}
