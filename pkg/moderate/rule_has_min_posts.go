//ff:func feature=moderate type=rule control=sequence
//ff:what HasMinPosts: spec(MinPostsSpec)으로 지정된 최소 게시 수 이상인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasMinPosts checks if the author has at least the minimum post count.
// spec is *MinPostsSpec.
func HasMinPosts(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	author, _ := ctx.Get("author")
	if len(specs) == 0 {
		return false, nil
	}
	mp := specs[0].(*MinPostsSpec)
	return author.(*Author).PostCount >= mp.MinPosts, nil
}
