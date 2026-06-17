//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what fillAll — calc's every not-yet-run node to fill inactive state (full pass)
package toulmin

// fillAll runs calc on each node that has not yet been evaluated, so that
// inactive nodes also have their active/evidence state populated.
func (ec *evalContext) fillAll(rules []RuleMeta, ctx Context) {
	for _, r := range rules {
		if !ec.ran[r.Name] {
			ec.calc(r.Name, ctx)
		}
	}
}
