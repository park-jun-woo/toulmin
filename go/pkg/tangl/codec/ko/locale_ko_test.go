//ff:func feature=tangl type=parser control=sequence
//ff:what TestLocaleKo — the codec reports the "ko" locale tag
package ko

import "testing"

func TestLocaleKo(t *testing.T) {
	if got := (KoCodec{}).Locale(); got != "ko" {
		t.Fatalf("Locale() = %q, want ko", got)
	}
}
