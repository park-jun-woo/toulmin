//ff:func feature=tangl type=util control=sequence
//ff:what captureStderrOutput — runs fn with stderr redirected and returns what was logged
package tangl

import "testing"

// captureStderrOutput runs fn with stderr redirected and returns what was logged.
func captureStderrOutput(t *testing.T, fn func() error) []byte {
	t.Helper()
	out, _ := captureStderrBoth(t, fn)
	return out
}
