//ff:func feature=moderate type=engine control=sequence
//ff:what TestClassifierSpec_SpecName — tests SpecName returns the fixed spec name
package moderate

import "testing"

func TestClassifierSpec_SpecName(t *testing.T) {
	spec := &ClassifierSpec{Classifier: &mockClassifier{}}
	if got := spec.SpecName(); got != "ClassifierSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "ClassifierSpec")
	}
}
