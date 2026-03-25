//ff:func feature=feature type=rule control=sequence
//ff:what IsRegion: backing(RegionBacking)으로 지정된 지역인지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsRegion checks if the user's region matches backing.Region.
func IsRegion(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	region, _ := ctx.Get("region")
	rb := backing.(*RegionBacking)
	return region.(string) == rb.Region, nil
}
