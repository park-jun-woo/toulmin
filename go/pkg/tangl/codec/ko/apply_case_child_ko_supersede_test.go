//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplyCaseChildKoSupersede — KO chains and edges append to Supersedes (slot 4) or propagate parse errors
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestApplyCaseChildKoSupersede(t *testing.T) {
	c := &ast.Case{}
	if _, err := applyCaseChildKo(c, scan.Item{Text: "`A` > `B` 순서로 배제한다", No: 2}, "t.md"); err == nil {
		t.Fatal("a two-name chain must be rejected")
	}
	if _, err := applyCaseChildKo(c, scan.Item{Text: "`A` 면 배제한다", No: 3}, "t.md"); err == nil {
		t.Fatal("an edge without a superseded target must be rejected")
	}
	if len(c.Supersedes) != 0 {
		t.Fatalf("failed statements must not add edges: %+v", c.Supersedes)
	}
	slot, err := applyCaseChildKo(c, scan.Item{Text: "`A` > `B` > `C` 순서로 배제한다", No: 4}, "t.md")
	if err != nil || slot != 4 || len(c.Supersedes) != 3 {
		t.Fatalf("chain: slot=%d err=%v edges=%+v", slot, err, c.Supersedes)
	}
	slot, err = applyCaseChildKo(c, scan.Item{Text: "`C` 면 `A` 배제한다", No: 5}, "t.md")
	if err != nil || slot != 4 || len(c.Supersedes) != 4 {
		t.Fatalf("edge: slot=%d err=%v edges=%+v", slot, err, c.Supersedes)
	}
	if last := c.Supersedes[3]; last.Superseder != "C" || last.Superseded != "A" || last.Line != 5 {
		t.Fatalf("appended edge = %+v", last)
	}
}
