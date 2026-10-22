//ff:func feature=tangl type=util control=sequence
//ff:what captureStderrBoth — redirects os.Stderr while fn runs and returns what was logged plus fn's error
package tangl

import (
	"io"
	"os"
	"testing"
)

// captureStderrBoth redirects os.Stderr to a pipe for the duration of fn,
// then restores it and returns everything written to stderr along with
// fn's own return value.
func captureStderrBoth(t *testing.T, fn func() error) ([]byte, error) {
	t.Helper()
	orig := os.Stderr
	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("failed to create pipe: %v", pipeErr)
	}
	os.Stderr = w

	fnErr := fn()

	w.Close()
	os.Stderr = orig
	out, _ := io.ReadAll(r)
	return out, fnErr
}
