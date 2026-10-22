//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildEquivCtx — builds a cross-test Context from one ctx-table map
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/toulmin"

// buildEquivCtx builds a Context from one JSON ctx-table entry, converting
// the authorTrustScore shorthand into an *equivAuthor under "author" so the
// struct-walking path is exercised. The generated runner defines an identical
// copy.
func buildEquivCtx(m map[string]any) toulmin.Context {
	ctx := toulmin.NewContext()
	for k, v := range m {
		if k == "authorTrustScore" {
			ctx.Set("author", &equivAuthor{TrustScore: v.(float64)})
			continue
		}
		ctx.Set(k, v)
	}
	return ctx
}
