//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestFieldTypeToAST — every code maps to its surface value; an unknown code fails
package codebook

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestFieldTypeToAST(t *testing.T) {
	for _, v := range []string{"Text", "Integer", "Number", "Boolean", "Email", "Date", "Time", "URL", "Currency", "Quantity"} {
		c, _ := FieldTypeFromAST(v)
		if got, err := FieldTypeToAST(c); err != nil || got != v {
			t.Fatalf("%v: got %v %v", v, got, err)
		}
	}
	if _, err := FieldTypeToAST(toulmin.FieldTypeCode(0x9999)); err == nil {
		t.Fatal("unknown code: want error")
	}
}
