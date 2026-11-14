//ff:func feature=tangl type=util control=sequence
//ff:what TestSetErr — the first error sticks; nil and later errors change nothing
package stream

import (
	"errors"
	"testing"
)

func TestSetErr(t *testing.T) {
	e := testEnc(nil)
	setErr(e, nil)
	if e.err != nil {
		t.Fatalf("nil: got %v", e.err)
	}
	first := errors.New("first")
	setErr(e, first)
	setErr(e, errors.New("second"))
	setErr(e, nil)
	if e.err != first {
		t.Fatalf("got %v, want first", e.err)
	}
}
