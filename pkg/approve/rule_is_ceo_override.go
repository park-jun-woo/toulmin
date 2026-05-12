//ff:func feature=approve type=rule control=sequence
//ff:what IsCEOOverride: CEO 직권 승인 여부 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsCEOOverride checks if the approver is the CEO.
func IsCEOOverride(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	approverRole, _ := ctx.Get("approverRole")
	s, ok := approverRole.(string)
	if !ok {
		return false, nil
	}
	return s == "ceo", nil
}
