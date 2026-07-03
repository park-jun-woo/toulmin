//ff:func feature=tangl type=util control=sequence
//ff:what closeCase — emit a trailing CASE_END if a case is still open at end of stream
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// closeCase emits a CASE_END for a case left open when the Blocks walk ends (the
// last section was Cases), so every CASE_BEGIN is balanced.
func closeCase(st *caseState) ([]word.Word, error) {
	if !st.open {
		return nil, nil
	}
	return caseControlWords(codebook.ControlCaseEnd, st.tid)
}
