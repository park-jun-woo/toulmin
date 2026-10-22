//ff:type feature=moderate type=model
//ff:what MinPostsSpec: HasMinPosts rule의 spec 타입
package moderate

// MinPostsSpec carries minimum post count criteria.
type MinPostsSpec struct {
	MinPosts int // minimum number of posts required
}
