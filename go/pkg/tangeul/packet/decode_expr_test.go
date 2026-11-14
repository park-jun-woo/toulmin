//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeExpr — DecodeExpr round trips every tree shape and rejects every malformed record
package packet

import "testing"

// TestDecodeExpr gathers the round-trip and error suites so one test exercises
// every branch of DecodeExpr.
func TestDecodeExpr(t *testing.T) {
	t.Run("roundtrip", TestExprRoundTrip)
	t.Run("errors", TestDecodeExprErrors)
}
