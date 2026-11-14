//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeExpr — EncodeExpr golden bytes, round trips and error propagation
package packet

import "testing"

// TestEncodeExpr gathers the byte-identity, round-trip and error suites so one
// test exercises every branch of EncodeExpr.
func TestEncodeExpr(t *testing.T) {
	t.Run("bytes", TestExprByteIdentity)
	t.Run("roundtrip", TestExprRoundTrip)
	t.Run("errors", TestEncodeExprErrors)
}
