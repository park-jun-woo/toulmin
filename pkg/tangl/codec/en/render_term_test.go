//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderTerm — renderTerm covers the prefix, Compare/Not/Either/default switch, and Not's inline-Compare sub-branch
package en

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderTerm exercises: the prefix present/absent head, a Compare term, a Not
// wrapping an inline Compare (single line "not `f` ..."), a Not wrapping a
// non-Compare (multi-line via renderCondList), an Either group, and an
// unrecognized Expr type (the defensive default branch).
func TestRenderTerm(t *testing.T) {
	cmp := ast.Compare{Field: "a", Op: "equals", Value: "1"}

	// No prefix, Compare term.
	got := renderTerm("", cmp, 0)
	want := []string{"- `a` equals 1"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("no-prefix Compare: got %v, want %v", got, want)
	}

	// Prefix present, Compare term.
	got = renderTerm("and", cmp, 0)
	want = []string{"- and `a` equals 1"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("prefix Compare: got %v, want %v", got, want)
	}

	// Not wrapping an inline Compare: single "not ..." line.
	got = renderTerm("", ast.Not{Term: cmp}, 0)
	want = []string{"- not `a` equals 1"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Not(Compare): got %v, want %v", got, want)
	}

	// Not wrapping a non-Compare (a Logic chain): head line plus a nested
	// condition list at depth+1.
	b := ast.Compare{Field: "b", Op: "equals", Value: "2"}
	logic := ast.Logic{Op: "and", Terms: []ast.Expr{cmp, b}}
	got = renderTerm("or", ast.Not{Term: logic}, 0)
	want = []string{"- or not", "  - `a` equals 1", "  - and `b` equals 2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Not(Logic): got %v, want %v", got, want)
	}

	// Either group: head line plus its first term's condition list.
	got = renderTerm("", ast.Either{Terms: []ast.Expr{logic}}, 1)
	want = []string{"  - either", "    - `a` equals 1", "    - and `b` equals 2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Either: got %v, want %v", got, want)
	}

	// Unrecognized Expr type: the defensive default branch.
	got = renderTerm("", logic, 0)
	want = []string{"- `?`"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("default: got %v, want %v", got, want)
	}
}
