//ff:func feature=feature type=rule control=sequence
//ff:what IsRegion: backing(string)으로 지정된 지역인지 판정
package feature

// IsRegion checks if the user's region matches backing (string).
func IsRegion(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*UserContext)
	region := backing.(string)
	return ctx.Region == region, nil
}
