//ff:func feature=tangl type=util control=sequence
//ff:what issueID — hand out the next free Edge Id for a non-primary record
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
)

// issueID returns the next free Edge Id for a non-primary record and advances
// the counter. It is the sole issuer after pass 1, so ids never collide.
func issueID(e *enc) toulmin.TID {
	id := e.next
	e.next++
	return id
}
