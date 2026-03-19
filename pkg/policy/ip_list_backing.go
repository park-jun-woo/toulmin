//ff:type feature=policy type=model
//ff:what IPListBacking: IP 목록 판정 기준 (용도 + 판정 함수)
package policy

// IPListBacking carries the judgment criteria for IP list checks.
// Purpose identifies the list's role (e.g., "blocklist", "whitelist").
// Check is the actual lookup function.
type IPListBacking struct {
	Purpose string
	Check   func(string) bool
}
