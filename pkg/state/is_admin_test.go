//ff:func feature=state type=rule control=sequence
//ff:what isAdmin — test helper admin function stub
package state

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func isAdmin(claim any, ground any, backing toulmin.Backing) (bool, any) { return true, nil }
