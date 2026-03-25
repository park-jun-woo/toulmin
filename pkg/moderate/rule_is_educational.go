//ff:func feature=moderate type=rule control=sequence
//ff:what IsEducational: 교육 채널인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsEducational returns true if the channel type is "education".
func IsEducational(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*ContentContext)
	return ctx.Channel.Type == "education", nil
}
