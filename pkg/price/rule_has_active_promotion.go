//ff:func feature=price type=rule control=iteration dimension=1
//ff:what HasActivePromotion: spec(DiscountSpec).Name의 프로모션이 활성인지 판정
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasActivePromotion checks if the promotion named by spec.Name is active.
func HasActivePromotion(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	promotions, _ := ctx.Get("promotions")
	if len(specs) == 0 {
		return false, nil
	}
	db := specs[0].(*DiscountSpec)
	for _, p := range promotions.([]Promotion) {
		if p.Name == db.Name && p.Active {
			return true, db
		}
	}
	return false, nil
}
