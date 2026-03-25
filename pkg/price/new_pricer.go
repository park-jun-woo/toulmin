//ff:func feature=price type=engine control=sequence
//ff:what NewPricer: Pricer 생성
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// NewPricer creates a Pricer. totalCap is optional (nil for no cap).
func NewPricer(g *toulmin.Graph, totalCap *DiscountSpec) *Pricer {
	return &Pricer{graph: g, totalCap: totalCap}
}
