//ff:func feature=feature type=rule control=sequence
//ff:what IsLegacyBrowser: 레거시 브라우저인지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsLegacyBrowser returns true if the user has the "legacy_browser" attribute.
func IsLegacyBrowser(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*UserContext)
	legacy, _ := ctx.Attributes["legacy_browser"].(bool)
	return legacy, nil
}
