//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildCheckingWrappers — tests buildCheckingWrappers for empty and non-empty checking-target sets
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildCheckingWrappers(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		doc := &ast.Document{}
		var w strings.Builder
		wrappers := buildCheckingWrappers(&w, doc)
		if len(wrappers) != 0 {
			t.Errorf("expected no wrappers, got %+v", wrappers)
		}
		if w.String() != "" {
			t.Errorf("expected no output, got %q", w.String())
		}
	})

	t.Run("non-empty sorted", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Nodes: []ast.Node{
						{Name: "n1", Checking: "zebra case"},
						{Name: "n2", Checking: "alpha case"},
						{Name: "n3"}, // no checking clause: not a target
					},
				},
			},
		}
		var w strings.Builder
		wrappers := buildCheckingWrappers(&w, doc)
		if len(wrappers) != 2 {
			t.Fatalf("expected 2 wrappers, got %+v", wrappers)
		}
		zebraFn, ok := wrappers["zebra case"]
		if !ok {
			t.Fatal("expected wrapper for 'zebra case'")
		}
		alphaFn, ok := wrappers["alpha case"]
		if !ok {
			t.Fatal("expected wrapper for 'alpha case'")
		}
		out := w.String()
		if !strings.Contains(out, "func "+alphaFn) || !strings.Contains(out, "func "+zebraFn) {
			t.Errorf("expected both wrapper functions rendered, got:\n%s", out)
		}
		// sorted order: "alpha case" < "zebra case", so alpha's func appears first
		if strings.Index(out, alphaFn) > strings.Index(out, zebraFn) {
			t.Errorf("expected alpha wrapper rendered before zebra wrapper, got:\n%s", out)
		}
	})
}
