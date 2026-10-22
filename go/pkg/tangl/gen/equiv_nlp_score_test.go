//ff:func feature=tangl type=codegen control=sequence
//ff:what equivNlpScore — cross-test host value function (result of nlpScore)
package gen

import (
	"errors"

	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// equivNlpScore is the cross-test host value function bound to the bare
// `result of nlpScore` call: it fails when ctx nlpErr is true, else returns
// the ctx nlpScore value. The generated runner defines an identical copy
// named nlpScore.
func equivNlpScore(ctx toulmin.Context, specs toulmin.Specs) (any, error) {
	if v, ok := ctx.Get("nlpErr"); ok && v == true {
		return nil, errors.New("nlp unavailable")
	}
	v, _ := ctx.Get("nlpScore")
	return v, nil
}
