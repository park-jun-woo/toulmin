//ff:func feature=engine type=util control=iteration dimension=1
//ff:what TestFuncID — tests funcID for resolved-name and FuncForPC-nil-fallback branches
package toulmin

import (
	"strings"
	"testing"
)

func TestFuncID(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"ResolvesName", func(t *testing.T) {
			fn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			id := funcID(fn)
			if id == "" {
				t.Fatal("expected non-empty funcID")
			}
			if !strings.Contains(id, ".") {
				t.Errorf("expected dotted full path name, got %q", id)
			}
		}},
		{"FallbackOnNilFunc", func(t *testing.T) {
			// A nil func value has pointer 0, for which runtime.FuncForPC returns
			// nil, exercising funcID's "unknown_%d" fallback branch.
			var fn func(ctx Context, specs Specs) (bool, any)
			id := funcID(fn)
			if !strings.HasPrefix(id, "unknown_") {
				t.Errorf("expected unknown_ fallback prefix, got %q", id)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
