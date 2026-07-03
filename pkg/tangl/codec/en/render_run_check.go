//ff:func feature=tangl type=parser control=selection
//ff:what renderRunCheck — render a Provides/Internal body run/check line (`run `C`` / `check `C``)
package en

import "fmt"

// renderRunCheck renders a run/check body item back to "run `<case>`" (slotKind 1)
// or "check `<case>`" (slotKind 2). Runs/Checks are Line-less []string, so the
// StmtRef's slotKind — not the resolved value — selects the keyword.
func renderRunCheck(slotKind int, name string) (string, error) {
	switch slotKind {
	case 1:
		return "run `" + name + "`", nil
	case 2:
		return "check `" + name + "`", nil
	default:
		return "", fmt.Errorf("en: unknown run/check slotKind %d", slotKind)
	}
}
