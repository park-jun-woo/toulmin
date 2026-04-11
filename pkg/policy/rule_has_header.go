//ff:func feature=policy type=rule control=sequence
//ff:what HasHeader: spec(*HeaderSpec)으로 지정된 헤더가 존재하는지 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasHeader checks if the request has a non-empty header specified by spec.
func HasHeader(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	headers, _ := ctx.Get("headers")
	if len(specs) == 0 {
		return false, nil
	}
	b := specs[0].(*HeaderSpec)
	return headers.(map[string]string)[b.Header] != "", nil
}
