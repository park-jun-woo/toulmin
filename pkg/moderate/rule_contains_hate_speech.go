//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsHateSpeech: backing(ClassifierBacking)으로 혐오 표현 감지
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ContainsHateSpeech checks if the content contains hate speech.
// backing is *ClassifierBacking. Returns score as evidence.
func ContainsHateSpeech(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	body, _ := ctx.Get("body")
	cb := backing.(*ClassifierBacking)
	score := cb.Classifier.Predict(body.(string), "hate_speech")
	return score > 0.8, score
}
