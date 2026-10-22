//ff:func feature=tangl type=util control=sequence
//ff:what TestRulesExprStreamRoundTrip — the new judgment-expression fixture survives a full .geul round-trip (mode-1 bytes + AST up to Line)
package tangeul

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/en"
)

// TestRulesExprStreamRoundTrip closes the Phase003 end-to-end loop that Phase005
// deferred: the rules_expr.md fixture exercises every new judgment-expression
// node (Exists, predicate Call, operand-carrying CompareX, and Path/SpecRef/Call
// operands). It encodes the source, serializes it through the unmodified stream
// Writer (Save → stream.WriteDoc), reloads it (Load → stream.ReadDoc), and
// asserts both mode-1 byte identity (RenderVerbatim == source) and AST identity
// up to Line (the stream never stores authored line numbers). This proves the
// new EXPR encoding round-trips through the whole stack with stream/ untouched.
func TestRulesExprStreamRoundTrip(t *testing.T) {
	md := filepath.Join("..", "tangl", "parser", "testdata", "rules_expr.md")
	src, err := os.ReadFile(md)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	enc, err := en.ENCodec{}.Encode(src, md)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	out := filepath.Join(t.TempDir(), "rules_expr.geul")
	if err := Save(enc, out); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load(out)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !bytes.Equal(doc.RenderVerbatim(loaded), src) {
		t.Error("mode-1 RenderVerbatim after stream round-trip is not byte-identical to source")
	}
	loaded.AST.Path = enc.AST.Path
	if !EqualIgnoringLine(enc.AST, loaded.AST) {
		t.Error("reloaded AST differs from the encoded AST (ignoring Line)")
	}
}
