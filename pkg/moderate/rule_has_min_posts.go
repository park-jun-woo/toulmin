//ff:func feature=moderate type=rule control=sequence
//ff:what HasMinPosts: backing(MinPostsBacking)으로 지정된 최소 게시 수 이상인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasMinPosts checks if the author has at least the minimum post count.
// backing is *MinPostsBacking.
func HasMinPosts(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	author, _ := ctx.Get("author")
	mp := backing.(*MinPostsBacking)
	return author.(*Author).PostCount >= mp.MinPosts, nil
}
