//ff:func feature=engine type=validator control=iteration dimension=1
//ff:what TestValidateSpecFields — tests validateSpecFields for pointer/value receiver, func-field error, and no-func-field success branches
package toulmin

import (
	"strings"
	"testing"
)

func TestValidateSpecFields(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"PointerNoFuncField", func(t *testing.T) {
			err := validateSpecFields(&testSpec{Value: "x"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}},
		{"ValueReceiverNoFuncField", func(t *testing.T) {
			err := validateSpecFields(validateSpecFieldsValueSpec{Value: "x"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}},
		{"FuncFieldError", func(t *testing.T) {
			err := validateSpecFields(&ruleIDUnmarshalableSpec{})
			if err == nil {
				t.Fatal("expected error for spec with func field")
			}
			if !strings.Contains(err.Error(), "func field") || !strings.Contains(err.Error(), "Fn") {
				t.Errorf("unexpected error message: %v", err)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
