//ff:func feature=cli type=command control=sequence
//ff:what checkDetectLocale — assert detectLocale returns the expected locale, or an error containing wantErr
package tanglcli

import (
	"strings"
	"testing"
)

// checkDetectLocale runs detectLocale over src and asserts the outcome: when
// wantErr is non-empty it requires an error containing that substring, otherwise
// it requires no error and locale == wantLocale.
func checkDetectLocale(t *testing.T, src, wantLocale, wantErr string) {
	t.Helper()
	got, err := detectLocale([]byte(src))
	if wantErr != "" {
		if err == nil || !strings.Contains(err.Error(), wantErr) {
			t.Fatalf("detectLocale err = %v, want containing %q", err, wantErr)
		}
		return
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != wantLocale {
		t.Errorf("locale = %q, want %q", got, wantLocale)
	}
}
