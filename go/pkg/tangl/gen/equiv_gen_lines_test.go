//ff:func feature=tangl type=codegen control=sequence
//ff:what equivGenLines — compiles and runs the generated cross-test code over the ctx table
package gen

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// equivGenLines generates Go code for the cross-test document, rewrites its
// package clause to main, pairs it with the runner source (host stubs + rule
// map + ctx table), go-runs the pair in a throwaway module that replaces
// this repo locally, and returns the printed "<ctx>\t<rule>\t<canon>" lines —
// the gen side of the eval≡gen equivalence check.
func equivGenLines(t *testing.T, doc *ast.Document, ctxJSON string) []string {
	t.Helper()
	src, err := Generate(doc)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !strings.Contains(src, "package equiv") {
		t.Fatalf("expected generated package equiv, got:\n%s", src)
	}
	src = strings.Replace(src, "package equiv", "package main", 1)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "gen.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	runner := equivRunnerSrc(doc.Rules, ctxJSON)
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(runner), 0o644); err != nil {
		t.Fatal(err)
	}
	repoRoot, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	gomod := "module equivrun\n\ngo 1.18\n\nrequire github.com/park-jun-woo/toulmin/go v0.0.0\n\nreplace github.com/park-jun-woo/toulmin/go => " + repoRoot + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOPROXY=off")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\n%s\n--- generated source ---\n%s", err, out, src)
	}
	return strings.Split(strings.TrimRight(string(out), "\n"), "\n")
}
