//ff:type feature=tangl type=interface
//ff:what Codec — the surface encoder/decoder contract implemented by the EN/KO locales
package doc

// Codec is the surface (text) encoder/decoder contract that each locale
// implements (README contract §3). Mode 1 (verbatim) needs no codec — it is a
// Source dump via RenderVerbatim. Mode 2 (cross-locale/normalizing) is driven by
// the core renderer Render, which calls RenderStmt for statement text only; the
// markdown skeleton (section headers, list markers, indentation) is the core's
// concern.
type Codec interface {
	// Locale returns the codec's locale tag, "en" or "ko".
	Locale() string
	// Encode absorbs surface text into a Doc; src becomes Doc.Source in full.
	Encode(src []byte, path string) (*Doc, error)
	// RenderStmt regenerates one AST element as this locale's canonical statement
	// text: a single line by default, a multi-line subtree for a Rules inline rule.
	RenderStmt(d *Doc, ref StmtRef) (string, error)
}
