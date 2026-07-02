//ff:func feature=tangl type=codegen control=sequence
//ff:what renderEndpoint — dispatches one Provides entry to its run or check renderer
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderEndpoint dispatches one Provides entry to its Run-mode or
// Check-mode renderer: Runs (if any) always drives the exported
// function's signature and side effects — any Checks are then appended
// as extra Evaluate results; a Checks-only endpoint is a pure Evaluate
// passthrough.
func renderEndpoint(w *strings.Builder, ep ast.Endpoint) {
	fnName := goIdentExported(ep.Name)
	fields := requireFields(ep.Requires)
	if len(ep.Runs) > 0 {
		renderEndpointRun(w, fnName, fields, ep.Runs, ep.Checks)
		return
	}
	if len(ep.Checks) > 0 {
		renderEndpointCheck(w, fnName, fields, ep.Checks)
	}
}
