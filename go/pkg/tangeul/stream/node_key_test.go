//ff:func feature=tangl type=util control=sequence
//ff:what TestNodeKey — nodeKey formats "<caseIdx>\x00<name>", scoping node names to their case
package stream

import "testing"

// TestNodeKey directly checks nodeKey's NUL-separated key format (already
// exercised indirectly through assignTIDs in tid_table_test.go), including
// that distinct case indices with the same name produce distinct keys.
func TestNodeKey(t *testing.T) {
	if got, want := nodeKey(0, "balance sufficient"), "0\x00balance sufficient"; got != want {
		t.Errorf("nodeKey(0, ...) = %q, want %q", got, want)
	}
	if got, want := nodeKey(4, "balance sufficient"), "4\x00balance sufficient"; got != want {
		t.Errorf("nodeKey(4, ...) = %q, want %q", got, want)
	}
	if got, want := nodeKey(0, ""), "0\x00"; got != want {
		t.Errorf("nodeKey(0, \"\") = %q, want %q", got, want)
	}
	if nodeKey(0, "balance sufficient") == nodeKey(4, "balance sufficient") {
		t.Errorf("keys for different case indices collided")
	}
}
