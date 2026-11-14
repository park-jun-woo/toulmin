//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateWeights — Validate reports a qualifier outside [0, 1] carried by an AST that bypassed the parser
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestValidateWeights(t *testing.T) {
	neg := -0.5
	doc := &ast.Document{Path: "v.md", Subject: "v", Cases: []ast.Case{{
		Name:  "c",
		Nodes: []ast.Node{{Name: "n", Qualified: &neg, Line: 2}},
	}}}
	err := Validate(doc)
	if err == nil || !strings.Contains(err.Error(), `node "n" qualified -0.5: weight out of range [0, 1]`) {
		t.Fatalf("Validate = %v, want the qualifier range error", err)
	}
}
