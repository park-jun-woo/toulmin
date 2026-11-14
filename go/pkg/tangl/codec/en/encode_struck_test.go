//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestEncodeStruckProse — a struck list item is lowered to prose, not a statement
package en

import (
	"bytes"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestEncodeStruckProse encodes a case containing a `~~...~~` struck item and
// asserts the struck line is covered by a ProseBlock (not a StmtBlock, matching the
// parser's filterStruck), so the StmtBlock count still equals the four real AST
// terminals and the tiling holds.
func TestEncodeStruckProse(t *testing.T) {
	src := []byte("## tangl:Subject\n" +
		"- this document is `x`\n" +
		"\n" +
		"## tangl:Cases\n" +
		"- in case of `c`\n" +
		"  - `r` is a general rule using `p`.`fn`\n" +
		"  - ~~commented out~~\n" +
		"  - do `p`.`act` when `r`\n")

	d, err := ENCodec{}.Encode(src, "struck.md")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	stmts := 0
	for _, b := range d.Blocks {
		if b.Kind == doc.StmtBlock {
			stmts++
		}
	}
	if stmts != 4 { // Subject + case + node + exec
		t.Errorf("StmtBlocks = %d, want 4", stmts)
	}

	needle := []byte("~~commented out~~")
	inProse := false
	for _, b := range d.Blocks {
		if b.Kind != doc.ProseBlock {
			continue
		}
		seg := d.Source[b.Span.Off : b.Span.Off+b.Span.Len]
		if bytes.Contains(seg, needle) {
			inProse = true
		}
	}
	if !inProse {
		t.Error("struck item was not lowered into a ProseBlock")
	}
}
