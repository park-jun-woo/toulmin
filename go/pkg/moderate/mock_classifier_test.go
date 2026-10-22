//ff:type feature=moderate type=model
//ff:what mockClassifier — test helper classifier mock
package moderate

type mockClassifier struct {
	scores map[string]float64
}
