//ff:func feature=tangl type=util control=sequence
//ff:what closeCase — emit a trailing CASE_END if a case is still open at end of stream
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// closeCase emits a CASE_END for a case left open when the Blocks walk ends (the
// last section was Cases), so every CASE_BEGIN is balanced.
func closeCase(e *enc, st *caseState) []core.Word {
	if !st.open {
		return nil
	}
	return controlWords(e, toulmin.ControlCaseEnd, uint64(st.tid), nil)
}
