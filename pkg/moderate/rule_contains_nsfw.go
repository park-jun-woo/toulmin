//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsNSFW: backing(Classifier)으로 NSFW 감지
package moderate

// ContainsNSFW checks if the content contains NSFW material.
// backing is Classifier. Returns score as evidence.
func ContainsNSFW(claim any, ground any, backing any) (bool, any) {
	content := claim.(*Content)
	classifier := backing.(Classifier)
	score := classifier.Predict(content.Body, "nsfw")
	return score > 0.8, score
}
