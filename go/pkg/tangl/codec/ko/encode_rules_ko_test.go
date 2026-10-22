//ff:func feature=tangl type=parser control=sequence
//ff:what TestEncodeRulesKo — a 규칙 doc (nested condition tree) tiles and matches its EN counterpart's AST
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

// TestEncodeRulesKo exercises the tangl:규칙 section (absent from the three example
// fixtures): a simple ~이상이면 rule and a 다음이면 rule whose nested 다음 중/그리고/또는
// tree becomes one absorbed StmtBlock. It checks byte-faithful tiling and that the
// KO-encoded AST matches the EN counterpart's parse (Line ignored, Path aligned).
func TestEncodeRulesKo(t *testing.T) {
	koPath := filepath.Join("testdata", "rules-ko.md")
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
	if len(d.AST.Rules) != 2 {
		t.Fatalf("Rules = %d, want 2", len(d.AST.Rules))
	}
	enSrc, err := os.ReadFile(filepath.Join("testdata", "rules-en.md"))
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
