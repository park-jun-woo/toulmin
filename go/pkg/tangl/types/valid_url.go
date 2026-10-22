//ff:func feature=tangl type=util control=sequence
//ff:what validURL — reports whether v is a string parseable as an absolute URL
package types

import "net/url"

// validURL reports whether v is a string parseable as an absolute URL
// (non-empty scheme and host).
func validURL(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}
