//ff:type feature=moderate type=model
//ff:what MinPostsBacking: HasMinPosts rule의 backing 타입
package moderate

// MinPostsBacking carries minimum post count criteria.
type MinPostsBacking struct {
	MinPosts int // minimum number of posts required
}
