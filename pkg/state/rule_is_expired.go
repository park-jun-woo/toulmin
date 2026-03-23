//ff:func feature=state type=rule control=sequence
//ff:what IsExpired: backing(ExpiryBacking)의 만료 시간으로 만료 판정
package state

import "time"

// IsExpired checks if the resource is expired using backing (*ExpiryBacking).
func IsExpired(claim any, ground any, backing any) (bool, any) {
	b := backing.(*ExpiryBacking)
	return time.Now().After(b.ExpiresAt), nil
}
