//ff:func feature=engine type=engine control=sequence
//ff:what authenticate — test helper: warrant active when context "authenticated" is true
package toulmin

func authenticate(ctx Context, specs Specs) (bool, any) {
	v, _ := ctx.Get("authenticated")
	ok, _ := v.(bool)
	return ok, "authenticated"
}
