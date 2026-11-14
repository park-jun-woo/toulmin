//ff:type feature=state type=model
//ff:what ExpirySpec: IsExpired rule의 spec 타입 (만료 판정 기준)
package state

import "time"

// ExpirySpec carries the expiry time for expiration checks.
type ExpirySpec struct {
	ExpiresAt time.Time
}
