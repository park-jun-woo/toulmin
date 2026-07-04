//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what groundsFingerprint — joins every node's Grounds into one stable string for determinism comparison
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// groundsFingerprint renders every node's Grounds slice, in document order,
// into a single string so two parses of the same fixture can be compared for
// byte-identical ground lists.
func groundsFingerprint(doc *ast.Document) string {
	var b strings.Builder
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			b.WriteString(strings.Join(n.Grounds, "\x00"))
			b.WriteByte('\n')
		}
	}
	return b.String()
}
