//ff:func feature=state type=rule control=sequence
//ff:what IsExpired: backing(func)의 만료 시간 추출로 만료 판정
package state

import "time"

// IsExpired checks if the resource is expired using backing (func(any) time.Time).
func IsExpired(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*TransitionContext)
	expiryFunc := backing.(func(any) time.Time)
	return time.Now().After(expiryFunc(ctx.Resource)), nil
}
