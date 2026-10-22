//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestToRuleFunc — tests toRuleFunc for the matching-type success branch and the default panic branch
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
		{"Success", func(t *testing.T) {
			fn := func(ctx Context, specs Specs) (bool, any) { return true, "ev" }
			wrapped := toRuleFunc(fn)
			activated, evidence := wrapped(NewContext(), nil)
			if !activated {
				t.Errorf("expected activated=true")
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
				if !strings.Contains(msg, "toulmin:") {
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
