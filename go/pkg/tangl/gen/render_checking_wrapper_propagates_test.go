//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestRenderCheckingWrapperPropagates — a failing checked case surfaces as the caller's Evaluate error; an explicit force compiles
package gen

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// TestRenderCheckingWrapperPropagates generates a document whose checked case
// `can inner` has a host rule that panics, builds it in a throwaway module and
// runs a test there: Evaluating `can review` (which checks `can inner`) must
// return an error carrying the inner failure — never a silent false (−1). The
// same document carries an explicit-force attack, so the build also proves
// `Attacks(x, 0.5)` compiles against the engine.
func TestRenderCheckingWrapperPropagates(t *testing.T) {
	src := "" +
		"## tangl:Subject\n- this document is `demo`\n\n" +
		"## tangl:Rules\n- `explode` uses `explode`\n- `ok` when `score` is at least 5\n\n" +
		"## tangl:Cases\n" +
		"- in case of `can inner`\n" +
		"  - `explode` is a general rule\n" +
		"  - `ok` is a counter rule\n" +
		"  - `ok` attacks `explode` with force 0.5\n\n" +
		"- in case of `can review`\n" +
		"  - `approved` is a general rule checking `can inner`\n"
	doc, err := parser.ParseSource(src, "propagate.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	out, err := Generate(doc)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	dir := t.TempDir()
	files := map[string]string{
		"demo.go": out,
		"stub.go": "package demo\n\nimport \"github.com/park-jun-woo/toulmin/go/pkg/toulmin\"\n\n" +
			"func explode(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { panic(\"boom\") }\n",
		"demo_test.go": "package demo\n\nimport (\n\t\"strings\"\n\t\"testing\"\n\n\t\"github.com/park-jun-woo/toulmin/go/pkg/toulmin\"\n)\n\n" +
			"func TestPropagates(t *testing.T) {\n" +
			"\t_, err := canReviewGraph.Evaluate(toulmin.NewContext())\n" +
			"\tif err == nil || !strings.Contains(err.Error(), \"boom\") {\n" +
			"\t\tt.Fatalf(\"err = %v, want the inner failure\", err)\n\t}\n}\n",
	}
	repoRoot, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	files["go.mod"] = "module demo\n\ngo 1.18\n\nrequire github.com/park-jun-woo/toulmin/go v0.0.0\n\nreplace github.com/park-jun-woo/toulmin/go => " + repoRoot + "\n"
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOPROXY=off", "GOWORK=off")
	if cmdOut, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated test failed: %v\n%s\n--- generated source ---\n%s", err, cmdOut, out)
	}
}
