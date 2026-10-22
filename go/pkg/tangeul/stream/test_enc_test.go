//ff:func feature=tangl type=util control=sequence
//ff:what testEnc — test helper: a Writer state over tt at the narrow test widths
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// testEnc returns a Writer state over tt (an empty table when nil) at testWire,
// issuing non-primary Edge Ids from tt.primaryEnd (or 1).
func testEnc(tt *tidTable) *enc {
	if tt == nil {
		tt = &tidTable{byRef: map[string]toulmin.TID{}, aliasByName: map[string]toulmin.TID{}, caseByName: map[string]toulmin.TID{}, nodeByCase: map[string]toulmin.TID{}, primaryEnd: 1}
	}
	return &enc{tt: tt, wire: testWire, next: tt.primaryEnd}
}
