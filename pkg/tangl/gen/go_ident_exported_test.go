//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestGoIdentExported — tests goIdentExported for underscore-passthrough and uppercase-first-letter branches
package gen

import "testing"

func TestGoIdentExported(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty input yields underscore passthrough", "", "_"},
		{"punctuation-only input yields underscore passthrough", "!!!", "_"},
		{"single word gets uppercase first letter", "order", "Order"},
		{"multi word gets uppercase first letter", "order received", "OrderReceived"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := goIdentExported(tt.in)
			if got != tt.want {
				t.Errorf("goIdentExported(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
