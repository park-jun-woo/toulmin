//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsNSFW: backing(ClassifierBacking)으로 NSFW 감지
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ContainsNSFW checks if the content contains NSFW material.
// backing is *ClassifierBacking. Returns score as evidence.
func ContainsNSFW(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	body, _ := ctx.Get("body")
	cb := backing.(*ClassifierBacking)
	score := cb.Classifier.Predict(body.(string), "nsfw")
	return score > 0.8, score
}
