//ff:func feature=tangl type=model control=iteration dimension=1
//ff:what TestCollectSpecNames — every Expr kind contributes its SpecRef operands in order; others are ignored
package ast

import (
	"reflect"
	"testing"
)

func TestCollectSpecNames(t *testing.T) {
	spec := func(steps ...string) SpecRef { return SpecRef{Steps: steps} }
	cases := []struct {
		name string
		expr Expr
		want []string
	}{
		{"compare lhs then rhs", Compare{Lhs: spec("a"), Rhs: spec("b", "x")}, []string{"a", "b"}},
		{"compare nil operands", Compare{Field: "f", Op: "equals", Value: "1"}, nil},
		{"compare path operands ignored", Compare{Lhs: Path{Steps: []string{"p"}}, Rhs: Call{Fn: Ref{Name: "f"}}}, nil},
		{"logic terms in order", Logic{Op: "or", Terms: []Expr{Compare{Lhs: spec("b")}, Compare{Lhs: spec("a")}}}, []string{"b", "a"}},
		{"either terms", Either{Terms: []Expr{Compare{Rhs: spec("e")}}}, []string{"e"}},
		{"not term", Not{Term: Compare{Lhs: spec("n")}}, []string{"n"}},
		{"exists source before cond", Exists{Source: spec("src"), Binder: "x", Cond: Compare{Rhs: spec("lim")}}, []string{"src", "lim"}},
		{"exists path source", Exists{Source: Path{Steps: []string{"items"}}, Cond: Compare{Lhs: spec("k")}}, []string{"k"}},
		{"unhandled expr kind", Call{Fn: Ref{Name: "holds"}}, nil},
		{"nil expr", nil, nil},
		{"duplicates collapse", Logic{Terms: []Expr{Compare{Lhs: spec("d"), Rhs: spec("d")}, Not{Term: Compare{Lhs: spec("d")}}}}, []string{"d"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := collectSpecNames(c.expr, map[string]bool{}, nil)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("collectSpecNames = %v, want %v", got, c.want)
			}
		})
	}
	seen := map[string]bool{"pre": true}
	if got := collectSpecNames(Compare{Lhs: spec("pre"), Rhs: spec("new")}, seen, []string{"kept"}); !reflect.DeepEqual(got, []string{"kept", "new"}) || !seen["new"] {
		t.Fatalf("pre-seeded collect = %v seen=%v", got, seen)
	}
}
