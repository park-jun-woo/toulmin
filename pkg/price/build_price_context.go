//ff:func feature=price type=adapter control=sequence
//ff:what buildPriceContext: PurchaseRequest + PriceContext → toulmin.Context 변환
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// buildPriceContext converts a PurchaseRequest and PriceContext into a toulmin.Context.
func buildPriceContext(req *PurchaseRequest, pc *PriceContext) toulmin.Context {
	ctx := toulmin.NewContext()
	ctx.Set("productID", req.ProductID)
	ctx.Set("quantity", req.Quantity)
	ctx.Set("basePrice", req.BasePrice)
	ctx.Set("requestMetadata", req.Metadata)
	ctx.Set("user", pc.User)
	ctx.Set("membership", pc.Membership)
	ctx.Set("coupons", pc.Coupons)
	ctx.Set("promotions", pc.Promotions)
	ctx.Set("metadata", pc.Metadata)
	return ctx
}
