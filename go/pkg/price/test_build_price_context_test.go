//ff:func feature=price type=engine control=iteration dimension=1
//ff:what TestBuildPriceContext — verifies buildPriceContext maps all PurchaseRequest and PriceContext fields into toulmin.Context
package price

import "testing"

func TestBuildPriceContext(t *testing.T) {
	req := &PurchaseRequest{
		ProductID: "p-1",
		Quantity:  3,
		BasePrice: 9.99,
		Metadata:  map[string]any{"reqKey": "reqVal"},
	}
	pc := &PriceContext{
		User:       "alice",
		Membership: "gold",
		Coupons:    []Coupon{{Code: "SAVE10"}},
		Promotions: []Promotion{{Name: "summer"}},
		Metadata:   map[string]any{"pcKey": "pcVal"},
	}

	ctx := buildPriceContext(req, pc)

	cases := []struct {
		key  string
		want any
	}{
		{"productID", req.ProductID},
		{"quantity", req.Quantity},
		{"basePrice", req.BasePrice},
		{"user", pc.User},
		{"membership", pc.Membership},
	}
	for _, c := range cases {
		v, ok := ctx.Get(c.key)
		if !ok {
			t.Fatalf("key %q not set", c.key)
		}
		if v != c.want {
			t.Fatalf("key %q = %v, want %v", c.key, v, c.want)
		}
	}

	rm, ok := ctx.Get("requestMetadata")
	if !ok {
		t.Fatal("requestMetadata not set")
	}
	if m, ok := rm.(map[string]any); !ok || m["reqKey"] != "reqVal" {
		t.Fatalf("requestMetadata not set correctly: %v", rm)
	}

	coupons, ok := ctx.Get("coupons")
	if !ok {
		t.Fatal("coupons not set")
	}
	if cs, ok := coupons.([]Coupon); !ok || len(cs) != 1 || cs[0].Code != "SAVE10" {
		t.Fatalf("coupons not set correctly: %v", coupons)
	}

	promos, ok := ctx.Get("promotions")
	if !ok {
		t.Fatal("promotions not set")
	}
	if ps, ok := promos.([]Promotion); !ok || len(ps) != 1 || ps[0].Name != "summer" {
		t.Fatalf("promotions not set correctly: %v", promos)
	}

	md, ok := ctx.Get("metadata")
	if !ok {
		t.Fatal("metadata not set")
	}
	if m, ok := md.(map[string]any); !ok || m["pcKey"] != "pcVal" {
		t.Fatalf("metadata not set correctly: %v", md)
	}
}
