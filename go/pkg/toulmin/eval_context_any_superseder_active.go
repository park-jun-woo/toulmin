//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what anySupersederActive — force-calc every superseder, report whether any is Active
package toulmin

// anySupersederActive force-evaluates every superseder (ec.calc is idempotent via ec.ran
// caching) so its Active state is known, and reports whether at least one is Active. Every
// superseder is calc'd unconditionally — the loop never short-circuits — so an unreached
// superseder is still marked ran even when an earlier one already decided the result.
func (ec *evalContext) anySupersederActive(superseders []string, ctx Context) bool {
	active := false
	for _, from := range superseders {
		ec.calc(from, ctx)
		if ec.active[from] {
			active = true
		}
	}
	return active
}
