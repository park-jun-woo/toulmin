//ff:func feature=approve type=rule control=sequence
//ff:what IsAboveLevel: spec(ApproverSpec)으로 지정된 최소 레벨 이상인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsAboveLevel checks if the approver's level is at or above spec.Level.
func IsAboveLevel(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	approverLevel, _ := ctx.Get("approverLevel")
	if len(specs) == 0 {
		return false, nil
	}
	ab := specs[0].(*ApproverSpec)
	return approverLevel.(int) >= ab.Level, nil
}
