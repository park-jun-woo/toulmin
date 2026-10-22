//ff:type feature=price type=model
//ff:what PurchaseRequest: 구매 요청
package price

// PurchaseRequest represents a purchase to be evaluated for discounts.
type PurchaseRequest struct {
	ProductID string
	Quantity  int
	BasePrice float64
	Metadata  map[string]any
}
