//ff:func feature=approve type=rule control=sequence
//ff:what IsUrgent: 긴급 승인 요청인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsUrgent checks if the request is marked as urgent.
func IsUrgent(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	reqMeta, _ := ctx.Get("requestMetadata")
	m, ok := reqMeta.(map[string]any)
	if !ok {
		return false, nil
	}
	urgent, _ := m["urgent"].(bool)
	return urgent, nil
}
