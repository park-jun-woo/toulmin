//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestGoFieldType — tests goFieldType for every case in its switch, including the default
package gen

import "testing"

func TestGoFieldType(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Text", "string"},
		{"Email", "string"},
		{"Date", "string"},
		{"Time", "string"},
		{"URL", "string"},
		{"Integer", "int"},
		{"Number", "float64"},
		{"Currency", "float64"},
		{"Quantity", "float64"},
		{"Boolean", "bool"},
		{"", "any"},
		{"custom type", "CustomType"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := goFieldType(tt.in)
			if got != tt.want {
				t.Errorf("goFieldType(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
