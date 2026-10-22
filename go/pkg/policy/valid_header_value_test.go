//ff:func feature=policy type=engine control=sequence
//ff:what validHeaderValue — checks whether v is a map[string]string with X-Test header equal to "1"
package policy

func validHeaderValue(v any) bool {
	h, ok := v.(map[string]string)
	return ok && h["X-Test"] == "1"
}
