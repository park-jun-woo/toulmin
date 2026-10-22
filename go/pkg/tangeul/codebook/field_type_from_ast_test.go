//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestFieldTypeFromAST — every surface value maps to a code that maps back; an unknown value fails
package codebook

import (
	"testing"
)

func TestFieldTypeFromAST(t *testing.T) {
	for _, v := range []string{"Text", "Integer", "Number", "Boolean", "Email", "Date", "Time", "URL", "Currency", "Quantity"} {
		c, err := FieldTypeFromAST(v)
		if err != nil {
			t.Fatalf("%v: %v", v, err)
		}
		if back, err := FieldTypeToAST(c); err != nil || back != v {
			t.Fatalf("%v: round trip got %v %v", v, back, err)
		}
	}
	if _, err := FieldTypeFromAST("Bogus"); err == nil {
		t.Fatal("unknown value: want error")
	}
}
