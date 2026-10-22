//ff:func feature=tangl type=util control=sequence
//ff:what TestLookupCodec — a registered locale is found; an unregistered one is not
package doc

import "testing"

func TestLookupCodec(t *testing.T) {
	RegisterCodec(stubCodec{locale: "zz-lookup"})
	if c, ok := LookupCodec("zz-lookup"); !ok || c.Locale() != "zz-lookup" {
		t.Fatalf("registered: %v %v", c, ok)
	}
	if _, ok := LookupCodec("no-such-locale"); ok {
		t.Fatal("unregistered: want not found")
	}
}
