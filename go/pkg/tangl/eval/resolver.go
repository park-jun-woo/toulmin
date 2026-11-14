//ff:type feature=tangl type=interface
//ff:what Resolver — host lookup of predicate and value Call functions by alias/name
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/toulmin"

// Resolver resolves host-provided Call references to executable functions.
// Predicate resolves a `holds` predicate call and Value resolves a
// `result of` value call. The bool reports whether the reference resolved;
// an unresolved reference is a compile error, never a runtime lookup.
type Resolver interface {
	Predicate(alias, name string) (func(toulmin.Context, toulmin.Specs) (bool, any), bool)
	Value(alias, name string) (func(toulmin.Context, toulmin.Specs) (any, error), bool)
}
