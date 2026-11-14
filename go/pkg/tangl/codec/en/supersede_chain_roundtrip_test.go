//ff:func feature=tangl type=parser control=sequence
//ff:what TestSupersedeChainRoundTrip — English codec preserves chain surface while AST contains the closure
package en

import (
	"bytes"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

func TestSupersedeChainRoundTrip(t *testing.T) {
	src := []byte("## tangl:Subject\n- this document is `x`\n\n## tangl:Cases\n- in case of `c`\n  - `A` is a general rule\n  - `B` is a counter rule\n  - `C` is an except rule\n  - supersede in order: `A` > `B` > `C`\n")
	d, err := (ENCodec{}).Encode(src, "t.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(d.AST.Cases[0].Supersedes) != 3 {
		t.Fatalf("edges=%+v", d.AST.Cases[0].Supersedes)
	}
	rendered, err := doc.Render(d, ENCodec{})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rendered, src) {
		t.Fatalf("rendered:\n%s", rendered)
	}
}
