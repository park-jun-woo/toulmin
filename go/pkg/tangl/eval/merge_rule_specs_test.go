//ff:func feature=tangl type=engine control=sequence
//ff:what TestMergeRuleSpecsRejectsConflict — equal same-name Specs merge, new names append, different same-name Specs conflict
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestMergeRuleSpecsRejectsConflict(t *testing.T) {
	_, err := mergeRuleSpecs(toulmin.Specs{&derivedSpec{name: "x"}}, toulmin.Specs{&derivedSpec{name: "x"}})
	if err != nil {
		t.Fatalf("equal specs conflict: %v", err)
	}
	merged, err := mergeRuleSpecs(toulmin.Specs{&derivedSpec{name: "x"}}, toulmin.Specs{&testSpec{Name: "y"}})
	if err != nil || len(merged) != 2 || merged[1].SpecName() != "y" {
		t.Fatalf("new supplied spec must be appended: %v %v", merged, err)
	}
	_, err = mergeRuleSpecs(toulmin.Specs{&derivedSpec{name: "x"}}, toulmin.Specs{&testSpec{Name: "x"}})
	if err == nil {
		t.Fatal("different specs with the same SpecName must conflict")
	}
}
