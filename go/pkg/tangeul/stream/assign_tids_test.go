//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAssignTIDs — pass 1 issues ascending declaration-order TIDs and populates name maps
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestAssignTIDs verifies the Writer's first pass: every primary entity gets an
// ascending segment-local TID in document (Blocks) order, and the alias, case
// and in-case-node name maps resolve to those TIDs. Edges (attacks/execs) and
// run/check strings take no TID, so the counter skips them.
func TestAssignTIDs(t *testing.T) {
	tt, err := assignTIDs(buildTransferDoc())
	if err != nil {
		t.Fatalf("assignTIDs: %v", err)
	}
	got := func(sec toulmin.SectionCode, p ...int) toulmin.TID {
		return tt.byRef[tidKey(doc.StmtRef{Sec: sec, Path: p})]
	}
	// Declaration order: subject=1, sees 2/3, case0=4, its two requires 5/6, its
	// two rule nodes 7/8; case1=9 (attack/exec in case0 took none between 8 and 9).
	checks := []struct {
		name string
		got  toulmin.TID
		want toulmin.TID
	}{
		{"subject", got(toulmin.SectionSubject), 1},
		{"see0", got(toulmin.SectionSee, 0), 2},
		{"see1", got(toulmin.SectionSee, 1), 3},
		{"case0", got(toulmin.SectionCases, 0), 4},
		{"case0 req1", got(toulmin.SectionCases, 0, 0, 1), 6},
		{"case0 node0", got(toulmin.SectionCases, 0, 1, 0), 7},
		{"case1", got(toulmin.SectionCases, 1), 9},
		{"endpoint", got(toulmin.SectionProvides, 0), 12},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s TID = %d, want %d", c.name, c.got, c.want)
		}
	}
	if tt.aliasByName["bank"] != 2 || tt.aliasByName["log"] != 3 {
		t.Errorf("alias map = %v", tt.aliasByName)
	}
	if tt.caseByName["can withdraw"] != 4 || tt.caseByName["can deposit"] != 9 {
		t.Errorf("case map = %v", tt.caseByName)
	}
	if tt.nodeByCase[nodeKey(0, "balance sufficient")] != 7 {
		t.Errorf("node map missing balance sufficient: %v", tt.nodeByCase)
	}
	if tt.primaryEnd != 16 {
		t.Errorf("primaryEnd = %d, want 16", tt.primaryEnd)
	}
	t.Run("extra", TestAssignTIDsExtra)
}
