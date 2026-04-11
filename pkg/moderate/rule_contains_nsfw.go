//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsNSFW: spec(ClassifierSpec)으로 NSFW 감지
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ContainsNSFW checks if the content contains NSFW material.
// spec is *ClassifierSpec. Returns score as evidence.
func ContainsNSFW(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	body, _ := ctx.Get("body")
	if len(specs) == 0 {
		return false, nil
	}
	cb := specs[0].(*ClassifierSpec)
	score := cb.Classifier.Predict(body.(string), "nsfw")
	return score > 0.8, score
}
