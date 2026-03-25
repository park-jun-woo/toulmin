//ff:func feature=price type=rule control=iteration dimension=1
//ff:what HasActivePromotion: backing(DiscountBacking).Name의 프로모션이 활성인지 판정
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasActivePromotion checks if the promotion named by backing.Name is active.
func HasActivePromotion(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	promotions, _ := ctx.Get("promotions")
	db := backing.(*DiscountBacking)
	for _, p := range promotions.([]Promotion) {
		if p.Name == db.Name && p.Active {
			return true, db
		}
	}
	return false, nil
}
