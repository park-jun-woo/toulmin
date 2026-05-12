//ff:func feature=moderate type=rule control=sequence
//ff:what IsNewsContext: 뉴스 채널인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsNewsContext returns true if the channel type is "news".
func IsNewsContext(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	channel, _ := ctx.Get("channel")
	ch, ok := channel.(*Channel)
	if !ok {
		return false, nil
	}
	return ch.Type == "news", nil
}
