//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseStructDef — a struct-form Definitions entry, not used by the spec's three examples
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestParseStructDef exercises the "means" + nested "has `f` as Type" form.
func TestParseStructDef(t *testing.T) {
	src := "" +
		"## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:Definitions\n" +
		"- `applicant` means\n" +
		"  - has `annual income` as Currency\n" +
		"  - has `credit score` as Integer\n\n" +
		"## tangl:Cases\n- in case of `x`\n  - `a` is required\n"
	doc, err := ParseSource(src, "struct_def.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	if len(doc.Defs) != 1 || doc.Defs[0].Kind != ast.StructDef {
		t.Fatalf("Defs = %+v, want 1 StructDef", doc.Defs)
	}
	fields := doc.Defs[0].Fields
	if len(fields) != 2 || fields[0].Name != "annual income" || fields[0].Type != "Currency" {
		t.Errorf("Fields = %+v", fields)
	}
}
