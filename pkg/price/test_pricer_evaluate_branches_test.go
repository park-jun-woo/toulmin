//ff:func feature=price type=engine control=sequence
//ff:what TestPricerEvaluate_Branches — covers graph error, non-DiscountSpec/nil evidence, total-cap clamp, and base-price clamp branches of Pricer.Evaluate
package price

import "testing"

func TestPricerEvaluate_Branches(t *testing.T) {
	t.Run("graph error", pricerEvaluateGraphErrorCase)
	t.Run("non discount evidence", pricerEvaluateNonDiscountEvidenceCase)
	t.Run("nil discount evidence", pricerEvaluateNilDiscountEvidenceCase)
	t.Run("total cap max not positive", pricerEvaluateTotalCapMaxNotPositiveCase)
	t.Run("total cap not exceeded", pricerEvaluateTotalCapNotExceededCase)
	t.Run("clamped to base price", pricerEvaluateClampedToBasePriceCase)
}
