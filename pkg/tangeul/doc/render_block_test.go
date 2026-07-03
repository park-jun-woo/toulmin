//ff:func feature=tangl type=util control=sequence
//ff:what TestRenderBlock — every BlockKind's success and error branch, plus the unknown-kind default
package doc

import (
	"math"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
)

// TestRenderBlock drives renderBlock directly for each BlockKind: a ProseBlock
// with an in-bounds span (success) and both out-of-bounds shapes (overflowing
// Off+Len, and a span past end of Source); a SectionBlock with a known locale
// (success) and an unknown locale (sectionHeader error); a StmtBlock at the top
// level (level 0) and nested two deep (level 1), plus a section with no canned
// RenderStmt (error); and an unrecognized BlockKind hitting the default branch.
func TestRenderBlock(t *testing.T) {
	src := []byte("hello\n")
	d := &Doc{Source: src}
	c := stubCodec{locale: "en"}

	// ProseBlock: in-bounds span succeeds and returns the verbatim slice.
	got, err := renderBlock(d, c, Block{Kind: ProseBlock, Span: Span{Off: 0, Len: uint32(len(src))}}, 0)
	if err != nil {
		t.Fatalf("ProseBlock in-bounds: unexpected error: %v", err)
	}
	if string(got) != string(src) {
		t.Errorf("ProseBlock in-bounds = %q, want %q", got, src)
	}

	// ProseBlock: Off+Len overflows uint32, so end < Off.
	if _, err := renderBlock(d, c, Block{Kind: ProseBlock, Span: Span{Off: math.MaxUint32, Len: 2}}, 1); err == nil {
		t.Errorf("ProseBlock overflow span: expected error")
	}

	// ProseBlock: end within uint32 range but past len(Source).
	if _, err := renderBlock(d, c, Block{Kind: ProseBlock, Span: Span{Off: 0, Len: uint32(len(src) + 1)}}, 2); err == nil {
		t.Errorf("ProseBlock out-of-bounds span: expected error")
	}

	// SectionBlock: known locale succeeds.
	got, err = renderBlock(d, c, Block{Kind: SectionBlock, Ref: StmtRef{Sec: codebook.SectionCases}}, 3)
	if err != nil {
		t.Fatalf("SectionBlock known locale: unexpected error: %v", err)
	}
	if string(got) != "## tangl:Cases\n" {
		t.Errorf("SectionBlock known locale = %q, want %q", got, "## tangl:Cases\n")
	}

	// SectionBlock: unknown locale propagates sectionHeader's error.
	if _, err := renderBlock(d, stubCodec{locale: "xx"}, Block{Kind: SectionBlock, Ref: StmtRef{Sec: codebook.SectionCases}}, 4); err == nil {
		t.Errorf("SectionBlock unknown locale: expected error")
	}

	// StmtBlock: top-level path (len(Path) < 2) renders at indent level 0.
	got, err = renderBlock(d, c, Block{Kind: StmtBlock, Ref: StmtRef{Sec: codebook.SectionCases, Path: []int{0}}}, 5)
	if err != nil {
		t.Fatalf("StmtBlock level 0: unexpected error: %v", err)
	}
	if string(got) != "- in case of `x`\n" {
		t.Errorf("StmtBlock level 0 = %q, want %q", got, "- in case of `x`\n")
	}

	// StmtBlock: nested path (len(Path) >= 2) renders at indent level 1.
	got, err = renderBlock(d, c, Block{Kind: StmtBlock, Ref: StmtRef{Sec: codebook.SectionCases, Path: []int{0, 1, 0}}}, 6)
	if err != nil {
		t.Fatalf("StmtBlock level 1: unexpected error: %v", err)
	}
	if string(got) != "  - `a` is a general rule\n" {
		t.Errorf("StmtBlock level 1 = %q, want %q", got, "  - `a` is a general rule\n")
	}

	// StmtBlock: a section stubCodec has no canned text for propagates RenderStmt's error.
	if _, err := renderBlock(d, c, Block{Kind: StmtBlock, Ref: StmtRef{Sec: codebook.SectionSubject}}, 7); err == nil {
		t.Errorf("StmtBlock unrenderable section: expected error")
	}

	// default: an unrecognized BlockKind is an error.
	if _, err := renderBlock(d, c, Block{Kind: BlockKind(99)}, 8); err == nil {
		t.Errorf("unknown BlockKind: expected error")
	}
}
