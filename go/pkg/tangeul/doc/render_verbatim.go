//ff:func feature=tangl type=util control=sequence
//ff:what RenderVerbatim — mode 1 decode: dump Doc.Source verbatim (byte identity)
package doc

// RenderVerbatim is mode 1 decoding (authoring locale): it returns Doc.Source
// unchanged. Because the tiling covers the whole Source, replaying spans equals a
// Source dump, so byte identity holds structurally (README contract §3). Callers
// that must not alias the Doc should copy the result.
func RenderVerbatim(d *Doc) []byte {
	if d == nil {
		return nil
	}
	return d.Source
}
