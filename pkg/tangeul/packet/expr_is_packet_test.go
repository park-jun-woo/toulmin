//ff:func feature=tangl type=util control=sequence
//ff:what TestExprPacketIsPacket — ExprPacket satisfies the sealed Packet interface
package packet

import "testing"

// TestExprPacketIsPacket checks ExprPacket implements Packet via its
// isPacket marker method.
func TestExprPacketIsPacket(t *testing.T) {
	var p Packet = ExprPacket{}
	p.isPacket()
}
