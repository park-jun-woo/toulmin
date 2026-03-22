//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestIsRateLimited — tests IsRateLimited rule
package policy

import "testing"

func TestIsRateLimited(t *testing.T) {
	limiter := &mockLimiter{limited: map[string]bool{"1.2.3.4": true}}
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"limited", "1.2.3.4", true},
		{"not limited", "5.6.7.8", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsRateLimited(nil, &RequestContext{ClientIP: tt.ip, RateLimiter: limiter}, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
