//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — lazily evaluates rules by graph traversal and returns verdicts
package toulmin

// Evaluate traverses the defeats graph from each warrant node,
// lazily executing rule funcs only when reached. Returns verdicts
// for warrant nodes. Funcs are cached across warrant evaluations.
func (b *GraphBuilder) Evaluate(claim any, ground any) []EvalResult {
	fnMap := make(map[string]func(any, any) bool)
	qualMap := make(map[string]float64)
	strMap := make(map[string]Strength)
	for _, r := range b.rules {
		fnMap[r.Name] = r.Fn
		qualMap[r.Name] = r.Qualifier
		strMap[r.Name] = r.Strength
	}
	edges := make(map[string][]string)
	for _, d := range b.defeats {
		edges[d.to] = append(edges[d.to], d.from)
	}
	ran := make(map[string]bool)
	active := make(map[string]bool)
	var calc func(string, int) float64
	calc = func(id string, depth int) float64 {
		if depth >= maxDepth {
			return 0.0
		}
		if !ran[id] {
			ran[id] = true
			active[id] = fnMap[id](claim, ground)
		}
		if !active[id] {
			return -1.0
		}
		sum := 0.0
		if strMap[id] != Strict {
			for _, aid := range edges[id] {
				raw := (calc(aid, depth+1) + 1.0) / 2.0
				sum += raw
			}
		}
		raw := qualMap[id] / (1.0 + sum)
		return 2*raw - 1
	}
	attackerSet := collectAttackers(b.defeats)
	var results []EvalResult
	for _, r := range b.rules {
		if attackerSet[r.Name] || r.Strength == Defeater {
			continue
		}
		verdict := calc(r.Name, 0)
		if !active[r.Name] {
			continue
		}
		results = append(results, EvalResult{Name: r.Name, Verdict: verdict})
	}
	return results
}
