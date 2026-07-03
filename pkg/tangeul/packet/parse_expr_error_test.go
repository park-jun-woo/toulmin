//ff:func feature=tangl type=util control=sequence
//ff:what TestParseExprErrors — parseExpr rejects a missing/truncated record header and every op/field/child error per ExprKind
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestParseExprErrors exercises parseExpr's self-contained guards across all
// four ExprKind cases: a missing record header, a truncated Compare record,
// an unknown CompareOp code, an out-of-range field label, a propagated
// parseCompareValue error, a missing/unknown Logic op, a propagated
// parseChildren error (Logic/Not/Either) and a Not record whose child count
// is not exactly one.
func TestParseExprErrors(t *testing.T) {
	pos := 0
	if _, _, err := parseExpr(nil, &pos, nil); err == nil {
		t.Fatalf("missing record header: want error")
	}

	pos = 0
	if _, _, err := parseExpr([]word.Word{exprRecHeader(ExprCompare, 0)}, &pos, nil); err == nil {
		t.Fatalf("truncated Compare record: want error")
	}

	pos = 0
	recs := []word.Word{exprRecHeader(ExprCompare, 0), 0x9999, 0, word.Word(ValueNum)}
	if _, _, err := parseExpr(recs, &pos, nil); err == nil {
		t.Fatalf("unknown CompareOp: want error")
	}

	pos = 0
	recs = []word.Word{exprRecHeader(ExprCompare, 0), word.Word(codebook.CompareEquals), 5, word.Word(ValueNum)}
	if _, _, err := parseExpr(recs, &pos, nil); err == nil {
		t.Fatalf("out-of-range field label: want error")
	}

	pos = 0
	recs = []word.Word{exprRecHeader(ExprCompare, 0), word.Word(codebook.CompareEquals), 0, word.Word(ValueNum)}
	if _, _, err := parseExpr(recs, &pos, []string{"f"}); err == nil {
		t.Fatalf("truncated Compare NUM value: want error")
	}

	pos = 0
	if _, _, err := parseExpr([]word.Word{exprRecHeader(ExprLogic, 0)}, &pos, nil); err == nil {
		t.Fatalf("missing Logic op: want error")
	}

	pos = 0
	recs = []word.Word{exprRecHeader(ExprLogic, 0), 0x9999}
	if _, _, err := parseExpr(recs, &pos, nil); err == nil {
		t.Fatalf("unknown LogicOp: want error")
	}

	pos = 0
	recs = []word.Word{exprRecHeader(ExprLogic, 1), word.Word(codebook.LogicAnd)}
	if _, _, err := parseExpr(recs, &pos, nil); err == nil {
		t.Fatalf("Logic propagated child error: want error")
	}

	pos = 0
	if _, _, err := parseExpr([]word.Word{exprRecHeader(ExprNot, 1)}, &pos, nil); err == nil {
		t.Fatalf("Not propagated child error: want error")
	}

	pos = 0
	if _, _, err := parseExpr([]word.Word{exprRecHeader(ExprNot, 0)}, &pos, nil); err == nil {
		t.Fatalf("Not with zero children: want error")
	}

	pos = 0
	if _, _, err := parseExpr([]word.Word{exprRecHeader(ExprEither, 1)}, &pos, nil); err == nil {
		t.Fatalf("Either propagated child error: want error")
	}
}
