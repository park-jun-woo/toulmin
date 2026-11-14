//ff:func feature=tangl type=model control=sequence
//ff:what operandNode — marks Call as an Operand implementation (value position)
package ast

// operandNode marks Call as an Operand implementation (value position).
func (c Call) operandNode() {}
