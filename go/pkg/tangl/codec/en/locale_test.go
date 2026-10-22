//ff:func feature=tangl type=parser control=sequence
//ff:what TestLocale — the codec reports the "en" locale tag
package en

import "testing"

// TestLocale pins the codec's locale tag to "en".
func TestLocale(t *testing.T) {
	if got := (ENCodec{}).Locale(); got != "en" {
		t.Fatalf("Locale() = %q, want en", got)
	}
}
