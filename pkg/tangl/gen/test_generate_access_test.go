//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerateAccess — generates the spec's access control example and checks a pure run wrapper
package gen

import (
	"go/format"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestGenerateAccess compiles the spec's access-control example and
// checks the "access" endpoint compiles to the run wrapper with no
// Evaluate call (it has no "checking" node and no check endpoint), plus
// the counter-node deny action and the required-field guard.
func TestGenerateAccess(t *testing.T) {
	doc, err := parser.Parse("../parser/testdata/access.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	src, err := Generate(doc)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := format.Source([]byte(src)); err != nil {
		t.Fatalf("re-format.Source: %v", err)
	}
	if !strings.Contains(src, "package api") {
		t.Errorf("missing package clause; got:\n%s", src)
	}
	if !strings.Contains(src, `tangl.Required(ctx, "user")`) {
		t.Errorf("missing Required guard for the access endpoint")
	}
	if !strings.Contains(src, "func Access(ctx toulmin.Context) error {") {
		t.Errorf("missing the access endpoint function")
	}
	if strings.Contains(src, ".Evaluate(") {
		t.Errorf("access.md has no checking/check case; got an unexpected Evaluate call:\n%s", src)
	}
	if !strings.Contains(src, "blockIp.RunOn(func(self toulmin.TraceEntry, t toulmin.Trace) error {") {
		t.Errorf("missing the counter-node deny action RunOn handler")
	}
	if !strings.Contains(src, "policy.Deny(t.Ctx())") {
		t.Errorf("missing the policy.Deny call on the block ip counter node")
	}
}
