//ff:func feature=state type=rule control=iteration dimension=1
//ff:what TestIsExpired — tests IsExpired rule
package state

import (
	"testing"
	"time"
)

func TestIsExpired(t *testing.T) {
	tests := []struct {
		name string
		exp  time.Time
		want bool
	}{
		{"expired", time.Now().Add(-1 * time.Hour), true},
		{"not expired", time.Now().Add(1 * time.Hour), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ExpiryBacking{ExpiresAt: tt.exp}
			ctx := &TransitionContext{Resource: &testResource{ExpiresAt: tt.exp}}
			got, _ := IsExpired(nil, ctx, b)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
