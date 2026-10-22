//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerateCrossDoc — two documents where one alias-checks the other's case actually go build together
package gen

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// TestGenerateCrossDoc is the strongest verification of the multi-document
// loader: gate.md defines a case that gate itself never checks, and turn.md
// alias-checks it via `checking `gate`.`can pass“. Each document is generated
// into its own Go package directory inside a throwaway module, and the two
// packages must compile together — proving gen emits an exported
// `gate.CheckingCanPass` wrapper in the gate package and a matching
// cross-package call plus derived import in the turn package.
func TestGenerateCrossDoc(t *testing.T) {
	repoRoot, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	mod := t.TempDir()
	gomod := "module reinsdemo\n\ngo 1.18\n\nrequire github.com/park-jun-woo/toulmin/go v0.0.0\n\nreplace github.com/park-jun-woo/toulmin/go => " + repoRoot + "\n"
	if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}

	gateDir := filepath.Join(mod, "gate")
	turnDir := filepath.Join(mod, "turn")
	if err := os.MkdirAll(gateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(turnDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// gate defines `can pass` but never checks it itself — it is referenced
	// only from turn, exercising the "every case is externally addressable"
	// wrapper expansion.
	gateSrc := "" +
		"## tangl:Subject\n- this document is `gate`\n\n" +
		"## tangl:Rules\n- `qualifies` when `score` is at least 5\n\n" +
		"## tangl:Cases\n" +
		"- in case of `can pass`\n" +
		"  - `qualifies` is a general rule\n"

	// turn alias-checks gate's case through a See whose package path points at
	// the sibling gate document directory.
	turnSrc := "" +
		"## tangl:Subject\n- this document is `turn`\n\n" +
		"## tangl:See\n- see `gate` from `../gate/gate`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `submit`\n" +
		"  - `gate ok` is a general rule checking `gate`.`can pass`\n"

	gatePath := filepath.Join(gateDir, "gate.md")
	turnPath := filepath.Join(turnDir, "turn.md")

	gateDoc, err := parser.ParseSource(gateSrc, gatePath)
	if err != nil {
		t.Fatalf("parse gate: %v", err)
	}
	turnDoc, err := parser.ParseSource(turnSrc, turnPath)
	if err != nil {
		t.Fatalf("parse turn: %v", err)
	}

	gateOut, err := Generate(gateDoc)
	if err != nil {
		t.Fatalf("generate gate: %v", err)
	}
	turnOut, err := Generate(turnDoc)
	if err != nil {
		t.Fatalf("generate turn: %v", err)
	}

	// gate exports the wrapper the cross-package call needs.
	if !strings.Contains(gateOut, "func CheckingCanPass(ctx toulmin.Context, specs toulmin.Specs) (bool, any)") {
		t.Errorf("gate output missing exported CheckingCanPass wrapper:\n%s", gateOut)
	}
	// turn calls it cross-package and imports the derived gate package path.
	if !strings.Contains(turnOut, "gate.CheckingCanPass") {
		t.Errorf("turn output missing cross-package call gate.CheckingCanPass:\n%s", turnOut)
	}
	if !strings.Contains(turnOut, `"reinsdemo/gate"`) {
		t.Errorf("turn output missing derived import of the gate package:\n%s", turnOut)
	}

	if err := os.WriteFile(filepath.Join(gateDir, "gate.go"), []byte(gateOut), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(turnDir, "turn.go"), []byte(turnOut), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = mod
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOPROXY=off")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cross-package go build failed: %v\n%s\n--- gate ---\n%s\n--- turn ---\n%s", err, out, gateOut, turnOut)
	}
}
