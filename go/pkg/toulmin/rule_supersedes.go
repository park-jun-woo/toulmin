//ff:func feature=engine type=engine control=sequence
//ff:what Supersedes — declares that this rule supersedes (crisply excludes) the target rule
package toulmin

// Supersedes declares a supersede edge: when this rule's effective applicability
// a_eff > 0 (its own (v+1)/2 after Except gates; Counter rebuttals do not count),
// target is crisply excluded from the tally. An exception that fully gates this rule
// therefore lifts the exclusion. Unlike Attacks, this edge never enters the
// h-Categoriser traversal — it only sets target's TraceEntry.Excluded.
func (r *Rule) Supersedes(target *Rule) {
	r.graph.supersedes = append(r.graph.supersedes, defeatEdge{
		from: r.id,
		to:   target.id,
	})
}
