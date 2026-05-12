//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsSpam: spec(ClassifierSpec)으로 스팸 감지
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ContainsSpam checks if the content contains spam.
// spec is *ClassifierSpec. Returns score as evidence.
func ContainsSpam(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	body, _ := ctx.Get("body")
	if len(specs) == 0 {
		return false, nil
	}
	cb := specs[0].(*ClassifierSpec)
	s, ok := body.(string)
	if !ok {
		return false, nil
	}
	score := cb.Classifier.Predict(s, "spam")
	return score > 0.7, score
}
