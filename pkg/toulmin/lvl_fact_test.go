//ff:func feature=engine type=engine control=sequence
//ff:what lvlFact — test helper: always-active rule returning a Fact as evidence
package toulmin

func lvlFact(ctx Context, specs Specs) (bool, any) {
	return true, Fact{Where: "field.x", Expected: 10, Actual: 3}
}
