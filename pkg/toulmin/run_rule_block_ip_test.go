//ff:func feature=engine type=engine control=sequence
//ff:what blockIP — test helper: counter active when context "ip_blocked" is true
package toulmin

func blockIP(ctx Context, specs Specs) (bool, any) {
	v, _ := ctx.Get("ip_blocked")
	ok, _ := v.(bool)
	return ok, "ip_blocked"
}
