//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what IsIPInList: spec(IPListSpec)으로 전달된 IP 목록에 클라이언트 IP가 있는지 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsIPInList checks if the client IP is in the list provided by spec.
func IsIPInList(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	clientIP, _ := ctx.Get("clientIP")
	b := specs[0].(*IPListSpec)
	for _, ip := range b.List {
		if ip == clientIP.(string) {
			return true, nil
		}
	}
	return false, nil
}
