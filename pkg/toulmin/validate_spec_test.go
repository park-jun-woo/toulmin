//ff:func feature=engine type=validator control=iteration dimension=1
//ff:what TestValidateSpec — tests validateSpec for fields-error, Validate-error, and success branches
package toulmin

import "testing"

func TestValidateSpec(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"FieldsError", func(t *testing.T) {
			err := validateSpec(&ruleIDUnmarshalableSpec{})
			if err == nil {
				t.Fatal("expected error for spec with func field")
			}
		}},
		{"ValidateError", func(t *testing.T) {
			err := validateSpec(&validateSpecFailingSpec{})
			if err == nil {
				t.Fatal("expected error from spec.Validate()")
			}
			if err.Error() != "domain validation failed" {
				t.Errorf("unexpected error: %v", err)
			}
		}},
		{"Success", func(t *testing.T) {
			err := validateSpec(&testSpec{Value: "ok"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
