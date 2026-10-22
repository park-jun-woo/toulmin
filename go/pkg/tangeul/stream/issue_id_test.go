//ff:func feature=tangl type=util control=sequence
//ff:what TestIssueID — hands out consecutive ids from the counter
package stream

import "testing"

func TestIssueID(t *testing.T) {
	e := testEnc(nil)
	if a, b := issueID(e), issueID(e); a != 1 || b != 2 || e.next != 3 {
		t.Fatalf("got %d %d next=%d", a, b, e.next)
	}
}
