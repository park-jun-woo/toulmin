//ff:func feature=moderate type=rule control=sequence
//ff:what IsNewsContext: 뉴스 채널인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsNewsContext returns true if the channel type is "news".
func IsNewsContext(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*ContentContext)
	return ctx.Channel.Type == "news", nil
}
