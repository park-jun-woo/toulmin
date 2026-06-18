//ff:func feature=engine type=engine control=sequence
//ff:what TestRuleID — tests ruleID stability and collision avoidance across spec value, type, and nested pointer cases
package toulmin

import (
	"testing"
)

func TestRuleID(t *testing.T) {
	// (a) same func + different value specs → different IDs
	idA1 := ruleID(ruleIDTestFn, Specs{&testSpec{Value: "x"}})
	idA2 := ruleID(ruleIDTestFn, Specs{&testSpec{Value: "y"}})
	if idA1 == idA2 {
		t.Errorf("(a) same func + different value specs should differ: %q == %q", idA1, idA2)
	}

	// (b) different type (SpecName differs) + same field value → different IDs
	idB1 := ruleID(ruleIDTestFn, Specs{&testSpec{Value: "same"}})
	idB2 := ruleID(ruleIDTestFn, Specs{&ruleIDAltSpec{Value: "same"}})
	if idB1 == idB2 {
		t.Errorf("(b) different spec types with same field value should differ: %q == %q", idB1, idB2)
	}

	// (c) same func + identical spec twice → same ID (stability)
	idC1 := ruleID(ruleIDTestFn, Specs{&testSpec{Value: "stable"}})
	idC2 := ruleID(ruleIDTestFn, Specs{&testSpec{Value: "stable"}})
	if idC1 != idC2 {
		t.Errorf("(c) identical specs should produce same ID: %q != %q", idC1, idC2)
	}

	// (d) nested pointer field, two equal-value specs (distinct instances) → same ID
	idD1 := ruleID(ruleIDTestFn, Specs{&ruleIDMemberSpec{
		Level:    "basic",
		Discount: &ruleIDDiscountSpec{Name: "promo", Rate: 0.1},
	}})
	idD2 := ruleID(ruleIDTestFn, Specs{&ruleIDMemberSpec{
		Level:    "basic",
		Discount: &ruleIDDiscountSpec{Name: "promo", Rate: 0.1},
	}})
	if idD1 != idD2 {
		t.Errorf("(d) equal-value specs with nested pointer should produce same ID (no address leak): %q != %q", idD1, idD2)
	}
}
