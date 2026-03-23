//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsNSFW: backing(ClassifierBacking)으로 NSFW 감지
package moderate

// ContainsNSFW checks if the content contains NSFW material.
// backing is *ClassifierBacking. Returns score as evidence.
func ContainsNSFW(claim any, ground any, backing any) (bool, any) {
	content := claim.(*Content)
	cb := backing.(*ClassifierBacking)
	score := cb.Classifier.Predict(content.Body, "nsfw")
	return score > 0.8, score
}
