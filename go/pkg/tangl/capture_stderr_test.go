//ff:func feature=tangl type=util control=sequence
//ff:what captureStderr — runs fn with stderr redirected away from the test output and returns fn's error
package tangl

import "testing"

// captureStderr runs fn with stderr redirected away from the test output and
// returns fn's error.
func captureStderr(t *testing.T, fn func() error) error {
	t.Helper()
	_, err := captureStderrBoth(t, fn)
	return err
}
