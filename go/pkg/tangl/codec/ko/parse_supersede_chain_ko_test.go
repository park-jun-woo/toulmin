//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseSupersedeChainKo — Korean chain expands identically to canonical pair edges
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseSupersedeChainKo(t *testing.T) {
	edges, ok, err := parseSupersedeChainKo(scan.Item{Text: "`A` > `B` > `C` 순서로 배제한다", No: 9}, "t.md")
	if err != nil || !ok || len(edges) != 3 {
		t.Fatalf("ok=%v err=%v edges=%+v", ok, err, edges)
	}
	want := [][2]string{{"A", "B"}, {"A", "C"}, {"B", "C"}}
	for i, pair := range want {
		if edges[i].Superseder != pair[0] || edges[i].Superseded != pair[1] {
			t.Fatalf("edge[%d]=%+v want=%v", i, edges[i], pair)
		}
	}
}
