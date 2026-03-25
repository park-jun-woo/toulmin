//ff:func feature=price type=engine control=iteration dimension=1
//ff:what Evaluate: 할인 판정 + 합산 + cap 적용 + 최종 가격
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Evaluate runs the discount graph and computes the final price.
func (p *Pricer) Evaluate(req *PurchaseRequest, pctx *PriceContext) (*PriceResult, error) {
	ctx := buildPriceContext(req, pctx)
	results, err := p.graph.Evaluate(ctx, toulmin.EvalOption{Trace: true})
	if err != nil {
		return nil, err
	}
	var applied []*DiscountSpec
	var allTrace []toulmin.TraceEntry
	totalDiscount := 0.0

	for _, r := range results {
		if len(r.Trace) > 0 {
			allTrace = append(allTrace, r.Trace...)
		}
		db, ok := r.Evidence.(*DiscountSpec)
		if !ok || db == nil {
			continue
		}
		discount := calcDiscount(req.BasePrice, db)
		totalDiscount += discount
		applied = append(applied, db)
	}

	if p.totalCap != nil && p.totalCap.Max > 0 && totalDiscount > p.totalCap.Max {
		totalDiscount = p.totalCap.Max
	}

	if totalDiscount > req.BasePrice {
		totalDiscount = req.BasePrice
	}

	return &PriceResult{
		BasePrice:        req.BasePrice,
		TotalDiscount:    totalDiscount,
		FinalPrice:       req.BasePrice - totalDiscount,
		AppliedDiscounts: applied,
		Trace:            allTrace,
	}, nil
}
