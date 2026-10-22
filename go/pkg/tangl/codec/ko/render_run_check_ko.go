//ff:func feature=tangl type=parser control=selection
//ff:what renderRunCheckKo — render a 업무/내부 body run/check line (`` `C` 진행한다 `` / `` `C` 판정한다 ``)
package ko

import "fmt"

// renderRunCheckKo renders a run/check body item back to "`<case>` 진행한다" (slotKind 1)
// or "`<case>` 판정한다" (slotKind 2). Runs/Checks are Line-less []string, so the
// StmtRef's slotKind — not the resolved value — selects the keyword. Unlike the
// case-scoped run edge, this body form has no 면 condition node.
func renderRunCheckKo(slotKind int, name string) (string, error) {
	switch slotKind {
	case 1:
		return "`" + name + "` 진행한다", nil
	case 2:
		return "`" + name + "` 판정한다", nil
	default:
		return "", fmt.Errorf("ko: unknown run/check slotKind %d", slotKind)
	}
}
