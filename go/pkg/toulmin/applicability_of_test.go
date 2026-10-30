//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what applicabilityOf — test helper: Applicability of the named warrant result (-1 if absent)
package toulmin

// applicabilityOf returns the Applicability of the result named name, or -1 when absent.
func applicabilityOf(results []EvalResult, name string) float64 {
	for _, r := range results {
		if r.Name == name {
			return r.Applicability
		}
	}
	return -1
}
