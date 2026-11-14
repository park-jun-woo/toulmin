//ff:type feature=tangl type=model
//ff:what docResolver — a Resolver that resolves bare predicates to sibling doc rules first
package eval

// docResolver wraps a host Resolver so that a bare (alias-less) predicate Call
// resolves to a compiled sibling document rule when one exists, delegating
// everything else to the host. ruleSet names the document rules; compiled holds
// their already-built closures.
type docResolver struct {
	host     Resolver
	ruleSet  map[string]bool
	compiled map[string]evalFunc
}
