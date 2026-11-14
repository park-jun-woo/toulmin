//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestShortNameSpec — verifies shortName handles # spec suffix correctly
package toulmin

import "testing"

func TestShortNameSpec(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"github.com/example/pkg.IsInRole#&{0.5}", "IsInRole#&{0.5}"},
		{"github.com/example/pkg.IsAdult", "IsAdult"},
		{"IsAdult", "IsAdult"},
		{"pkg.Fn#specA+specB", "Fn#specA+specB"},
		{"IsAdult#specA", "IsAdult#specA"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := shortName(tt.input)
			if got != tt.want {
				t.Errorf("shortName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
