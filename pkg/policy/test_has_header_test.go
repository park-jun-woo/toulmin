//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestHasHeader — tests HasHeader rule
package policy

import "testing"

func TestHasHeader(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		header  string
		want    bool
	}{
		{"has token", map[string]string{"X-Internal-Token": "secret"}, "X-Internal-Token", true},
		{"no token", map[string]string{}, "X-Internal-Token", false},
		{"empty token", map[string]string{"X-Internal-Token": ""}, "X-Internal-Token", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := HasHeader(nil, &RequestContext{Headers: tt.headers}, tt.header)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
