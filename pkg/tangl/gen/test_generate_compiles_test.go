//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCompileSyntheticDoc — Generate output actually go build's, not just gofmt's
package gen

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestCompileSyntheticDoc exercises the code paths the three spec
// fixtures don't (tangl:Rules, "checking", a "check" endpoint) with a
// synthetic document that has no external `see` packages, so the
// generated file — plus one hand-written stub for its one local `do`
// function — can be fully go build'd in a throwaway module that replaces
// this repo's module locally (no network access required).
func TestCompileSyntheticDoc(t *testing.T) {
	src := "" +
		"## tangl:Subject\n- this document is `demo`\n\n" +
		"## tangl:Definitions\n- `threshold` means 5\n\n" +
		"## tangl:Rules\n- `qualifies` when `score` is at least 5\n\n" +
		"## tangl:Cases\n" +
		"- in case of `can approve`\n" +
		"  - `amount` is required\n" +
		"  - `qualifies` is a general rule\n" +
		"  - do `logIt` when `qualifies`\n\n" +
		"- in case of `can review`\n" +
		"  - `approved` is a general rule checking `can approve`\n\n" +
		"## tangl:Provides\n" +
		"- provides `approve`\n" +
		"  - `amount` is required\n" +
		"  - run `can approve`\n" +
		"- provides `review`\n" +
		"  - check `can review`\n"
	doc, err := parser.ParseSource(src, "synthetic.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	out, err := Generate(doc)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "demo.go"), []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}
	stub := "package demo\n\nimport \"github.com/park-jun-woo/toulmin/pkg/toulmin\"\n\nfunc logIt(ctx toulmin.Context) error { return nil }\n"
	if err := os.WriteFile(filepath.Join(dir, "stub.go"), []byte(stub), 0o644); err != nil {
		t.Fatal(err)
	}
	repoRoot, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	gomod := "module demo\n\ngo 1.18\n\nrequire github.com/park-jun-woo/toulmin v0.0.0\n\nreplace github.com/park-jun-woo/toulmin => " + repoRoot + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOPROXY=off")
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\n%s\n--- generated source ---\n%s", err, cmdOut, out)
	}
	t.Logf("go build ok:\n%s", cmdOut)
}
