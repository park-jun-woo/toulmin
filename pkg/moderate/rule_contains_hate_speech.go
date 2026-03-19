//ff:func feature=moderate type=rule control=sequence
//ff:what ContainsHateSpeech: backing(Classifier)으로 혐오 표현 감지
package moderate

// ContainsHateSpeech checks if the content contains hate speech.
// backing is Classifier. Returns score as evidence.
func ContainsHateSpeech(claim any, ground any, backing any) (bool, any) {
	content := claim.(*Content)
	classifier := backing.(Classifier)
	score := classifier.Predict(content.Body, "hate_speech")
	return score > 0.8, score
}
