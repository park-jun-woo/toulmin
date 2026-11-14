//ff:func feature=tangl type=engine control=sequence
//ff:what testCtx — a Context seeded with x="1" and y="9" for compile tests
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/toulmin"

// testCtx returns a Context seeded so that `x` equals "1" (truthy) and `y`
// equals "9" (falsy against "2").
func testCtx() *toulmin.MapContext {
	c := toulmin.NewContext()
	c.Set("x", "1")
	c.Set("y", "9")
	return c
}
