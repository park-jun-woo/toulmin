//ff:func feature=moderate type=rule control=sequence
//ff:what IsNewsContext: 뉴스 채널인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsNewsContext returns true if the channel type is "news".
func IsNewsContext(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	channel, _ := ctx.Get("channel")
	return channel.(*Channel).Type == "news", nil
}
