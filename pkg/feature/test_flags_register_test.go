//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_Register — tests registering new and re-registering existing feature names
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlags_Register(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		flags := NewFlags()

		g := toulmin.NewGraph("feature:a")
		flags.Register("a", g)

		if len(flags.order) != 1 || flags.order[0] != "a" {
			t.Fatalf("expected order [a], got %v", flags.order)
		}
		if flags.features["a"] != g {
			t.Fatalf("expected registered graph to match")
		}
	})

	t.Run("Overwrite", func(t *testing.T) {
		flags := NewFlags()

		g1 := toulmin.NewGraph("feature:a-v1")
		flags.Register("a", g1)

		g2 := toulmin.NewGraph("feature:a-v2")
		flags.Register("a", g2)

		if len(flags.order) != 1 || flags.order[0] != "a" {
			t.Fatalf("expected order to remain [a] without duplicate, got %v", flags.order)
		}
		if flags.features["a"] != g2 {
			t.Fatalf("expected re-registration to overwrite graph")
		}
	})
}
