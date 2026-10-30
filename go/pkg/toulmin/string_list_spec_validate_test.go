//ff:func feature=engine type=engine control=sequence
//ff:what TestStringListSpecValidate — Validate accepts a named list Spec and rejects an empty name
package toulmin

import (
	"strings"
	"testing"
)

func TestStringListSpecValidate(t *testing.T) {
	if err := NewStringListSpec("list", []string{"a"}).Validate(); err != nil {
		t.Fatalf("named spec should validate, got %v", err)
	}
	err := NewStringListSpec("", nil).Validate()
	if err == nil || !strings.Contains(err.Error(), "name is empty") {
		t.Fatalf("expected empty-name error, got %v", err)
	}
}
