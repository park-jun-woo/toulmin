//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplyNodeSource — using / checking / neither source clauses and their error branches
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestApplyNodeSource(t *testing.T) {
	n := &ast.Node{}
	rest, err := applyNodeSource(n, "using `p`.`f` with `t` qualified 0.25 grounded on `g`, level fail", 1, "t.md")
	if err != nil || n.Using == nil || n.Using.Alias != "p" || len(n.With) != 1 || n.Qualified == nil || *n.Qualified != 0.25 || len(n.Grounds) != 1 {
		t.Fatalf("using: got %+v err %v", n, err)
	}
	if rest != ", level fail" {
		t.Errorf("using rest = %q", rest)
	}
	n = &ast.Node{}
	rest, err = applyNodeSource(n, "checking `loan`.`can renew` qualified 0.5", 1, "t.md")
	if err != nil || n.Checking == nil || n.Checking.Alias != "loan" || n.Checking.Name != "can renew" || rest != " qualified 0.5" {
		t.Fatalf("checking: got %+v rest %q err %v", n, rest, err)
	}
	n = &ast.Node{}
	rest, err = applyNodeSource(n, "qualified 0.5", 1, "t.md")
	if err != nil || rest != "qualified 0.5" || n.Using != nil || n.Checking != nil {
		t.Fatalf("neither: got %+v rest %q err %v", n, rest, err)
	}
	if _, err := applyNodeSource(&ast.Node{}, "using nope", 1, "t.md"); err == nil {
		t.Error("using without a reference must fail")
	}
	if _, err := applyNodeSource(&ast.Node{}, "checking nope", 1, "t.md"); err == nil {
		t.Error("checking without a case reference must fail")
	}
}
