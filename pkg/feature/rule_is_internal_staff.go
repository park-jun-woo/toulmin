//ff:func feature=feature type=rule control=sequence
//ff:what IsInternalStaff: 내부 직원인지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsInternalStaff checks if the user is internal staff.
// Checks Attributes["internal"].
func IsInternalStaff(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	attributes, _ := ctx.Get("attributes")
	internal, _ := attributes.(map[string]any)["internal"].(bool)
	return internal, nil
}
