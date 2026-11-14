//ff:func feature=policy type=engine control=sequence
//ff:what TestHeaderSpec_SpecName — verifies HeaderSpec.SpecName returns fixed name
package policy

import "testing"

func TestHeaderSpec_SpecName(t *testing.T) {
	if got := (&HeaderSpec{}).SpecName(); got != "HeaderSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "HeaderSpec")
	}
}
