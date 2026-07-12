//ff:type feature=tangl type=model
//ff:what testAuthor — a *Author-style struct exercising case-insensitive field access
package eval

// testAuthor mirrors the demo *Author: a struct whose TrustScore field is
// reached from a `trustScore` backtick path step.
type testAuthor struct {
	TrustScore float64
}
