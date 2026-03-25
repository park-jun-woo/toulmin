//ff:func feature=state type=adapter control=sequence
//ff:what buildStateContext: TransitionRequest + TransitionContext → toulmin.Context 변환
package state

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// buildStateContext converts a TransitionRequest and TransitionContext into a toulmin.Context.
func buildStateContext(req *TransitionRequest, tc *TransitionContext) toulmin.Context {
	ctx := toulmin.NewContext()
	ctx.Set("from", req.From)
	ctx.Set("to", req.To)
	ctx.Set("event", req.Event)
	ctx.Set("currentState", tc.CurrentState)
	ctx.Set("user", tc.User)
	ctx.Set("resource", tc.Resource)
	ctx.Set("metadata", tc.Metadata)
	ctx.Set("userID", tc.UserID)
	ctx.Set("resourceOwnerID", tc.ResourceOwnerID)
	return ctx
}
