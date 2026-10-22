//ff:func feature=tangl type=parser control=sequence
//ff:what TestEncodeRulesExprKo — the new-syntax 규칙 fixture tiles, renders verbatim, and matches its EN AST
package ko

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// TestEncodeRulesExprKo exercises the판정식 extensions (field paths, spec refs, host
// value/predicate calls, and one-line/nested exists) in the tangl:규칙 section: the KO
// fixture must pass VerifyTiling, reproduce its source byte-for-byte (mode 1), carry
// all five rules, and encode to the same AST as its EN counterpart's direct parse
// (EqualIgnoringLine, Path aligned).
func TestEncodeRulesExprKo(t *testing.T) {
	koPath := filepath.Join("testdata", "rules_expr-ko.md")
	src, err := os.ReadFile(koPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	d, err := KoCodec{}.Encode(src, koPath)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if err := doc.VerifyTiling(d); err != nil {
		t.Fatalf("VerifyTiling: %v", err)
	}
	if !bytes.Equal(doc.RenderVerbatim(d), src) {
		t.Error("RenderVerbatim is not byte-identical to source")
	}
	if len(d.AST.Rules) != 5 {
		t.Fatalf("Rules = %d, want 5", len(d.AST.Rules))
	}
	enSrc, err := os.ReadFile(filepath.Join("testdata", "rules_expr-en.md"))
	if err != nil {
		t.Fatalf("read EN: %v", err)
	}
	enAST, err := parser.ParseSource(string(enSrc), koPath)
	if err != nil {
		t.Fatalf("parse EN counterpart: %v", err)
	}
	if !tangeul.EqualIgnoringLine(d.AST, enAST) {
		t.Error("KO-encoded AST differs from EN parse")
	}
}
