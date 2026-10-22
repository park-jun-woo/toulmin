//ff:func feature=engine type=engine control=sequence
//ff:what Supersedes — declares that this rule supersedes (crisply excludes) the target rule
package toulmin

// Supersedes declares a supersede edge: when this rule is Active, target is
// crisply excluded from the tally. Unlike Attacks, this edge never enters the
// h-Categoriser traversal — it only sets target's TraceEntry.Excluded.
func (r *Rule) Supersedes(target *Rule) {
	r.graph.supersedes = append(r.graph.supersedes, defeatEdge{
		from: r.id,
		to:   target.id,
	})
}
