//ff:func feature=state type=rule control=sequence
//ff:what IsExpired: spec(ExpirySpec)의 만료 시간으로 만료 판정
package state

import (
	"time"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// IsExpired checks if the resource is expired using spec (*ExpirySpec).
func IsExpired(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	if len(specs) == 0 {
		return false, nil
	}
	b := specs[0].(*ExpirySpec)
	return time.Now().After(b.ExpiresAt), nil
}
