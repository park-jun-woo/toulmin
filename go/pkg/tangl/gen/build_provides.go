//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildProvides — renders every tangl:Provides entry in document order
package gen

import "strings"

// buildProvides renders every tangl:Provides entry as an exported
// endpoint function, in document order.
func buildProvides(w *strings.Builder, gc *genContext) {
	for _, ep := range gc.Doc.Provides {
		renderEndpoint(w, ep)
	}
}
