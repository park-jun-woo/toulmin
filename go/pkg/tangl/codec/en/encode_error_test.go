//ff:func feature=tangl type=parser control=sequence
//ff:what TestEncodeParserError — a parser error is propagated unchanged by Encode
package en

import (
	"strings"
	"testing"
)

// TestEncodeParserError checks that when the parser rejects the source (here a
// document with no tangl:Subject section), Encode returns that error and no Doc,
// rather than proceeding to tiling.
func TestEncodeParserError(t *testing.T) {
	src := []byte("## tangl:Cases\n- in case of `c`\n  - `r` is a general rule\n")
	d, err := ENCodec{}.Encode(src, "bad.md")
	if err == nil {
		t.Fatal("Encode accepted source with no Subject section")
	}
	if d != nil {
		t.Fatalf("Encode returned a Doc alongside an error: %+v", d)
	}
	if !strings.Contains(err.Error(), "Subject") {
		t.Fatalf("Encode error = %v, want a Subject-required parser error", err)
	}
}
