//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeCondPrefix — tests takeCondPrefix for and/or/unmatched branches
package parser

import "testing"

func TestTakeCondPrefix(t *testing.T) {
	t.Run("And", func(t *testing.T) {
		op, rest, ok := takeCondPrefix("and `x` equals 1")
		if !ok {
			t.Fatal("expected ok=true for 'and'")
		}
		if op != "and" {
			t.Errorf("expected op='and', got %q", op)
		}
		if rest != " `x` equals 1" {
			t.Errorf("expected rest=' `x` equals 1', got %q", rest)
		}
	})

	t.Run("Or", func(t *testing.T) {
		op, rest, ok := takeCondPrefix("or `x` equals 1")
		if !ok {
			t.Fatal("expected ok=true for 'or'")
		}
		if op != "or" {
			t.Errorf("expected op='or', got %q", op)
		}
		if rest != " `x` equals 1" {
			t.Errorf("expected rest=' `x` equals 1', got %q", rest)
		}
	})

	t.Run("Unmatched", func(t *testing.T) {
		op, rest, ok := takeCondPrefix("`x` equals 1")
		if ok {
			t.Fatal("expected ok=false for an unrecognized prefix")
		}
		if op != "" {
			t.Errorf("expected empty op, got %q", op)
		}
		if rest != "`x` equals 1" {
			t.Errorf("expected rest='`x` equals 1', got %q", rest)
		}
	})
}
