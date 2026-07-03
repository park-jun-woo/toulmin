//ff:func feature=cli type=command control=selection
//ff:what fakeDecodeCodec.RenderStmt — canned statement text, or a forced failure when fail is set
package tanglcli

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// RenderStmt returns a canned statement, or a forced error when f.fail is set —
// the switch that drives runDecode's Render-success and Render-error branches.
func (f fakeDecodeCodec) RenderStmt(d *doc.Doc, ref doc.StmtRef) (string, error) {
	switch f.fail {
	case true:
		return "", errors.New("fakeDecodeCodec: forced RenderStmt failure")
	default:
		return "stmt", nil
	}
}
