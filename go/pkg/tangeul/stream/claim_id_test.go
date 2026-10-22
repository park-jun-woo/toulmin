//ff:func feature=tangl type=util control=sequence
//ff:what TestClaimID — a fresh nonzero id registers; zero and repeats fail
package stream

import "testing"

func TestClaimID(t *testing.T) {
	d := newDec(testWire)
	if err := claimID(d, 5); err != nil || !d.seen[5] {
		t.Fatalf("fresh: %v", err)
	}
	if err := claimID(d, 5); err == nil {
		t.Fatal("repeat: want error")
	}
	if err := claimID(d, 0); err == nil {
		t.Fatal("zero: want error")
	}
}
