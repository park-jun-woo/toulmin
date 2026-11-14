//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitBlock — ProseBlock/SectionBlock/StmtBlock success and error paths, plus unknown BlockKind
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitBlock exercises every BlockKind branch of emitBlock: a ProseBlock, a
// SectionBlock (success and an oversized-span proseWords failure), a StmtBlock
// (success, a ResolveStmt failure, and an oversized-span toWireSpan failure),
// and the default unknown-BlockKind hard error.
func TestEmitBlock(t *testing.T) {
	d := &doc.Doc{AST: &ast.Document{Subject: "sys"}}
	small := doc.Span{Off: 0, Len: 3}
	huge := doc.Span{Off: 0, Len: 0x10000}

	cases := []struct {
		name    string
		b       doc.Block
		wantErr bool
	}{
		{"prose", doc.Block{Kind: doc.ProseBlock, Span: small}, false},
		{"section success", doc.Block{Kind: doc.SectionBlock, Span: small, Ref: doc.StmtRef{Sec: toulmin.SectionSubject}}, false},
		{"section oversized span", doc.Block{Kind: doc.SectionBlock, Span: huge, Ref: doc.StmtRef{Sec: toulmin.SectionSubject}}, true},
		{"stmt success", doc.Block{Kind: doc.StmtBlock, Span: small, Ref: doc.StmtRef{Sec: toulmin.SectionSubject}}, false},
		{"stmt resolve error", doc.Block{Kind: doc.StmtBlock, Span: small, Ref: doc.StmtRef{Sec: toulmin.SectionSubject, Path: []int{0}}}, true},
		{"stmt oversized span", doc.Block{Kind: doc.StmtBlock, Span: huge, Ref: doc.StmtRef{Sec: toulmin.SectionSubject}}, true},
		{"unknown kind", doc.Block{Kind: doc.BlockKind(99), Span: small}, true},
	}
	for _, c := range cases {
		e := testEnc(&tidTable{})
		got, err := emitBlock(e, d, c.b)
		if err == nil {
			err = e.err
		}
		if c.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
		}
		if c.wantErr {
			continue
		}
		if err != nil {
			t.Fatalf("%s: emitBlock: %v", c.name, err)
		}
		if len(got) == 0 {
			t.Errorf("%s: expected non-empty words", c.name)
		}
	}
}
