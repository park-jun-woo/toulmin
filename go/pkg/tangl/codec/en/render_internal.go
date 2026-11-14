//ff:func feature=tangl type=parser control=selection
//ff:what renderInternal — render an Internal trigger head (`on <event>` / `every <interval> [until `C`]`)
package en

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderInternal renders a tangl:Internal trigger head back to "on <event>" or
// "every <interval>" with an optional "until `<case>`" clause (its run/check body
// items render as separate StmtBlocks). Event and Interval are stored verbatim, so
// they are emitted as-is. It inverts the parser's parseInternalHead.
func renderInternal(in ast.Internal) (string, error) {
	switch in.Kind {
	case ast.OnEvent:
		return "on " + in.Event, nil
	case ast.EveryTick:
		s := "every " + in.Interval
		if in.Until != "" {
			s += " until `" + in.Until + "`"
		}
		return s, nil
	default:
		return "", fmt.Errorf("en: unknown InternalKind %d", int(in.Kind))
	}
}
