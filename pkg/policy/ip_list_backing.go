//ff:type feature=policy type=model
//ff:what IPListBacking: IP 목록 판정 기준 (용도 + IP 목록)
package policy

// IPListBacking carries the judgment criteria for IP list checks.
// Purpose identifies the list's role (e.g., "blocklist", "whitelist").
// List holds the IP addresses.
type IPListBacking struct {
	Purpose string
	List    []string
}
