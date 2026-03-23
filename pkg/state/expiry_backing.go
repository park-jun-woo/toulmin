//ff:type feature=state type=model
//ff:what ExpiryBacking: IsExpired rule의 backing 타입 (만료 판정 기준)
package state

import "time"

// ExpiryBacking carries the expiry time for expiration checks.
type ExpiryBacking struct {
	ExpiresAt time.Time
}
