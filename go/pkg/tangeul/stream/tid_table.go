//ff:type feature=tangl type=model
//ff:what tidTable — declaration-order TIDs of primary entities plus name lookups for references
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// tidTable is the pass-1 result: every primary (referenceable) entity's
// segment-local TID, keyed both by its StmtRef coordinate (byRef, the canonical
// key from tidKey) and by the name maps the reference edges resolve through.
// aliasByName resolves `see` aliases used in function references; caseByName
// resolves run/checking/until/endpoint targets; nodeByCase resolves in-case node
// names used by attacks and execution triggers (keyed by caseIdx and name).
// primaryEnd is one past the highest primary TID, the base for auxiliary node
// TIDs. TIDs are issued in stream (declaration) order so §7.2 holds.
type tidTable struct {
	byRef       map[string]toulmin.TID
	aliasByName map[string]toulmin.TID
	caseByName  map[string]toulmin.TID
	nodeByCase  map[string]toulmin.TID
	primaryEnd  toulmin.TID
}
