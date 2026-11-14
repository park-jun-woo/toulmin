//ff:type feature=policy type=model
//ff:what IPListSpec: IP 목록 판정 기준 (용도 + IP 목록)
package policy

// IPListSpec carries the judgment criteria for IP list checks.
// Purpose identifies the list's role (e.g., "blocklist", "whitelist").
// List holds the IP addresses.
type IPListSpec struct {
	Purpose string
	List    []string
}
