//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestToRuleFunc — tests toRuleFunc for the verdict signature, the wrapped bool signature, and the default panic branch
package toulmin

import (
	"strings"
	"testing"
)

func TestToRuleFunc(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"VerdictSignature", func(t *testing.T) {
			fn := func(ctx Context, specs Specs) (float64, any) { return 0.5, "ev" }
			wrapped := toRuleFunc(fn)
			v, evidence := wrapped(NewContext(), nil)
			if v != 0.5 {
				t.Errorf("expected verdict=0.5, got %v", v)
			}
			if evidence != "ev" {
				t.Errorf("expected evidence=%q, got %v", "ev", evidence)
			}
		}},
		{"BoolSignature", func(t *testing.T) {
			fn := func(ctx Context, specs Specs) (bool, any) { return true, "ev" }
			wrapped := toRuleFunc(fn)
			v, evidence := wrapped(NewContext(), nil)
			if v != 1.0 {
				t.Errorf("expected true → +1, got %v", v)
			}
			if evidence != "ev" {
				t.Errorf("expected evidence=%q, got %v", "ev", evidence)
			}
		}},
		{"PanicsOnWrongType", func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected panic for wrong function type")
				}
				msg, ok := r.(string)
				if !ok {
					t.Fatalf("expected string panic, got %T", r)
				}
				if !strings.Contains(msg, "toulmin:") || !strings.Contains(msg, "(float64, any)") {
					t.Errorf("unexpected panic message: %s", msg)
				}
			}()
			toRuleFunc(func() {})
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
