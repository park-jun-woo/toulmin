//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestIsIPInList — tests IsIPInList rule
package policy

import "testing"

func TestIsIPInList(t *testing.T) {
	blocklist := &IPListBacking{Purpose: "blocklist", Check: func(ip string) bool { return ip == "1.2.3.4" }}
	whitelist := &IPListBacking{Purpose: "whitelist", Check: func(ip string) bool { return ip == "10.0.0.1" }}
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
			got, _ := IsIPInList(nil, &RequestContext{ClientIP: tt.ip}, tt.list)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
