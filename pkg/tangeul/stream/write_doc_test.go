//ff:func feature=tangl type=util control=sequence
//ff:what TestWriteDoc — WriteDoc's own branches: nil Doc/AST, assignTIDs failure, emitBlock failure, the empty-Blocks loop, and a full success
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestWriteDoc exercises every branch of WriteDoc itself that the existing
// round-trip fixtures do not: the nil-Doc/nil-AST guard, assignTIDs returning an
// error (a StmtBlock whose Ref does not resolve), the Blocks loop running zero
// times, and emitBlock returning an error mid-walk (a ProseBlock span too wide
// for the 16-bit wire width). A full multi-block success is also checked here so
// the file is self-contained; round_trip_test.go and round_trip_rich_test.go
// cover the success path in more depth (including case transitions).
func TestWriteDoc(t *testing.T) {
	t.Run("nil doc", func(t *testing.T) {
		if _, err := WriteDoc(nil); err == nil {
			t.Fatalf("WriteDoc(nil): expected error, got nil")
		}
	})

	t.Run("nil AST", func(t *testing.T) {
		if _, err := WriteDoc(&doc.Doc{}); err == nil {
			t.Fatalf("WriteDoc(&doc.Doc{}): expected error, got nil")
		}
	})

	t.Run("assignTIDs error", func(t *testing.T) {
		d := &doc.Doc{
			AST:    &ast.Document{Subject: "x"},
			Source: []byte("----"),
			Blocks: []doc.Block{
				{Kind: doc.StmtBlock, Span: doc.Span{Off: 0, Len: 4}, Ref: doc.StmtRef{Sec: codebook.SectionSee, Path: []int{0}}},
			},
		}
		if _, err := WriteDoc(d); err == nil {
			t.Fatalf("WriteDoc: expected assignTIDs error (See index out of range), got nil")
		}
	})

	t.Run("zero blocks", func(t *testing.T) {
		d := &doc.Doc{AST: &ast.Document{Subject: "x"}, Source: []byte("hi")}
		b, err := WriteDoc(d)
		if err != nil {
			t.Fatalf("WriteDoc: unexpected error: %v", err)
		}
		if len(b) == 0 {
			t.Errorf("WriteDoc: expected non-empty header+source bytes even with no Blocks")
		}
	})

	t.Run("emitBlock error", func(t *testing.T) {
		d := &doc.Doc{
			AST:    &ast.Document{Subject: "x"},
			Source: []byte("hi"),
			Blocks: []doc.Block{
				{Kind: doc.ProseBlock, Span: doc.Span{Off: 0x10000, Len: 1}},
			},
		}
		if _, err := WriteDoc(d); err == nil {
			t.Fatalf("WriteDoc: expected emitBlock error (span exceeds 16-bit wire width), got nil")
		}
	})

	t.Run("success", func(t *testing.T) {
		b, err := WriteDoc(buildTransferDoc())
		if err != nil {
			t.Fatalf("WriteDoc: unexpected error: %v", err)
		}
		if len(b) == 0 {
			t.Errorf("WriteDoc: expected non-empty output")
		}
	})
}
