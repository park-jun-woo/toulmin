//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderStmtSupersede — RenderStmt renders a case supersede slot back to its surface sentence
package en

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

func TestRenderStmtSupersede(t *testing.T) {
	src := []byte("## tangl:Subject\n- this document is `x`\n\n## tangl:Cases\n- in case of `c`\n" +
		"  - `A` is a general rule\n  - `B` is a counter rule\n  - `A` supersedes `B`\n")
	d, err := (ENCodec{}).Encode(src, "t.md")
	if err != nil {
		t.Fatal(err)
	}
	got, err := (ENCodec{}).RenderStmt(d, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 4, 0}})
	if err != nil || got != "`A` supersedes `B`" {
		t.Fatalf("RenderStmt = %q, %v", got, err)
	}
}
