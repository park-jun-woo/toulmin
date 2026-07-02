//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseAccess — parses the spec's access control example and checks key fields
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestParseAccess parses the spec's "접근 제어" example verbatim and checks
// section counts plus the with/attack fields it exercises.
func TestParseAccess(t *testing.T) {
	doc, err := Parse("testdata/access.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if doc.Subject != "api" {
		t.Errorf("Subject = %q, want api", doc.Subject)
	}
	if len(doc.Sees) != 2 {
		t.Fatalf("len(Sees) = %d, want 2", len(doc.Sees))
	}
	if len(doc.Cases) != 1 {
		t.Fatalf("len(Cases) = %d, want 1", len(doc.Cases))
	}
	c := doc.Cases[0]
	if c.Name != "can access" {
		t.Errorf("Cases[0].Name = %q, want 'can access'", c.Name)
	}
	if len(c.Requires) != 1 || c.Requires[0].Field != "user" {
		t.Fatalf("Requires = %+v, want [user]", c.Requires)
	}
	if len(c.Nodes) != 3 {
		t.Fatalf("len(Nodes) = %d, want 3", len(c.Nodes))
	}
	if c.Nodes[0].Using != nil {
		t.Errorf("authenticate.Using = %+v, want nil (local same-name function)", c.Nodes[0].Using)
	}
	blockIP := c.Nodes[1]
	if blockIP.Role != ast.CounterRule || len(blockIP.With) != 1 || blockIP.With[0] != "blocklist" {
		t.Errorf("block ip node = %+v, want CounterRule with [blocklist]", blockIP)
	}
	if len(c.Attacks) != 2 {
		t.Fatalf("len(Attacks) = %d, want 2", len(c.Attacks))
	}
	if len(c.Execs) != 3 {
		t.Fatalf("len(Execs) = %d, want 3", len(c.Execs))
	}
	for i, e := range c.Execs {
		if e.Kind != ast.DoExec || e.Once {
			t.Errorf("Execs[%d] = %+v, want DoExec without once", i, e)
		}
	}
	if len(doc.Provides) != 1 || len(doc.Provides[0].Requires) != 1 || doc.Provides[0].Runs[0] != "can access" {
		t.Fatalf("Provides = %+v", doc.Provides)
	}
}
