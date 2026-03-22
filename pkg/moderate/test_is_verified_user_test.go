//ff:func feature=moderate type=rule control=iteration dimension=1
//ff:what TestIsVerifiedUser — tests IsVerifiedUser rule
package moderate

import "testing"

func TestIsVerifiedUser(t *testing.T) {
	tests := []struct {
		name     string
		verified bool
		want     bool
	}{
		{"verified", true, true},
		{"not verified", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &ContentContext{Author: &Author{Verified: tt.verified}}
			got, _ := IsVerifiedUser(nil, ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
