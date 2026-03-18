//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what EvaluateTrace — lazily evaluates rules by graph traversal and returns verdicts with per-warrant trace
package toulmin

// EvaluateTrace traverses the defeats graph from each warrant node,
// lazily executing rule funcs only when reached. Returns verdicts with
// per-warrant trace containing only relevant rules. State is reset per warrant.
func (e *Engine) EvaluateTrace(claim any, ground any) []EvalResult {
	fnMap := make(map[string]func(any, any) (bool, any))
	qualMap := make(map[string]float64)
	strMap := make(map[string]Strength)
	edges := make(map[string][]string)
	for _, r := range e.rules {
		fnMap[r.Name] = r.Fn
		qualMap[r.Name] = r.Qualifier
		strMap[r.Name] = r.Strength
		for _, target := range r.Defeats {
			edges[target] = append(edges[target], r.Name)
		}
	}
	attackerSet := make(map[string]bool)
	for _, attackers := range edges {
		for _, aid := range attackers {
			attackerSet[aid] = true
		}
	}
	ran := make(map[string]bool)
	active := make(map[string]bool)
	evidence := make(map[string]any)
	var trace []TraceEntry
	var calc func(string, int) float64
	calc = func(id string, depth int) float64 {
		if depth >= maxDepth {
			return 0.0
		}
		if !ran[id] {
			ran[id] = true
			active[id], evidence[id] = fnMap[id](claim, ground)
			trace = append(trace, TraceEntry{
				Name:      id,
				Role:      inferRole(strMap, attackerSet, id),
				Activated: active[id],
				Qualifier: qualMap[id],
				Evidence:  evidence[id],
			})
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
	var results []EvalResult
	for _, r := range e.rules {
		if len(r.Defeats) > 0 || r.Strength == Defeater {
			continue
		}
		ran = make(map[string]bool)
		active = make(map[string]bool)
		evidence = make(map[string]any)
		trace = nil
		verdict := calc(r.Name, 0)
		if !active[r.Name] {
			continue
		}
		results = append(results, EvalResult{Name: r.Name, Verdict: verdict, Evidence: evidence[r.Name], Trace: trace})
	}
	return results
}
