//ff:func feature=tangl type=model control=sequence
//ff:what TestPath_OperandNode — tests Path.operandNode marker method and Operand interface satisfaction
package ast

import "testing"

func TestPath_OperandNode(t *testing.T) {
	p := Path{}
	p.operandNode()

	var o Operand = Path{}
	o.operandNode()
}
