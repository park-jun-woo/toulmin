//ff:func feature=moderate type=rule control=sequence
//ff:what IsAdultChannel: 성인 채널인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsAdultChannel returns true if the channel is age-gated.
func IsAdultChannel(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	channel, _ := ctx.Get("channel")
	return channel.(*Channel).AgeGated, nil
}
