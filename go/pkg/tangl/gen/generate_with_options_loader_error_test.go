//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerateWithOptionsLoaderError — a file definition without a loader registry aborts generation
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestGenerateWithOptionsLoaderError(t *testing.T) {
	doc := &ast.Document{Path: "/p.md", Subject: "list", Defs: []ast.Definition{{Name: "domains", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "x.json", Line: 3}}}
	out, err := GenerateWithOptions(doc, GenerateOptions{})
	if err == nil || !strings.Contains(err.Error(), `/p.md:3: definition "domains": loader registry is required`) {
		t.Fatalf("err = %v, want the missing-registry error", err)
	}
	if out != "" {
		t.Fatalf("out = %q, want empty on error", out)
	}
}
