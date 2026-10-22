//ff:func feature=engine type=engine control=sequence
//ff:what Grounds — records external ground fact names this rule's evaluation depends on (audit metadata only)
package toulmin

// Grounds records the names of external ground facts this rule's Fn resolves
// at evaluation time (e.g. a network lookup behind a predicate). Purely
// additive audit metadata — the engine never resolves, caches, or gates on
// grounds itself; that stays the registry/predicate layer's job (see reins
// Pred.Needs). Multiple calls append, mirroring With's compose semantics.
// Unlike With, Grounds does not affect ruleID — ground names are metadata,
// not judgment-identity-bearing spec state.
func (r *Rule) Grounds(names ...string) *Rule {
	meta := &r.graph.rules[r.idx]
	meta.Grounds = append(meta.Grounds, names...)
	return r
}
