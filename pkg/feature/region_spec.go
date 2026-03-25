//ff:type feature=feature type=model
//ff:what RegionSpec: IsRegion rule의 spec 타입
package feature

// RegionSpec carries region criteria for feature flag checks.
type RegionSpec struct {
	Region string // target region code ("KR", "US", etc.)
}
