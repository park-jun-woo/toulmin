//ff:type feature=price type=model
//ff:what PriceResult: 가격 판정 결과
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// PriceResult holds the price evaluation result.
type PriceResult struct {
	BasePrice        float64
	TotalDiscount    float64
	FinalPrice       float64
	AppliedDiscounts []*DiscountSpec
	Trace            []toulmin.TraceEntry
}
