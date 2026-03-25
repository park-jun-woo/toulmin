//ff:func feature=price type=rule control=sequence
//ff:what IsAlreadyDiscounted: 이미 할인이 적용되었는지 판정
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsAlreadyDiscounted checks if the purchase is already discounted.
func IsAlreadyDiscounted(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	reqMeta, _ := ctx.Get("requestMetadata")
	discounted, _ := reqMeta.(map[string]any)["discounted"].(bool)
	return discounted, nil
}
