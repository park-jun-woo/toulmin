//ff:func feature=feature type=rule control=sequence
//ff:what IsBetaUser: 베타 사용자인지 판정
package feature

// IsBetaUser returns true if the user has the "beta" attribute set to true.
func IsBetaUser(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*UserContext)
	beta, _ := ctx.Attributes["beta"].(bool)
	return beta, nil
}
