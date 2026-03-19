//ff:func feature=moderate type=rule control=sequence
//ff:what IsAdultChannel: 성인 채널인지 판정
package moderate

// IsAdultChannel returns true if the channel is age-gated.
func IsAdultChannel(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ContentContext)
	return ctx.Channel.AgeGated, nil
}
