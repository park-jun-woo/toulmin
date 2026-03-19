//ff:func feature=moderate type=rule control=sequence
//ff:what HasMinPosts: backing(int)으로 지정된 최소 게시 수 이상인지 판정
package moderate

// HasMinPosts checks if the author has at least the minimum post count.
// backing is int (minimum posts).
func HasMinPosts(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ContentContext)
	min := backing.(int)
	return ctx.Author.PostCount >= min, nil
}
