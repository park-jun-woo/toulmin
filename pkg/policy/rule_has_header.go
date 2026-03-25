//ff:func feature=policy type=rule control=sequence
//ff:what HasHeader: backing(*HeaderBacking)으로 지정된 헤더가 존재하는지 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasHeader checks if the request has a non-empty header specified by backing.
func HasHeader(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	headers, _ := ctx.Get("headers")
	b := backing.(*HeaderBacking)
	return headers.(map[string]string)[b.Header] != "", nil
}
