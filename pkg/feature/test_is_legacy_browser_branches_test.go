//ff:func feature=feature type=rule control=sequence
//ff:what TestIsLegacyBrowser_NonMapAttributes — covers non-map attributes branch of IsLegacyBrowser
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsLegacyBrowser_NonMapAttributes(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("attributes", "not-a-map")
	got, _ := IsLegacyBrowser(ctx, nil)
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}
