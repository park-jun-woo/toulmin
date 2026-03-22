//ff:func feature=moderate type=model control=sequence
//ff:what mockClassifier.Predict — returns mock prediction score for category
package moderate

func (m *mockClassifier) Predict(text string, category string) float64 {
	return m.scores[category]
}
