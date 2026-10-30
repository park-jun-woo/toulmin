//ff:func feature=tangl type=model control=sequence
//ff:what addSpecOperand — append a SpecRef operand's first step once (non-SpecRef operands ignored)
package ast

// addSpecOperand appends the first step of op when op is a SpecRef with at
// least one step that has not been seen yet, marking it seen.
func addSpecOperand(op Operand, seen map[string]bool, names []string) []string {
	if ref, ok := op.(SpecRef); ok && len(ref.Steps) > 0 && !seen[ref.Steps[0]] {
		seen[ref.Steps[0]] = true
		names = append(names, ref.Steps[0])
	}
	return names
}
