//ff:func feature=tangl type=model control=sequence
//ff:what TestCall_OperandNode — tests Call.operandNode marker method and Operand interface satisfaction
package ast

import "testing"

func TestCall_OperandNode(t *testing.T) {
	c := Call{}
	c.operandNode()

	var o Operand = Call{}
	o.operandNode()
}
