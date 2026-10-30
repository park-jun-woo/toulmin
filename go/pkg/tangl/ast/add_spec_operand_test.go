//ff:func feature=tangl type=model control=iteration dimension=1
//ff:what TestAddSpecOperand — only unseen SpecRefs with at least one step append their first step
package ast

import (
	"reflect"
	"testing"
)

func TestAddSpecOperand(t *testing.T) {
	cases := []struct {
		name string
		op   Operand
		seen map[string]bool
		want []string
	}{
		{"new spec ref", SpecRef{Steps: []string{"limit", "Value"}}, map[string]bool{}, []string{"limit"}},
		{"already seen", SpecRef{Steps: []string{"limit"}}, map[string]bool{"limit": true}, nil},
		{"empty steps", SpecRef{}, map[string]bool{}, nil},
		{"path operand", Path{Steps: []string{"limit"}}, map[string]bool{}, nil},
		{"call operand", Call{Fn: Ref{Name: "f"}}, map[string]bool{}, nil},
		{"nil operand", nil, map[string]bool{}, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := addSpecOperand(c.op, c.seen, nil)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("addSpecOperand = %v, want %v", got, c.want)
			}
			if len(c.want) == 1 && !c.seen[c.want[0]] {
				t.Fatalf("%q not marked seen", c.want[0])
			}
		})
	}
}
