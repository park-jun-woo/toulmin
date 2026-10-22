//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what resolveSupersedes — post-hoc pass computing which targets are crisply excluded
package toulmin

// resolveSupersedes forces every superseder to be evaluated (ec.calc is idempotent via
// ec.ran caching) so its Active state is known, then returns excluded (target → true when
// any of its superseders is Active). This never touches the h-Categoriser traversal; it is
// a reporting-only pass whose result populates TraceEntry.Excluded. Keys are full ids.
func (ec *evalContext) resolveSupersedes(ctx Context) map[string]bool {
	excluded := make(map[string]bool)
	for target, superseders := range ec.supersedeBy {
		if ec.anySupersederActive(superseders, ctx) {
			excluded[target] = true
		}
	}
	return excluded
}
