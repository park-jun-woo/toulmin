//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what anySupersederApplicable — force-calc every superseder, report whether any has a_eff > 0
package toulmin

// anySupersederApplicable force-evaluates every superseder (ec.calc is idempotent via ec.ran
// caching) so its effective applicability is known, and reports whether at least one has
// a_eff > 0. a_eff = (v+1)/2 · Π_except (1 − F·(v_j+1)/2): exceptions that gate the
// superseder count, rebuttals (Counter) do not, and a Strict superseder keeps a_eff = a.
// A non-Active superseder never reaches combine, so its a_eff stays 0. Every superseder is
// calc'd unconditionally — the loop never short-circuits — so an unreached superseder is
// still marked ran even when an earlier one already decided the result.
func (ec *evalContext) anySupersederApplicable(superseders []string, ctx Context) bool {
	applicable := false
	for _, from := range superseders {
		ec.calc(from, ctx)
		if ec.applicability[from] > 0 {
			applicable = true
		}
	}
	return applicable
}
