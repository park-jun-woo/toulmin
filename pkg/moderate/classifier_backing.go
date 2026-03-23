//ff:type feature=moderate type=model
//ff:what ClassifierBacking: AI 분류기 기반 rule의 backing 타입
package moderate

// ClassifierBacking carries a Classifier for content classification rules.
type ClassifierBacking struct {
	Classifier Classifier // AI classifier implementation
}
