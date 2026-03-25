//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestIsIPInList — tests IsIPInList rule
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsIPInList(t *testing.T) {
	blocklist := &IPListBacking{Purpose: "blocklist", List: []string{"1.2.3.4"}}
	whitelist := &IPListBacking{Purpose: "whitelist", List: []string{"10.0.0.1"}}
	tests := []struct {
		name string
		ip   string
		list *IPListBacking
		want bool
	}{
		{"blocked", "1.2.3.4", blocklist, true},
		{"not blocked", "5.6.7.8", blocklist, false},
		{"whitelisted", "10.0.0.1", whitelist, true},
		{"not whitelisted", "5.6.7.8", whitelist, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("clientIP", tt.ip)
			got, _ := IsIPInList(ctx, tt.list)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
