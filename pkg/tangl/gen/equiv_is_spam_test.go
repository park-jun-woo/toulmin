//ff:func feature=tangl type=codegen control=sequence
//ff:what equivIsSpam — cross-test host predicate function (isSpam holds)
package gen

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// equivIsSpam is the cross-test host predicate bound to the bare
// `isSpam holds` call: true when ctx spamFlag is true, returning the ctx
// spamEvidence value (nil or a host evidence to pass through). The generated
// runner defines an identical copy named isSpam.
func equivIsSpam(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	v, _ := ctx.Get("spamFlag")
	ev, _ := ctx.Get("spamEvidence")
	return v == true, ev
}
