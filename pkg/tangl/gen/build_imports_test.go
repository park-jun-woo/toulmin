//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildImports — tests buildImports for all-flags-off and all-flags-on branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildImports(t *testing.T) {
	t.Run("all flags off", func(t *testing.T) {
		gc := &genContext{Doc: &ast.Document{}, Flags: genFlags{}}
		specs, err := buildImports(gc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(specs) != 1 || specs[0].Alias != "toulmin" {
			t.Errorf("expected only toulmin import, got %+v", specs)
		}
	})

	t.Run("all flags on", func(t *testing.T) {
		gc := &genContext{
			Doc: &ast.Document{},
			Flags: genFlags{
				NeedsTangl:   true,
				NeedsTime:    true,
				NeedsCompare: true,
			},
		}
		specs, err := buildImports(gc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// toulmin, tangl, time, fmt, strconv, strings = 6 entries (no alias imports)
		if len(specs) != 6 {
			t.Fatalf("expected 6 import specs, got %+v", specs)
		}
		wantPaths := []string{
			"github.com/park-jun-woo/toulmin/pkg/toulmin",
			"github.com/park-jun-woo/toulmin/pkg/tangl",
			"time",
			"fmt",
			"strconv",
			"strings",
		}
		for i, p := range wantPaths {
			if specs[i].Path != p {
				t.Errorf("specs[%d].Path = %q, want %q", i, specs[i].Path, p)
			}
		}
	})
}
