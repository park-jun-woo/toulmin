//ff:func feature=approve type=rule control=sequence
//ff:what IsUrgent: 긴급 승인 요청인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsUrgent checks if the request is marked as urgent.
func IsUrgent(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	reqMeta, _ := ctx.Get("requestMetadata")
	urgent, _ := reqMeta.(map[string]any)["urgent"].(bool)
	return urgent, nil
}
