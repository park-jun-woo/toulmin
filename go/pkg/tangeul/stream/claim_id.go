//ff:func feature=tangl type=util control=sequence
//ff:what claimID — register a record's Edge Id, rejecting 0 and duplicates
package stream

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// claimID records id as used. Edge Ids start at 1 (0 means "no reference") and
// are unique within a document (GEUL t002 TOULMIN 공통 필드), so 0 or a repeat is
// an error.
func claimID(d *dec, id toulmin.TID) error {
	if id == 0 || d.seen[id] {
		return fmt.Errorf("tangeul/stream: Edge Id %d is zero or used twice", id)
	}
	d.seen[id] = true
	return nil
}
