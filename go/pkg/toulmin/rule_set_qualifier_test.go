//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRuleQualifier — tests Rule.Qualifier for below-range panic, above-range panic, and valid-set branches
package toulmin

import (
	"strings"
	"testing"
)

func TestRuleQualifier(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"PanicBelow", func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected panic for qualifier < 0.0")
				}
				msg, ok := r.(string)
				if !ok {
					t.Fatalf("expected string panic, got %T", r)
				}
				if !strings.Contains(msg, "qualifier must be between 0.0 and 1.0") {
					t.Errorf("unexpected panic message: %s", msg)
				}
			}()
			fn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			g := NewGraph("test")
			g.Rule(fn).Qualifier(-0.1)
		}},
		{"PanicAbove", func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected panic for qualifier > 1.0")
				}
				msg, ok := r.(string)
				if !ok {
					t.Fatalf("expected string panic, got %T", r)
				}
				if !strings.Contains(msg, "qualifier must be between 0.0 and 1.0") {
					t.Errorf("unexpected panic message: %s", msg)
				}
			}()
			g := NewGraph("test")
			g.Rule(WarrantA).Qualifier(1.5)
		}},
		{"Valid", func(t *testing.T) {
			g := NewGraph("test")
			r := g.Rule(WarrantA)
			got := r.Qualifier(0.5)
			if got != r {
				t.Errorf("Qualifier must return the receiver for chaining, got %v want %v", got, r)
			}
			if g.rules[r.idx].Qualifier != 0.5 {
				t.Errorf("expected Qualifier=0.5, got %f", g.rules[r.idx].Qualifier)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
