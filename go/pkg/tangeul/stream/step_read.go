//ff:func feature=tangl type=util control=selection
//ff:what stepRead — consume one control record or one whole statement, returning the next index
package stream

import "errors"

// stepRead processes the record at packets[i]: a framing control is applied via
// handleControl (advancing one record), otherwise a whole statement is collected
// and interpreted. A statement while a SECTION heading is still pending, and a
// second SOURCE, are errors. It returns the index to continue from.
func stepRead(d *dec, packets []rawPacket, i int, pendingSection *int) (int, error) {
	switch packets[i].Class {
	case classControl:
		if err := handleControl(d, pendingSection, packets[i]); err != nil {
			return 0, err
		}
		return i + 1, nil
	case classSource:
		return 0, errors.New("tangeul/stream: a document has exactly one SOURCE record")
	}
	if *pendingSection >= 0 {
		return 0, errors.New("tangeul/stream: SECTION control must be followed by its heading PROSE")
	}
	return readStatement(d, packets, i)
}
