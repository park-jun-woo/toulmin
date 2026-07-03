//ff:func feature=tangl type=util control=sequence
//ff:what TestRegisterCodec — RegisterCodec + LookupCodec round-trip; duplicate locale panics
package doc

import "testing"

// TestRegisterCodec verifies a registered codec is found by locale, an unknown
// locale is not found, and registering the same locale twice panics.
func TestRegisterCodec(t *testing.T) {
	c := stubCodec{locale: "regtest"}
	RegisterCodec(c)

	got, ok := LookupCodec("regtest")
	if !ok || got.Locale() != "regtest" {
		t.Fatalf("LookupCodec(regtest) = %v, %v", got, ok)
	}
	if _, ok := LookupCodec("nope"); ok {
		t.Fatalf("LookupCodec(nope) = found, want missing")
	}

	defer func() {
		if recover() == nil {
			t.Errorf("duplicate RegisterCodec: expected panic")
		}
	}()
	RegisterCodec(c)
}
