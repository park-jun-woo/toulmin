//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenCompileKo — a KO-authored Hangul-named document encodes, generates valid romanized Go, and go build's in an isolated module
package ko

import (
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/gen"
)

// TestGenCompileKo proves the full KO authoring → codegen → compile path
// (plans/tangl-ko §7 step 7, exercising the Phase 006 RR romanization): a
// self-contained Korean surface whose case, rule, definition and provides names are
// all Hangul is encoded by the KO codec, its AST fed to gen.Generate, and the result
// must be gofmt-clean, contain the romanized exported endpoints, and actually go
// build in a throwaway module that replaces this repo locally (its one local `do`
// target `logIt` is satisfied by a hand-written stub). Before the RR fix every
// Hangul name collapsed to "_" and the output could not compile.
func TestGenCompileKo(t *testing.T) {
	src := strings.Join([]string{
		"## tangl:주제", "- 이 문서는 `demo` 이다", "",
		"## tangl:정의", "- `기준값` 는 5 를 뜻한다", "",
		"## tangl:규칙", "- `점수` 가 5 이상이면 `자격` 이다", "",
		"## tangl:사안",
		"- `출금 가능` 의 경우에",
		"  - `금액` 는 필수다",
		"  - `자격` 은 일반 규칙이다",
		"  - `자격` 이면 `logIt` 실행한다", "",
		"- `입금 확인` 의 경우에",
		"  - `승인` 은 `출금 가능` 판정 쓰는 일반 규칙이다", "",
		"## tangl:업무",
		"- `출금` 를 한다면",
		"  - `금액` 는 필수다",
		"  - `출금 가능` 진행한다",
		"- `조회` 를 한다면",
		"  - `입금 확인` 판정한다", "",
	}, "\n")

	d, err := KoCodec{}.Encode([]byte(src), "hangul-ko.md")
	if err != nil {
		t.Fatalf("KO Encode: %v", err)
	}
	out, err := gen.Generate(d.AST)
	if err != nil {
		t.Fatalf("gen.Generate: %v", err)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("generated source is not gofmt-parseable: %v\n--- source ---\n%s", err, out)
	}
	if !strings.Contains(out, "func Chulgeum(") || !strings.Contains(out, "func Johoe(") {
		t.Fatalf("expected romanized exported endpoints Chulgeum and Johoe in output:\n%s", out)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "demo.go"), []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}
	stub := "package demo\n\nimport \"github.com/park-jun-woo/toulmin/pkg/toulmin\"\n\nfunc logIt(ctx toulmin.Context) error { return nil }\n"
	if err := os.WriteFile(filepath.Join(dir, "stub.go"), []byte(stub), 0o644); err != nil {
		t.Fatal(err)
	}
	repoRoot, err := filepath.Abs("../../../..")
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
}
