//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsHateSpeech: spec(ClassifierSpec)으로 혐오 표현 감지
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ContainsHateSpeech checks if the content contains hate speech.
// spec is *ClassifierSpec. Returns score as evidence.
func ContainsHateSpeech(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	body, _ := ctx.Get("body")
	if len(specs) == 0 {
		return false, nil
	}
	cb := specs[0].(*ClassifierSpec)
	score := cb.Classifier.Predict(body.(string), "hate_speech")
	return score > 0.8, score
}
