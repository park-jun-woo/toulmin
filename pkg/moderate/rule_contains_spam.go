//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsSpam: backing(ClassifierBacking)으로 스팸 감지
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ContainsSpam checks if the content contains spam.
// backing is *ClassifierBacking. Returns score as evidence.
func ContainsSpam(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	body, _ := ctx.Get("body")
	cb := backing.(*ClassifierBacking)
	score := cb.Classifier.Predict(body.(string), "spam")
	return score > 0.7, score
}
