//ff:func feature=tangl type=codegen control=sequence
//ff:what equivCanon — canonical string form of a (bool, evidence) result
package gen

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// equivCanon renders a rule result as a canonical comparison string carrying
// the verdict, the witness Fact fields (with the Actual's dynamic type), or
// the passed-through evidence. The generated runner defines an identical copy.
func equivCanon(b bool, ev any) string {
	if f, ok := ev.(toulmin.Fact); ok {
		return fmt.Sprintf("%t|fact|%v|%v|%T|%v", b, f.Where, f.Expected, f.Actual, f.Actual)
	}
	return fmt.Sprintf("%t|ev|%T|%v", b, ev, ev)
}
