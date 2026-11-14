//ff:func feature=tangl type=model control=sequence
//ff:what TestRuleSpecNames — rule conditions yield first-step SpecRef names once, in source order
package ast

import (
	"reflect"
	"testing"
)

func TestRuleSpecNames(t *testing.T) {
	rule := InlineRule{Name: "r", Cond: Logic{Op: "and", Terms: []Expr{
		Compare{Lhs: SpecRef{Steps: []string{"limit", "Value"}}, Op: "is at most", Rhs: SpecRef{Steps: []string{"cap"}}},
		Not{Term: Compare{Lhs: Path{Steps: []string{"score"}}, Op: "is in", Rhs: SpecRef{Steps: []string{"limit"}}}},
	}}}
	if got, want := RuleSpecNames(rule), []string{"limit", "cap"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("RuleSpecNames = %v, want %v", got, want)
	}
	if got := RuleSpecNames(InlineRule{Name: "host", Using: &Ref{Alias: "a", Name: "f"}}); got != nil {
		t.Fatalf("host rule without Cond = %v, want nil", got)
	}
}
