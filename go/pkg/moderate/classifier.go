//ff:type feature=moderate type=interface
//ff:what Classifier: AI 콘텐츠 분류 추상화 인터페이스
package moderate

// Classifier abstracts content classification (LLM, ML model, keyword matching, etc.).
type Classifier interface {
	Predict(text string, category string) float64
}
