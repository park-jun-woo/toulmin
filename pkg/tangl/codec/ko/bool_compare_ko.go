//ff:func feature=tangl type=parser control=selection
//ff:what boolCompareKo — build a bare-boolean Compare{Op:equals, Value:true} from path steps
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// boolCompareKo builds the EN-canonical bare-boolean comparison a bare name maps to:
// Compare{Op:"equals", Value:"true"}. A single-step path is demoted to Field (Lhs nil,
// preserving the byte-identical AST of pre-existing `x` 이면 documents); a multi-step
// 의-chain path rides in Lhs so `작성자` 의 `활성` 이면 matches EN `author`.`active`
// equals true. It is shared by parseCondTermKo (이면/아니고) and parseCompareAdnKo
// (인/이 아닌).
func boolCompareKo(steps []string) ast.Compare {
	cmp := ast.Compare{Op: "equals", Value: "true"}
	switch len(steps) {
	case 1:
		cmp.Field = steps[0]
	default:
		cmp.Lhs = ast.Path{Steps: steps}
	}
	return cmp
}
