//ff:func feature=tangl type=util control=sequence
//ff:what stepRead — advance the reader by one control or one statement, returning the next index
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// stepRead processes the packet at packets[i]: a control is applied via
// handleControl (advancing one packet), otherwise a whole statement is collected
// and interpreted. A statement encountered while a SECTION heading is still
// pending is a hard error. It returns the index to continue from.
func stepRead(d *dec, packets []rawPacket, i int, pendingSection *int) (int, error) {
	if packets[i].Class == word.ClassControl {
		if err := handleControl(d, pendingSection, packets[i]); err != nil {
			return 0, err
		}
		return i + 1, nil
	}
	if *pendingSection >= 0 {
		return 0, errors.New("tangeul/stream: SECTION control must be followed by its heading PROSE")
	}
	return readStatement(d, packets, i)
}
