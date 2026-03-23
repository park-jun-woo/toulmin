//ff:func feature=feature type=rule control=sequence
//ff:what IsRegion: backing(RegionBacking)으로 지정된 지역인지 판정
package feature

// IsRegion checks if the user's region matches backing.Region.
func IsRegion(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*UserContext)
	rb := backing.(*RegionBacking)
	return ctx.Region == rb.Region, nil
}
