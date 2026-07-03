//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what Render — mode 2 decode: walk Blocks, PROSE verbatim + StmtBlock via codec + section headers
package doc

import (
	"bytes"
	"errors"
)

// Render is mode 2 decoding (cross-locale / normalizing): it walks Blocks in
// stream order and appends each block's rendering — ProseBlock bytes verbatim
// from Source, SectionBlock headings regenerated in the codec's locale, and
// StmtBlock statements via the codec's RenderStmt with core-owned list markers
// and indentation (README contract §3). The codec supplies statement text only.
func Render(d *Doc, c Codec) ([]byte, error) {
	if d == nil {
		return nil, errors.New("tangeul/doc: Render on nil Doc")
	}
	if c == nil {
		return nil, errors.New("tangeul/doc: Render with nil Codec")
	}
	var out bytes.Buffer
	for i, b := range d.Blocks {
		chunk, err := renderBlock(d, c, b, i)
		if err != nil {
			return nil, err
		}
		out.Write(chunk)
	}
	return out.Bytes(), nil
}
