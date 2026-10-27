//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestFieldTypeRoundTrip — all 10 field-type names map and roundtrip, unknown errors
package codebook

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestFieldTypeRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		code toulmin.FieldTypeCode
	}{
		{"Text", toulmin.FieldText},
		{"Integer", toulmin.FieldInteger},
		{"Number", toulmin.FieldNumber},
		{"Boolean", toulmin.FieldBoolean},
		{"Email", toulmin.FieldEmail},
		{"Date", toulmin.FieldDate},
		{"Time", toulmin.FieldTime},
		{"URL", toulmin.FieldURL},
		{"Currency", toulmin.FieldCurrency},
		{"Quantity", toulmin.FieldQuantity},
	}
	for _, c := range cases {
		code, err := FieldTypeFromAST(c.name)
		if err != nil || code != c.code {
			t.Fatalf("FieldTypeFromAST(%q): got 0x%04X err=%v, want 0x%04X", c.name, uint16(code), err, uint16(c.code))
		}
		name, err := FieldTypeToAST(c.code)
		if err != nil || name != c.name {
			t.Fatalf("FieldTypeToAST(0x%04X): got %q err=%v, want %q", uint16(c.code), name, err, c.name)
		}
	}
	if _, err := FieldTypeFromAST("Blob"); err == nil {
		t.Error("FieldTypeFromAST(unknown): expected error")
	}
	if _, err := FieldTypeToAST(toulmin.FieldTypeCode(0xFFFF)); err == nil {
		t.Error("FieldTypeToAST(unknown): expected error")
	}
}
