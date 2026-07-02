//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what requireFields — extracts ctx field names from is-required clauses
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// requireFields extracts the ctx field names from an endpoint's or
// internal's "is required" clauses, in document order.
func requireFields(reqs []ast.Require) []string {
	fields := make([]string, len(reqs))
	for i, r := range reqs {
		fields[i] = r.Field
	}
	return fields
}
