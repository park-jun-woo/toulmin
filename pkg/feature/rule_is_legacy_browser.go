//ff:func feature=feature type=rule control=sequence
//ff:what IsLegacyBrowser: 레거시 브라우저인지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsLegacyBrowser returns true if the user has the "legacy_browser" attribute.
func IsLegacyBrowser(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	attributes, _ := ctx.Get("attributes")
	legacy, _ := attributes.(map[string]any)["legacy_browser"].(bool)
	return legacy, nil
}
