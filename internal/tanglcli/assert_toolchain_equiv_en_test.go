//ff:func feature=cli type=command control=sequence
//ff:what assertToolchainEquiv — one fixture's gen/effects output must match across the direct-parse and .tangeul paths
package tanglcli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/en"
	"github.com/park-jun-woo/toulmin/pkg/tangl/effects"
	"github.com/park-jun-woo/toulmin/pkg/tangl/gen"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// assertToolchainEquiv parses the parser-testdata fixture named `name` directly
// (astA), then encodes+saves+loads it into astB via the codec and stream, and
// asserts gen.Generate and effects.Closure for the fixture's single endpoint agree
// across the two paths — the .tangeul path must reproduce the old text-direct
// output byte-for-byte.
func assertToolchainEquiv(t *testing.T, name, endpoint string) {
	t.Helper()
	p := filepath.Join("..", "..", "pkg", "tangl", "parser", "testdata", name)
	astA, err := parser.Parse(p)
	if err != nil {
		t.Fatalf("%s: direct parse: %v", name, err)
	}
	src, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("%s: read: %v", name, err)
	}
	d, err := en.ENCodec{}.Encode(src, p)
	if err != nil {
		t.Fatalf("%s: Encode: %v", name, err)
	}
	out := filepath.Join(t.TempDir(), "equiv.tangeul")
	if err := tangeul.Save(d, out); err != nil {
		t.Fatalf("%s: Save: %v", name, err)
	}
	loaded, err := tangeul.Load(out)
	if err != nil {
		t.Fatalf("%s: Load: %v", name, err)
	}
	astB := loaded.AST
	genA, err := gen.Generate(astA)
	if err != nil {
		t.Fatalf("%s: gen (direct): %v", name, err)
	}
	genB, err := gen.Generate(astB)
	if err != nil {
		t.Fatalf("%s: gen (.tangeul): %v", name, err)
	}
	if genA != genB {
		t.Errorf("%s: gen output differs between direct parse and .tangeul path", name)
	}
	eA, err := effects.Closure(astA, endpoint)
	if err != nil {
		t.Fatalf("%s: effects (direct): %v", name, err)
	}
	eB, err := effects.Closure(astB, endpoint)
	if err != nil {
		t.Fatalf("%s: effects (.tangeul): %v", name, err)
	}
	if formatEffects(eA) != formatEffects(eB) {
		t.Errorf("%s: effects for %q differ between direct parse and .tangeul path", name, endpoint)
	}
}
