//ff:func feature=engine type=engine control=sequence
//ff:what settle — turns applicability and rebuttal sum into a cached verdict
package toulmin

// settle records applicability aEff for id, computes raw = q·aEff/(1+sum),
// caches and returns the verdict 2·raw − 1.
func (ec *evalContext) settle(id string, aEff, sum float64) float64 {
	ec.applicability[id] = aEff
	raw := ec.qualMap[id] * aEff / (1.0 + sum)
	v := 2*raw - 1
	ec.verdictCache[id] = v
	return v
}
