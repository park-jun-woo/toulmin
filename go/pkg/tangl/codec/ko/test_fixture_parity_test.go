//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestFixtureParity — the KO example docs encode, tile, and match their EN counterparts' AST
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

// TestFixtureParity encodes each shared example document (아메리카노/접근 제어/계좌 이체
// and the quest_gate supersede+level fixture) with KoCodec and checks: the Blocks pass VerifyTiling,
// Doc.Locale is "ko", RenderVerbatim reproduces the source byte-for-byte (mode 1),
// and the encoded AST equals its EN counterpart's parse (EqualIgnoringLine, with
// Path aligned since each source stamps its own path).
func TestFixtureParity(t *testing.T) {
	for _, name := range []string{"americano", "access", "transfer", "quest_gate"} {
		koPath := filepath.Join("testdata", name+"-ko.md")
		enPath := filepath.Join("testdata", name+"-en.md")
		src, err := os.ReadFile(koPath)
		if err != nil {
			t.Fatalf("read %s: %v", koPath, err)
		}
		d, err := KoCodec{}.Encode(src, koPath)
		if err != nil {
			t.Fatalf("%s: Encode: %v", name, err)
		}
		if d.Locale != "ko" {
			t.Errorf("%s: Locale = %q, want ko", name, d.Locale)
		}
		if err := doc.VerifyTiling(d); err != nil {
			t.Errorf("%s: VerifyTiling: %v", name, err)
		}
		if !bytes.Equal(doc.RenderVerbatim(d), src) {
			t.Errorf("%s: RenderVerbatim is not byte-identical to source", name)
		}
		enSrc, err := os.ReadFile(enPath)
		if err != nil {
			t.Fatalf("read %s: %v", enPath, err)
		}
		enAST, err := parser.ParseSource(string(enSrc), koPath)
		if err != nil {
			t.Fatalf("%s: parse EN counterpart: %v", name, err)
		}
		if !tangeul.EqualIgnoringLine(d.AST, enAST) {
			t.Errorf("%s: KO-encoded AST differs from EN parse", name)
		}
	}
}
