//ff:func feature=tangl type=model control=sequence
//ff:what TestSpecRef_OperandNode — tests SpecRef.operandNode marker method and Operand interface satisfaction
package ast

import "testing"

func TestSpecRef_OperandNode(t *testing.T) {
	s := SpecRef{}
	s.operandNode()

	var o Operand = SpecRef{}
	o.operandNode()
}
