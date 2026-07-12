//ff:func feature=tangl type=util control=sequence
//ff:what TestQuantOpToAST — QuantAny maps back to "any", an unknown QuantOp code is a hard error
package codebook

import "testing"

// TestQuantOpToAST exercises both branches of QuantOpToAST: the sole assigned
// code QuantAny maps back to "any", and any other code is a hard error —
// never a silent fallback.
func TestQuantOpToAST(t *testing.T) {
	op, err := QuantOpToAST(QuantAny)
	if err != nil || op != "any" {
		t.Fatalf("QuantOpToAST(QuantAny): got %q err=%v, want %q", op, err, "any")
	}
	if _, err := QuantOpToAST(QuantOp(0xFFFF)); err == nil {
		t.Error("QuantOpToAST(unknown): expected error")
	}
}
