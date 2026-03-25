//ff:func feature=moderate type=rule control=sequence
//ff:what IsAdultChannel: 성인 채널인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsAdultChannel returns true if the channel is age-gated.
func IsAdultChannel(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*ContentContext)
	return ctx.Channel.AgeGated, nil
}
