//ff:func feature=policy type=rule control=sequence
//ff:what HasHeader: backing(string)으로 지정된 헤더가 존재하는지 판정
package policy

// HasHeader checks if the request has a non-empty header specified by backing (string).
func HasHeader(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RequestContext)
	header := backing.(string)
	return ctx.Headers[header] != "", nil
}
