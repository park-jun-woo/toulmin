//ff:func feature=tangl type=util control=sequence
//ff:what TestParseChildren — parseChildren rebuilds n subtrees in order, handles n=0 and propagates a child error
package packet

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestParseChildren checks parseChildren returns an empty slice for n=0,
// rebuilds two Compare children in order for n=2, and propagates an error
// from a malformed child record.
func TestParseChildren(t *testing.T) {
	pos := 0
	terms, count, err := parseChildren(nil, &pos, nil, 0)
	if err != nil || len(terms) != 0 || count != 0 {
		t.Fatalf("n=0: terms=%v count=%d err=%v", terms, count, err)
	}

	rec := exprCompareRecord()
	recs := append(append([]core.Word(nil), rec...), rec...)
	pos = 0
	terms, count, err = parseChildren(recs, &pos, []string{"f"}, 2)
	if err != nil || len(terms) != 2 || count != 2 || pos != len(recs) {
		t.Fatalf("n=2: terms=%v count=%d pos=%d err=%v", terms, count, pos, err)
	}

	pos = 0
	if _, _, err := parseChildren([]core.Word{exprRecHeader(ExprCompare, 0)}, &pos, nil, 1); err == nil {
		t.Fatalf("malformed child: want error")
	}
}
