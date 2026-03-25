//ff:func feature=feature type=rule control=sequence
//ff:what IsRegion: spec(RegionSpec)으로 지정된 지역인지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsRegion checks if the user's region matches spec.Region.
func IsRegion(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	region, _ := ctx.Get("region")
	rb := specs[0].(*RegionSpec)
	return region.(string) == rb.Region, nil
}
