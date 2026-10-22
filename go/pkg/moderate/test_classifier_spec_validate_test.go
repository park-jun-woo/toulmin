//ff:func feature=moderate type=engine control=sequence
//ff:what TestClassifierSpec_Validate — tests validation of ClassifierSpec classifier presence
package moderate

import "testing"

func TestClassifierSpec_Validate(t *testing.T) {
	if err := (&ClassifierSpec{Classifier: nil}).Validate(); err == nil {
		t.Fatal("expected error for nil classifier")
	}
	if err := (&ClassifierSpec{Classifier: &mockClassifier{}}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
