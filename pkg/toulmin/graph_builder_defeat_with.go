//ff:func feature=engine type=engine control=sequence
//ff:what DefeatWith — declares a defeat edge with explicit backing for both sides
package toulmin

// DefeatWith declares a defeat edge between rules with specific backing.
// Use when same function is registered with different backing values.
func (b *GraphBuilder) DefeatWith(fromFn any, fromBacking any, toFn any, toBacking any) *GraphBuilder {
	b.defeats = append(b.defeats, defeatEdge{
		from: ruleID(fromFn, fromBacking),
		to:   ruleID(toFn, toBacking),
	})
	return b
}
