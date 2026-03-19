//ff:func feature=moderate type=rule control=sequence
//ff:what IsEducational: 교육 채널인지 판정
package moderate

// IsEducational returns true if the channel type is "education".
func IsEducational(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ContentContext)
	return ctx.Channel.Type == "education", nil
}
