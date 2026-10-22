//ff:func feature=tangl type=util control=sequence
//ff:what TestTextsOf — TEXT ids resolve to their strings in order; an empty list is nil
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

func TestTextsOf(t *testing.T) {
	d := newDec(testWire)
	d.texts[3], d.texts[4] = "a", "b"
	if got := textsOf(d, []toulmin.TID{4, 3}); !reflect.DeepEqual(got, []string{"b", "a"}) {
		t.Fatalf("got %v", got)
	}
	if got := textsOf(d, nil); got != nil {
		t.Fatalf("empty: %v", got)
	}
}
