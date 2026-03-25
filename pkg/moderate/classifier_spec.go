//ff:type feature=moderate type=model
//ff:what ClassifierSpec: AI 분류기 기반 rule의 spec 타입
package moderate

// ClassifierSpec carries a Classifier for content classification rules.
type ClassifierSpec struct {
	Classifier Classifier // AI classifier implementation
}
