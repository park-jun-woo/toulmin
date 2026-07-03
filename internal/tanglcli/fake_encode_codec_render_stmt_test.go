//go:build fakecodec

//ff:func feature=cli type=command control=sequence
//ff:what fakeEncodeCodec.RenderStmt — unused stub; runEncode never renders statement text
package tanglcli

import "github.com/park-jun-woo/toulmin/pkg/tangeul/doc"

// RenderStmt is unused by runEncode and always returns an empty statement and
// nil error.
func (f *fakeEncodeCodec) RenderStmt(d *doc.Doc, ref doc.StmtRef) (string, error) {
	return "", nil
}
