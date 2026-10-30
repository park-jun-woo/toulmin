//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestSafeCall — tests safeCall for normal return, partial verdict, panic recovery, and out-of-range/NaN verdict rejection
package toulmin

import (
	"math"
	"strings"
	"testing"
)

func TestSafeCall(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Normal", func(t *testing.T) {
			fn := func(ctx Context, specs Specs) (float64, any) { return 1.0, "evidence" }
			v, evidence, err := safeCall(fn, NewContext(), nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v != 1.0 {
				t.Errorf("expected verdict=1.0, got %v", v)
			}
			if evidence != "evidence" {
				t.Errorf("expected evidence=%q, got %v", "evidence", evidence)
			}
		}},
		{"BoundariesAndPartial", func(t *testing.T) {
			for _, want := range []float64{-1.0, 0.0, 0.25, 1.0} {
				w := want
				fn := func(ctx Context, specs Specs) (float64, any) { return w, w }
				v, evidence, err := safeCall(fn, NewContext(), nil)
				if err != nil || v != w || evidence != w {
					t.Errorf("verdict %v: got (%v, %v, %v)", w, v, evidence, err)
				}
			}
		}},
		{"Panic", func(t *testing.T) {
			fn := func(ctx Context, specs Specs) (float64, any) {
				panic("boom")
			}
			v, evidence, err := safeCall(fn, NewContext(), nil)
			if err == nil {
				t.Fatal("expected error from recovered panic")
			}
			if !strings.Contains(err.Error(), "rule panicked") || !strings.Contains(err.Error(), "boom") {
				t.Errorf("unexpected error message: %v", err)
			}
			if v != -1.0 {
				t.Errorf("expected verdict=-1 on panic, got %v", v)
			}
			if evidence != nil {
				t.Errorf("expected evidence=nil on panic, got %v", evidence)
			}
		}},
		{"OutOfRangeAndNaN", func(t *testing.T) {
			for _, bad := range []float64{1.0000001, -1.0000001, 2, -5, math.NaN(), math.Inf(1), math.Inf(-1)} {
				b := bad
				fn := func(ctx Context, specs Specs) (float64, any) { return b, "leaked" }
				v, evidence, err := safeCall(fn, NewContext(), nil)
				if err == nil {
					t.Fatalf("verdict %v: expected out-of-range error", b)
				}
				if !strings.Contains(err.Error(), "out of range [-1, 1]") {
					t.Errorf("verdict %v: unexpected error message: %v", b, err)
				}
				if v != -1.0 || evidence != nil {
					t.Errorf("verdict %v: expected (-1, nil), got (%v, %v)", b, v, evidence)
				}
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
