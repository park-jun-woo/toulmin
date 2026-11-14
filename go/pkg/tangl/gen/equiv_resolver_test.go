//ff:type feature=tangl type=codegen
//ff:what equivResolver — cross-test eval Resolver binding the two bare host functions
package gen

// equivResolver is the eval Resolver for the cross-test document: it binds
// the bare predicate isSpam and the bare value nlpScore to the same host
// behaviors the generated runner compiles in as Go functions.
type equivResolver struct{}
