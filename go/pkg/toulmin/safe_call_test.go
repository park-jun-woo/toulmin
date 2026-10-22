//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestSafeCall — tests safeCall for normal return and panic-recovery branches
package toulmin

import (
	"strings"
	"testing"
)

func TestSafeCall(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Normal", func(t *testing.T) {
			fn := func(ctx Context, specs Specs) (bool, any) { return true, "evidence" }
			activated, evidence, err := safeCall(fn, NewContext(), nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !activated {
				t.Errorf("expected activated=true")
			}
			if evidence != "evidence" {
				t.Errorf("expected evidence=%q, got %v", "evidence", evidence)
			}
		}},
		{"Panic", func(t *testing.T) {
			fn := func(ctx Context, specs Specs) (bool, any) {
				panic("boom")
			}
			activated, evidence, err := safeCall(fn, NewContext(), nil)
			if err == nil {
				t.Fatal("expected error from recovered panic")
			}
			if !strings.Contains(err.Error(), "rule panicked") || !strings.Contains(err.Error(), "boom") {
				t.Errorf("unexpected error message: %v", err)
			}
			if activated {
				t.Errorf("expected activated=false on panic")
			}
			if evidence != nil {
				t.Errorf("expected evidence=nil on panic, got %v", evidence)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
