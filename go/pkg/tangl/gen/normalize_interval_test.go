//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestNormalizeInterval — tests normalizeInterval for valid word-form durations and invalid input
package gen

import (
	"testing"
	"time"
)

func TestNormalizeInterval(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		want   time.Duration
		wantOk bool
	}{
		{"seconds word form", "30seconds", 30 * time.Second, true},
		{"minutes word form", "5minutes", 5 * time.Minute, true},
		{"hours word form", "2hours", 2 * time.Hour, true},
		{"already valid go duration", "1h30m", 90 * time.Minute, true},
		{"invalid duration string", "not-a-duration", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := normalizeInterval(tt.in)
			if ok != tt.wantOk {
				t.Fatalf("normalizeInterval(%q) ok = %v, want %v", tt.in, ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("normalizeInterval(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
