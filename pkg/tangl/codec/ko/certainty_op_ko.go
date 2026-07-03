//ff:func feature=tangl type=parser control=selection
//ff:what certaintyOpKo — map a KO certainty threshold token to its canonical English Op
package ko

// certaintyOpKo maps the compound threshold token of a 확신 gate to its canonical
// English Certainty.Op (문법 v0.3 §확신도 게이트): 이상이면→"at least", 초과면→"above",
// 미만이면→"less than", 이하면→"at most". ok is false for any other token.
func certaintyOpKo(token string) (string, bool) {
	switch token {
	case "이상이면":
		return "at least", true
	case "초과면":
		return "above", true
	case "미만이면":
		return "less than", true
	case "이하면":
		return "at most", true
	default:
		return "", false
	}
}
