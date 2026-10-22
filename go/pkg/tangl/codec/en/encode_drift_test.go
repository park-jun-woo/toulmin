//ff:func feature=tangl type=parser control=sequence
//ff:what TestEncodeDrift — a scanner item with no AST ref propagates the collectAnchors drift error through Encode
package en

import "testing"

// TestEncodeDrift exercises Encode's second error return (the collectAnchors call):
// a Cases `requires` line accepts arbitrary trailing text unconditionally
// (parseRequireItem never inspects it.Children, so the parser silently drops a
// nested child under it), while the encoder's collectItemAnchors still walks that
// orphaned nested line looking for a StmtRef. The document is otherwise valid and
// the parser accepts it, so the resulting error can only come from collectAnchors,
// not from parser.ParseSource (already covered by TestEncodeParserError) nor from
// matchBodyRefs (this document has no Provides/Internal section).
func TestEncodeDrift(t *testing.T) {
	src := []byte("## tangl:Subject\n" +
		"- this document is `x`\n" +
		"## tangl:Cases\n" +
		"- in case of `c`\n" +
		"  - `q` is required\n" +
		"    - orphaned nested line under a require\n")
	d, err := ENCodec{}.Encode(src, "drift.md")
	if err == nil {
		t.Fatal("Encode accepted an orphaned nested line under a require with no AST ref")
	}
	if d != nil {
		t.Fatalf("Encode returned a Doc alongside an error: %+v", d)
	}
}
