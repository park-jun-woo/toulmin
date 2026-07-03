//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what assignTIDs — pass 1: issue declaration-order TIDs to every primary entity
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// assignTIDs is the Writer's first pass: it walks Blocks in stream order and, via
// classifyTID, issues an ascending segment-local TID to every primary
// (referenceable) entity so forward references (a `run <case>` before its case,
// a `when <node>` cascade) resolve in pass 2. Edges and run/check strings declare
// no entity and take no TID. primaryEnd is the base for auxiliary node TIDs.
func assignTIDs(d *doc.Doc) (*tidTable, error) {
	tt := &tidTable{
		byRef:       map[string]uint16{},
		aliasByName: map[string]uint16{},
		caseByName:  map[string]uint16{},
		nodeByCase:  map[string]uint16{},
	}
	next := uint16(1)
	for _, b := range d.Blocks {
		if b.Kind != doc.StmtBlock {
			continue
		}
		elem, err := doc.ResolveStmt(d, b.Ref)
		if err != nil {
			return nil, err
		}
		if err := classifyTID(tt, &next, b.Ref, elem); err != nil {
			return nil, err
		}
	}
	tt.primaryEnd = next
	return tt, nil
}
