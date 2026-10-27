//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what emitBody — encode the SOURCE record and every Block's records at the given widths
package stream

import (
	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// emitBody encodes the BODY records of d at widths w: the SOURCE record first,
// then one record run per Block in stream order, with CASE-BEGIN/CASE-END
// inserted by transitionCase as the walk enters and leaves cases. Declaration
// order is the Blocks order, so the TANGEUL order rule holds by construction.
// A record that fails to encode is recorded on e and reported once the walk
// ends.
func emitBody(d *doc.Doc, tt *tidTable, w toulmin.Wire) ([]core.Word, error) {
	e := &enc{tt: tt, wire: w, next: tt.primaryEnd}
	out, err := toulmin.EncodeSource(toulmin.TextPacket{TID: issueID(e), Bytes: d.Source}, w)
	if err != nil {
		return nil, err
	}
	st := &caseState{}
	for _, b := range d.Blocks {
		trans := transitionCase(e, st, b)
		bw, err := emitBlock(e, d, b)
		if err != nil {
			return nil, err
		}
		out = append(append(out, trans...), bw...)
	}
	if out = append(out, closeCase(e, st)...); e.err != nil {
		return nil, e.err
	}
	return out, nil
}
