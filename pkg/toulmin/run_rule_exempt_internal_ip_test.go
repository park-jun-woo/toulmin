//ff:func feature=engine type=engine control=sequence
//ff:what exemptInternalIP — test helper: except active when context "ip_internal" is true
package toulmin

func exemptInternalIP(ctx Context, specs Specs) (bool, any) {
	v, _ := ctx.Get("ip_internal")
	ok, _ := v.(bool)
	return ok, "ip_internal"
}
