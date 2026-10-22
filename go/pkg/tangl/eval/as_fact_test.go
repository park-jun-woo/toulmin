//ff:func feature=tangl type=engine control=sequence
//ff:what asFact — asserts an evidence value is a toulmin.Fact and returns it
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// asFact asserts that ev is a toulmin.Fact and returns it, failing the test
// otherwise.
func asFact(t *testing.T, ev any) toulmin.Fact {
	t.Helper()
	f, ok := ev.(toulmin.Fact)
	if !ok {
		t.Fatalf("evidence is not a Fact: %T", ev)
	}
	return f
}
