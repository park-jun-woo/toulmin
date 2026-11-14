//ff:func feature=tangl type=util control=sequence
//ff:what TestRoundTrip — ReadDoc(WriteDoc(d)) equals d (AST/Source/Blocks) and re-write is byte-identical
package stream

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestRoundTrip serializes the §10 account-transfer fixture, reads it back, and
// asserts the recovered Doc matches on AST (reflect.DeepEqual — the fixture uses
// Line 0 throughout, which TANGEUL does not encode), Source bytes and Blocks. It
// then re-serializes the recovered Doc and requires byte identity (§9 norm), and
// checks the recovered Blocks still tile the Source.
func TestRoundTrip(t *testing.T) {
	d := buildTransferDoc()

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
		t.Errorf("Source mismatch: got %q want %q", got.Source, d.Source)
	}
	if !reflect.DeepEqual(got.Blocks, d.Blocks) {
		t.Errorf("Blocks mismatch:\n got %+v\nwant %+v", got.Blocks, d.Blocks)
	}
	if err := doc.VerifyTiling(got); err != nil {
		t.Errorf("VerifyTiling(recovered): %v", err)
	}

	b2, err := WriteDoc(got)
	if err != nil {
		t.Fatalf("WriteDoc(recovered): %v", err)
	}
	if !bytes.Equal(b1, b2) {
		t.Errorf("re-write not byte-identical: %d vs %d bytes", len(b1), len(b2))
	}
}
