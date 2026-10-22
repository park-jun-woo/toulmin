//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what assignTIDs — pass 1: issue declaration-order TIDs to every primary entity
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// assignTIDs is the Writer's first pass: it walks Blocks in stream order and, via
// classifyTID, issues an ascending segment-local TID to every primary
// (referenceable) entity so forward references (a `run <case>` before its case,
// a `when <node>` cascade) resolve in pass 2. Edges and run/check strings declare
// no entity and take no TID. primaryEnd is the base for auxiliary node TIDs.
func assignTIDs(d *doc.Doc) (*tidTable, error) {
	tt := &tidTable{
		byRef:       map[string]toulmin.TID{},
		aliasByName: map[string]toulmin.TID{},
		caseByName:  map[string]toulmin.TID{},
		nodeByCase:  map[string]toulmin.TID{},
	}
	next := toulmin.TID(1)
	for _, b := range d.Blocks {
		if b.Kind != doc.StmtBlock {
			continue
		}
		elem, err := doc.ResolveStmt(d, b.Ref)
		if err != nil {
			return nil, err
		}
		// ResolveStmt has already rejected unknown sections, the only
		// classifyTID error, so it cannot fail here.
		_ = classifyTID(tt, &next, b.Ref, elem)
	}
	tt.primaryEnd = next
	return tt, nil
}
