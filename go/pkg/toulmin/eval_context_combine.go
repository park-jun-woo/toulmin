//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what combine — folds attackers into an active node's verdict (rebuttal denominator, exception gate)
package toulmin

// combine computes the verdict of an active node id from its own rule verdict v
// and its attackers, evaluated through eval:
//
//	a      = (v+1)/2
//	S      = Σ_rebuttal F·(v_j+1)/2
//	a_eff  = a · Π_except (1 − F·(v_j+1)/2)
//	raw    = q · a_eff / (1 + S)
//	result = 2·raw − 1
//
// Strict nodes ignore every attacker (no denominator, no gate).
func (ec *evalContext) combine(id string, eval func(string) float64) float64 {
	a := (ec.ruleVerdict[id] + 1.0) / 2.0
	if ec.strMap[id] == Strict {
		return ec.settle(id, a, 0.0)
	}
	sum := 0.0
	gate := 1.0
	for i, aid := range ec.edges[id] {
		x := ec.edgeForce(id, i) * (eval(aid) + 1.0) / 2.0
		if ec.isExceptAttacker(aid) {
			gate *= 1.0 - x
			continue
		}
		sum += x
	}
	return ec.settle(id, a*gate, sum)
}
