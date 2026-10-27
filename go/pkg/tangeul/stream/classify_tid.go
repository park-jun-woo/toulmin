//ff:func feature=tangl type=util control=selection
//ff:what classifyTID — issue a primary entity's TID and populate the reference name maps
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// classifyTID issues a TID for the primary entity a StmtRef declares (if any) and
// records the alias, case and in-case-node name maps that reference edges resolve
// through. Attacks, execs and run/check strings declare no entity and take no
// TID.
func classifyTID(tt *tidTable, next *toulmin.TID, ref doc.StmtRef, elem any) error {
	switch ref.Sec {
	case toulmin.SectionSubject, toulmin.SectionDefinitions, toulmin.SectionRules:
		assignTID(tt, next, ref)
	case toulmin.SectionSee:
		tt.aliasByName[elem.(ast.See).Alias] = assignTID(tt, next, ref)
	case toulmin.SectionCases:
		classifyCaseTID(tt, next, ref, elem)
	case toulmin.SectionProvides:
		if len(ref.Path) == 1 || ref.Path[1] == 0 {
			assignTID(tt, next, ref)
		}
	case toulmin.SectionInternal:
		if len(ref.Path) == 1 {
			assignTID(tt, next, ref)
		}
	default:
		return fmt.Errorf("tangeul/stream: unknown section %d in TID pass", uint16(ref.Sec))
	}
	return nil
}
