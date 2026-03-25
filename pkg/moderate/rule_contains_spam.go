//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsSpam: spec(ClassifierSpec)으로 스팸 감지
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ContainsSpam checks if the content contains spam.
// spec is *ClassifierSpec. Returns score as evidence.
func ContainsSpam(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	body, _ := ctx.Get("body")
	cb := specs[0].(*ClassifierSpec)
	score := cb.Classifier.Predict(body.(string), "spam")
	return score > 0.7, score
}
