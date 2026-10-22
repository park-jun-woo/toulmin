//ff:type feature=tangl type=interface
//ff:what Operand — a value-producing term inside a Rules condition (Path, SpecRef, Call)
package ast

// Operand is a value-producing term inside a Rules condition. Concrete
// implementations are Path, SpecRef, and Call.
type Operand interface{ operandNode() }
