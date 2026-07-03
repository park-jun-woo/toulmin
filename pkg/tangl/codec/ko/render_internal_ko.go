//ff:func feature=tangl type=parser control=selection
//ff:what renderInternalKo — render an Internal trigger head (`` `이벤트` 발생 시 `` / `` 매 <간격>마다[, `C` 까지] ``)
package ko

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderInternalKo renders a tangl:내부 trigger head back to "`<event>` 발생 시" (OnEvent)
// or "매 <interval>마다" with an optional ", `<case>` 까지" until clause (EveryTick). The
// canonical interval is denormalized back to its KO 초/분/시간 form; its run/check body
// items render as separate StmtBlocks. It inverts parseInternalHeadKo.
func renderInternalKo(in ast.Internal) (string, error) {
	switch in.Kind {
	case ast.OnEvent:
		return "`" + in.Event + "` 발생 시", nil
	case ast.EveryTick:
		s := "매 " + denormalizeIntervalKo(in.Interval) + "마다"
		if in.Until != "" {
			s += ", `" + in.Until + "` 까지"
		}
		return s, nil
	default:
		return "", fmt.Errorf("ko: unknown InternalKind %d", int(in.Kind))
	}
}
