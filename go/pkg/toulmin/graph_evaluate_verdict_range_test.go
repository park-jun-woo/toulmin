//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphEvaluateVerdictRange — an out-of-range or NaN rule verdict surfaces as an Evaluate/Run error naming the rule
package toulmin

import (
	"math"
	"strings"
	"testing"
)

func TestGraphEvaluateVerdictRange(t *testing.T) {
	var errs []error
	for _, bad := range []float64{1.5, -2.0, math.NaN()} {
		b := bad
		g := NewGraph("range")
		w := g.Rule(WarrantA)
		g.Counter(func(ctx Context, specs Specs) (float64, any) { return b, nil }).Attacks(w)
		_, errPlain := g.Evaluate(NewContext())
		_, errTrace := g.Evaluate(NewContext(), EvalOption{Trace: true})
		_, _, errRun := g.Run(NewContext())
		errs = append(errs, errPlain, errTrace, errRun)
	}
	for i, err := range errs {
		if err == nil || !strings.Contains(err.Error(), "out of range [-1, 1]") || !strings.Contains(err.Error(), "rule \"") {
			t.Errorf("case %d: expected range error naming the rule, got %v", i, err)
		}
	}
}
