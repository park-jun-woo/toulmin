//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCheckingQualified — checking/qualified node clauses, not used by the spec's three examples
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestParseCheckingQualified exercises "checking `case`" verdict composition
// and "qualified <float>" clauses, neither of which appear in the spec's
// three worked examples.
func TestParseCheckingQualified(t *testing.T) {
	src := "" +
		"## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `inner`\n" +
		"  - `a` is a general rule\n\n" +
		"- in case of `outer`\n" +
		"  - `b` is a general rule using `svc`.`check` qualified 0.8\n" +
		"  - `c` is a general rule checking `inner`\n"
	doc, err := ParseSource(src, "checking_qualified.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	if len(doc.Cases) != 2 {
		t.Fatalf("len(Cases) = %d, want 2", len(doc.Cases))
	}
	outer := doc.Cases[1]
	if len(outer.Nodes) != 2 {
		t.Fatalf("outer.Nodes = %+v, want 2", outer.Nodes)
	}
	b := outer.Nodes[0]
	if b.Qualified == nil || *b.Qualified != 0.8 {
		t.Errorf("b.Qualified = %v, want 0.8", b.Qualified)
	}
	c := outer.Nodes[1]
	if c.Checking != "inner" || c.Role != ast.GeneralRule {
		t.Errorf("c = %+v, want checking 'inner'", c)
	}
}
