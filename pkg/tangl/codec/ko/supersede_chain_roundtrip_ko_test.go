//ff:what TestSupersedeChainRoundTripKo — Korean codec preserves chain surface while AST contains the closure
package ko

import (
	"bytes"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

func TestSupersedeChainRoundTripKo(t *testing.T) {
	src := []byte("## tangl:주제\n- 이 문서는 `x` 이다\n\n## tangl:사안\n- `c` 의 경우에\n  - `A` 는 일반 규칙이다\n  - `B` 는 반박 규칙이다\n  - `C` 는 예외 규칙이다\n  - `A` > `B` > `C` 순서로 배제한다\n")
	d, err := (KoCodec{}).Encode(src, "t.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(d.AST.Cases[0].Supersedes) != 3 {
		t.Fatalf("edges=%+v", d.AST.Cases[0].Supersedes)
	}
	rendered, err := doc.Render(d, KoCodec{})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rendered, src) {
		t.Fatalf("rendered:\n%s", rendered)
	}
}
