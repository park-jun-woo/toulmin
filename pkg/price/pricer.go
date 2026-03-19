//ff:type feature=price type=engine
//ff:what Pricer: 할인 판정 + 최종 가격 계산
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Pricer evaluates discount rules and computes the final price.
type Pricer struct {
	graph    *toulmin.Graph
	totalCap *DiscountBacking
}
