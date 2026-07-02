//ff:func feature=price type=engine control=sequence
//ff:what TestNewPricer — verifies NewPricer constructs a Pricer with the given graph and total cap
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestNewPricer(t *testing.T) {
	t.Run("basic construction", func(t *testing.T) {
		g := toulmin.NewGraph("test:new-pricer")
		cap := &DiscountSpec{Rate: 0.5}

		p := NewPricer(g, cap)
		if p == nil {
			t.Fatal("expected non-nil Pricer")
		}
		if p.graph != g {
			t.Errorf("graph = %v, want %v", p.graph, g)
		}
		if p.totalCap != cap {
			t.Errorf("totalCap = %v, want %v", p.totalCap, cap)
		}
	})

	t.Run("nil total cap", func(t *testing.T) {
		g := toulmin.NewGraph("test:new-pricer-nil")

		p := NewPricer(g, nil)
		if p == nil {
			t.Fatal("expected non-nil Pricer")
		}
		if p.totalCap != nil {
			t.Errorf("totalCap = %v, want nil", p.totalCap)
		}
	})
}
