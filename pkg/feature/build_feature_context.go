//ff:func feature=feature type=adapter control=sequence
//ff:what buildFeatureContext: feature name + UserContext → toulmin.Context 변환
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// buildFeatureContext converts a feature name and UserContext into a toulmin.Context.
func buildFeatureContext(name string, uc *UserContext) toulmin.Context {
	ctx := toulmin.NewContext()
	ctx.Set("featureName", name)
	ctx.Set("user", uc.User)
	ctx.Set("id", uc.ID)
	ctx.Set("region", uc.Region)
	ctx.Set("attributes", uc.Attributes)
	return ctx
}
