//ff:func feature=state type=rule control=sequence
//ff:what IsCurrentState: 현재 상태가 전이 요청의 from과 일치하는지 판정
package state

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsCurrentState returns true if the current state matches the transition request's From.
func IsCurrentState(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	from, _ := ctx.Get("from")
	currentState, _ := ctx.Get("currentState")
	return currentState == from, nil
}
