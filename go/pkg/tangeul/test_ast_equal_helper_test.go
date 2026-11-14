//ff:func feature=tangl type=util control=sequence
//ff:what TestEqualIgnoringLine — EqualIgnoringLine reports equal, unequal, and both-nil Documents
package tangeul

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEqualIgnoringLine exercises the exported entry point directly (it has no
// branches of its own; it only wraps reflect.ValueOf + equalIgnoringLine): equal
// Documents that differ only in Line, Documents that differ in Subject, and the
// pointer-nil case both astPair-style callers never pass.
func TestEqualIgnoringLine(t *testing.T) {
	a := &ast.Document{Subject: "transfer", Cases: []ast.Case{{Name: "can withdraw", Line: 3}}}
	b := &ast.Document{Subject: "transfer", Cases: []ast.Case{{Name: "can withdraw", Line: 9}}}
	if !EqualIgnoringLine(a, b) {
		t.Fatal("expected true: Documents differ only in Line")
	}

	c := &ast.Document{Subject: "deposit", Cases: []ast.Case{{Name: "can withdraw", Line: 3}}}
	if EqualIgnoringLine(a, c) {
		t.Fatal("expected false: Subject differs")
	}

	var nilA, nilB *ast.Document
	if !EqualIgnoringLine(nilA, nilB) {
		t.Fatal("expected true: both nil Documents")
	}
}
