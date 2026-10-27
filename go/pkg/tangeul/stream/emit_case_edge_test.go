//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitCaseEdge — undeclared case target is a hard error, a declared target encodes an edge
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestEmitCaseEdge covers both emitCaseEdge branches: an undeclared case name
// (not present in caseByName) is a hard error, and a declared name resolves and
// encodes the edge. appendEdge's own `if err != nil` guard is not reachable here
// since PropTriggers/PropEvaluates/PropAttacks all take the fixed-width DstTID
// slot with no gate flag, so EncodeEdge cannot fail for these fixed shapes.
func TestEmitCaseEdge(t *testing.T) {
	e := testEnc(&tidTable{caseByName: map[string]toulmin.TID{"can deposit": 9}})
	span := toulmin.Span{Off: 0, Len: 1}

	cases := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{"undeclared", "missing", true},
		{"declared", "can deposit", false},
	}
	for _, c := range cases {
		got, err := emitCaseEdge(e, toulmin.PropTriggers, 4, c.target, span)
		if c.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
		}
		if c.wantErr {
			continue
		}
		if err != nil {
			t.Fatalf("%s: emitCaseEdge: %v", c.name, err)
		}
		if len(got) == 0 {
			t.Errorf("%s: expected non-empty words", c.name)
		}
	}
}
