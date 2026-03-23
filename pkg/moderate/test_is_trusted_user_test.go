//ff:func feature=moderate type=rule control=iteration dimension=1
//ff:what TestIsTrustedUser — tests IsTrustedUser rule
package moderate

import "testing"

func TestIsTrustedUser(t *testing.T) {
	tests := []struct {
		name  string
		score float64
		min   float64
		want  bool
	}{
		{"trusted", 0.95, 0.9, true},
		{"not trusted", 0.5, 0.9, false},
		{"equal", 0.9, 0.9, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &ContentContext{Author: &Author{TrustScore: tt.score}}
			got, _ := IsTrustedUser(nil, ctx, &TrustScoreBacking{MinScore: tt.min})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
