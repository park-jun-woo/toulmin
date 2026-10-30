//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseSupersedeChainKoErrors — non-chain text is unmatched; missing backticks, duplicates, bad separators and short chains fail
package ko

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseSupersedeChainKoErrors(t *testing.T) {
	if edges, ok, err := parseSupersedeChainKo(scan.Item{Text: "`A` 면 `B` 배제한다", No: 1}, "t.md"); ok || err != nil || edges != nil {
		t.Fatalf("pair edge: edges=%v ok=%v err=%v, want unmatched", edges, ok, err)
	}
	for _, c := range []struct{ text, want string }{
		{"A > B > C 순서로 배제한다", "expected backtick-quoted name"},
		{"`A` > `B` > `A` 순서로 배제한다", "duplicate name"},
		{"`A` < `B` > `C` 순서로 배제한다", "expected '>'"},
		{"`A` > `B` 순서로 배제한다", "at least 3 names"},
	} {
		_, ok, err := parseSupersedeChainKo(scan.Item{Text: c.text, No: 2}, "t.md")
		if !ok || err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: ok=%v err=%v, want %q", c.text, ok, err, c.want)
		}
	}
	edges, ok, err := parseSupersedeChainKo(scan.Item{Text: "`A` > `B` > `C` > `D` 순서로 배제한다", No: 3}, "t.md")
	if err != nil || !ok || len(edges) != 6 || len(edges[0].Chain) != 4 || edges[1].Chain != nil {
		t.Fatalf("four-name chain: edges=%+v ok=%v err=%v", edges, ok, err)
	}
}
