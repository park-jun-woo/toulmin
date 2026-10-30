//ff:func feature=tangl type=model control=sequence
//ff:what TestCollectSpecNameTerms — terms are visited in order sharing one seen set; no terms is a no-op
package ast

import (
	"reflect"
	"testing"
)

func TestCollectSpecNameTerms(t *testing.T) {
	seen := map[string]bool{}
	terms := []Expr{
		Compare{Lhs: SpecRef{Steps: []string{"b"}}},
		Either{Terms: []Expr{Compare{Rhs: SpecRef{Steps: []string{"a"}}}}},
		Compare{Lhs: SpecRef{Steps: []string{"b"}}},
	}
	got := collectSpecNameTerms(terms, seen, nil)
	if want := []string{"b", "a"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("collectSpecNameTerms = %v, want %v", got, want)
	}
	if !seen["a"] || !seen["b"] {
		t.Fatalf("seen = %v, want a and b marked", seen)
	}
	if got := collectSpecNameTerms(nil, map[string]bool{}, []string{"x"}); !reflect.DeepEqual(got, []string{"x"}) {
		t.Fatalf("empty terms = %v, want input unchanged", got)
	}
}
