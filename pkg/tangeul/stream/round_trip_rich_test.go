//ff:func feature=tangl type=util control=sequence
//ff:what TestRoundTripRich — every-section fixture round-trips on AST/Source/Blocks and re-writes identically
package stream

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestRoundTripRich runs the all-sections fixture (prose, Definitions const +
// struct, Rules, in-case with/checking/qualifier, a gated do, endpoint checks,
// both Internal kinds) through write, read and re-write, asserting AST, Source
// and Blocks equality, byte-identical re-serialization, and a valid tiling.
func TestRoundTripRich(t *testing.T) {
	d := buildRichDoc()

	b1, err := WriteDoc(d)
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	got, err := ReadDoc(b1)
	if err != nil {
		t.Fatalf("ReadDoc: %v", err)
	}
	if !reflect.DeepEqual(got.AST, d.AST) {
		t.Errorf("AST mismatch:\n got %+v\nwant %+v", got.AST, d.AST)
	}
	if !bytes.Equal(got.Source, d.Source) {
		t.Errorf("Source mismatch")
	}
	if !reflect.DeepEqual(got.Blocks, d.Blocks) {
		t.Errorf("Blocks mismatch:\n got %+v\nwant %+v", got.Blocks, d.Blocks)
	}
	if err := doc.VerifyTiling(got); err != nil {
		t.Errorf("VerifyTiling: %v", err)
	}
	b2, err := WriteDoc(got)
	if err != nil {
		t.Fatalf("WriteDoc(recovered): %v", err)
	}
	if !bytes.Equal(b1, b2) {
		t.Errorf("re-write not byte-identical")
	}
}
