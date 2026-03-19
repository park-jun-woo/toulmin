//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsSpam: backing(Classifier)으로 스팸 감지
package moderate

// ContainsSpam checks if the content contains spam.
// backing is Classifier. Returns score as evidence.
func ContainsSpam(claim any, ground any, backing any) (bool, any) {
	content := claim.(*Content)
	classifier := backing.(Classifier)
	score := classifier.Predict(content.Body, "spam")
	return score > 0.7, score
}
