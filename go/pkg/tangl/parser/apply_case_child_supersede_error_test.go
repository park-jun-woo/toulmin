//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplyCaseChildSupersedeErrors — supersede chains and edges append on success and propagate parse errors
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestApplyCaseChildSupersedeErrors(t *testing.T) {
	c := &ast.Case{}
	if err := applyCaseChild(c, item{Text: "supersede in order: `A` > `B`", Line: 3}, "t.md"); err == nil {
		t.Fatal("a two-name chain must be rejected")
	}
	if err := applyCaseChild(c, item{Text: "`A` supersedes nothing", Line: 4}, "t.md"); err == nil {
		t.Fatal("a supersede edge without a backtick target must be rejected")
	}
	if len(c.Supersedes) != 0 {
		t.Fatalf("failed statements must not add edges: %+v", c.Supersedes)
	}
	if err := applyCaseChild(c, item{Text: "supersede in order: `A` > `B` > `C`", Line: 5}, "t.md"); err != nil || len(c.Supersedes) != 3 {
		t.Fatalf("valid chain: err=%v edges=%+v", err, c.Supersedes)
	}
	if err := applyCaseChild(c, item{Text: "`C` supersedes `A`", Line: 6}, "t.md"); err != nil || len(c.Supersedes) != 4 {
		t.Fatalf("valid edge: err=%v edges=%+v", err, c.Supersedes)
	}
	if last := c.Supersedes[3]; last.Superseder != "C" || last.Superseded != "A" || last.Line != 6 || last.Chain != nil {
		t.Fatalf("appended edge = %+v", last)
	}
}
