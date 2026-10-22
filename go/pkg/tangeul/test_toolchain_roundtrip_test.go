//ff:func feature=tangl type=util control=sequence
//ff:what TestToolchainRoundtrip — Doc→Save→Load AST drives validate/effects/gen identically to the direct AST
package tangeul

import (
	"go/format"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/effects"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/gen"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/validate"
)

// TestToolchainRoundtrip confirms the retargeting contract: an AST recovered via
// Save→Load feeds the unmodified toolchain (validate.Validate, validate.Lint,
// effects.Closure, gen.Generate) with results identical to the direct in-memory
// AST. validate passes and lint is empty, the effect closure matches the §10
// expected do/undo entries, and gen output equals the direct-AST output and is
// gofmt-clean (Path/Line are unused by gen, so the .geul path stamp is inert).
func TestToolchainRoundtrip(t *testing.T) {
	direct := buildTransferDoc()

	path := filepath.Join(t.TempDir(), "transfer.geul")
	if err := Save(direct, path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if err := validate.Validate(loaded.AST); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if warns := validate.Lint(loaded.AST); len(warns) != 0 {
		t.Fatalf("Lint = %v, want none", warns)
	}

	cl, err := effects.Closure(loaded.AST, "transfer")
	if err != nil {
		t.Fatalf("Closure: %v", err)
	}
	want := []effects.Entry{
		{Kind: "do", Func: ast.Ref{Alias: "bank", Name: "withdraw"}, Once: true, Case: "can withdraw", Node: "balance sufficient"},
		{Kind: "undo", Func: ast.Ref{Alias: "bank", Name: "refund"}, Once: false, Case: "can withdraw", Node: "balance sufficient"},
		{Kind: "do", Func: ast.Ref{Alias: "bank", Name: "deposit"}, Once: true, Case: "can deposit", Node: "recipient valid"},
	}
	if !reflect.DeepEqual(cl, want) {
		t.Fatalf("Closure = %+v, want %+v", cl, want)
	}

	loadedCode, err := gen.Generate(loaded.AST)
	if err != nil {
		t.Fatalf("Generate(loaded): %v", err)
	}
	directCode, err := gen.Generate(direct.AST)
	if err != nil {
		t.Fatalf("Generate(direct): %v", err)
	}
	if loadedCode != directCode {
		t.Fatalf("gen output differs between loaded and direct AST")
	}
	formatted, err := format.Source([]byte(loadedCode))
	if err != nil {
		t.Fatalf("gofmt of gen output: %v", err)
	}
	if string(formatted) != loadedCode {
		t.Fatal("gen output is not gofmt-clean")
	}
}
