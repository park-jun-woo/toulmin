//ff:type feature=feature type=model
//ff:what RegionBacking: IsRegion rule의 backing 타입
package feature

// RegionBacking carries region criteria for feature flag checks.
type RegionBacking struct {
	Region string // target region code ("KR", "US", etc.)
}
