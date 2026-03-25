//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsSpam: backing(ClassifierBacking)으로 스팸 감지
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ContainsSpam checks if the content contains spam.
// backing is *ClassifierBacking. Returns score as evidence.
func ContainsSpam(claim any, ground any, backing toulmin.Backing) (bool, any) {
	content := claim.(*Content)
	cb := backing.(*ClassifierBacking)
	score := cb.Classifier.Predict(content.Body, "spam")
	return score > 0.7, score
}
