//ff:type feature=analyzer type=model
//ff:what graphCollector — mutable collector used during AST chain walk
package analyzer

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// graphCollector accumulates rules and defeats during AST traversal.
type graphCollector struct {
	def toulmin.GraphDef
}
