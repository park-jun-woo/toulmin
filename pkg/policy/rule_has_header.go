//ff:func feature=policy type=rule control=sequence
//ff:what HasHeader: backing(*HeaderBacking)으로 지정된 헤더가 존재하는지 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasHeader checks if the request has a non-empty header specified by backing.
func HasHeader(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*RequestContext)
	b := backing.(*HeaderBacking)
	return ctx.Headers[b.Header] != "", nil
}
