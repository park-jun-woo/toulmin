//ff:func feature=state type=rule control=sequence
//ff:what isAuth — test helper auth function stub
package state

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func isAuth(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
