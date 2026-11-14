//ff:func feature=tangl type=util control=sequence
//ff:what checkEncoderXCheck — one fixture: the tangl-en encoder's stream agrees structurally with the step-1 golden
package tangeul

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/en"
)

// checkEncoderXCheck runs the Phase 008 step-2 cross-check for one fixture. It builds
// the stream two independent ways — the tangl-en encoder (en.ENCodec.Encode) and the
// frozen step-1 golden (reloaded from testdata via Load) — and asserts they agree on
// the two axes that define stream-level equivalence here: (1) the recovered AST up to
// Line (EqualIgnoringLine, with Document.Path aligned since each path stamps its own),
// and (2) the ordered signatures of the non-prose blocks (stmtBlockSignatures) — same
// Section/Stmt tiles, same StmtRef coordinates, same order. Byte-level prose spans are
// deliberately NOT compared: the golden uses 4-byte placeholder spans (a pre-encoder
// stand-in, per assembleDoc) while the encoder emits the real authored prose, so the
// two streams differ in bytes by design. Mode-1 byte identity on the encoder's own
// output is covered separately by the CLI e2e (TestEncodeDecodeEN in tanglcli).
func checkEncoderXCheck(t *testing.T, name string) {
	t.Helper()
	md := filepath.Join("..", "tangl", "parser", "testdata", name+".md")
	src, err := os.ReadFile(md)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	enc, err := en.ENCodec{}.Encode(src, md)
	if err != nil {
		t.Fatalf("encode %s: %v", name, err)
	}
	golden, err := Load(filepath.Join("testdata", name+".geul"))
	if err != nil {
		t.Fatalf("load golden %s: %v", name, err)
	}
	enc.AST.Path = golden.AST.Path
	if !EqualIgnoringLine(enc.AST, golden.AST) {
		t.Errorf("%s: encoder AST differs from step-1 golden AST (ignoring Line)", name)
	}
	encSigs := stmtBlockSignatures(enc.Blocks)
	goldenSigs := stmtBlockSignatures(golden.Blocks)
	if !reflect.DeepEqual(encSigs, goldenSigs) {
		t.Errorf("%s: non-prose block sequence differs:\n encoder %v\n  golden %v", name, encSigs, goldenSigs)
	}
}
