//ff:func feature=tangl type=util control=sequence
//ff:what TestRender — mode 2 regenerates section headers, list skeleton, and multi-line Rules subtrees
package doc

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// TestRender verifies the core mode-2 renderer: it regenerates locale section
// headers, prefixes top-level and nested statements with the right list marker
// and indent, replays prose verbatim from Source, shifts a multi-line Rules
// subtree, and localizes the heading keyword for KO.
func TestRender(t *testing.T) {
	src := []byte("extra prose\n")
	d := &Doc{
		Source: src,
		Blocks: []Block{
			{Kind: SectionBlock, Ref: StmtRef{Sec: toulmin.SectionCases}},
			{Kind: StmtBlock, Ref: StmtRef{Sec: toulmin.SectionCases, Path: []int{0}}},
			{Kind: StmtBlock, Ref: StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 1, 0}}},
			{Kind: ProseBlock, Span: Span{Off: 0, Len: uint32(len(src))}},
			{Kind: SectionBlock, Ref: StmtRef{Sec: toulmin.SectionRules}},
			{Kind: StmtBlock, Ref: StmtRef{Sec: toulmin.SectionRules, Path: []int{0}}},
		},
	}

	want := "## tangl:Cases\n" +
		"- in case of `x`\n" +
		"  - `a` is a general rule\n" +
		"extra prose\n" +
		"## tangl:Rules\n" +
		"- `r` when\n" +
		"  `a` is 5\n" +
		"  and `b`\n"

	got, err := Render(d, stubCodec{locale: "en"})
	if err != nil {
		t.Fatalf("Render(en): unexpected error: %v", err)
	}
	if string(got) != want {
		t.Errorf("Render(en) mismatch:\n got %q\nwant %q", got, want)
	}

	ko := &Doc{Blocks: []Block{{Kind: SectionBlock, Ref: StmtRef{Sec: toulmin.SectionCases}}}}
	gotKo, err := Render(ko, stubCodec{locale: "ko"})
	if err != nil {
		t.Fatalf("Render(ko): unexpected error: %v", err)
	}
	if string(gotKo) != "## tangl:사안\n" {
		t.Errorf("Render(ko) header = %q, want %q", gotKo, "## tangl:사안\n")
	}

	if _, err := Render(nil, stubCodec{locale: "en"}); err == nil {
		t.Errorf("Render(nil Doc): expected error")
	}
	if _, err := Render(d, nil); err == nil {
		t.Errorf("Render(nil Codec): expected error")
	}

	// A block that fails to render (out-of-bounds ProseBlock span) propagates
	// its error out of Render's loop instead of being swallowed.
	bad := &Doc{
		Source: src,
		Blocks: []Block{
			{Kind: ProseBlock, Span: Span{Off: 0, Len: uint32(len(src) + 1)}},
		},
	}
	if _, err := Render(bad, stubCodec{locale: "en"}); err == nil {
		t.Errorf("Render(bad block): expected error")
	}
}
