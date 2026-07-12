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
				NeedsTangl:      true,
				NeedsTime:       true,
				NeedsCompareVal: true,
				NeedsWalk:       true,
				NeedsItoa:       true,
			},
		}
		specs, err := buildImports(gc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		wantPaths := []string{
			"github.com/park-jun-woo/toulmin/pkg/toulmin",
			"github.com/park-jun-woo/toulmin/pkg/tangl",
			"time",
			"fmt",
			"reflect",
			"strconv",
			"strings",
		}
		if len(specs) != len(wantPaths) {
			t.Fatalf("expected %d import specs, got %+v", len(wantPaths), specs)
		}
		for i, p := range wantPaths {
			if specs[i].Path != p {
				t.Errorf("specs[%d].Path = %q, want %q", i, specs[i].Path, p)
			}
		}
	})

	t.Run("itoa alone gates strconv only", func(t *testing.T) {
		gc := &genContext{Doc: &ast.Document{}, Flags: genFlags{NeedsItoa: true}}
		specs, err := buildImports(gc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(specs) != 2 || specs[1].Path != "strconv" {
			t.Errorf("expected toulmin+strconv, got %+v", specs)
		}
	})
}
