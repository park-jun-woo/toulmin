//ff:func feature=tangl type=util control=sequence
//ff:what withRemovedCwd — test helper that removes the working directory so os.Getwd fails
package gen

import (
	"os"
	"path/filepath"
	"testing"
)

// withRemovedCwd chdirs into a fresh directory and removes it, so any
// subsequent os.Getwd (as filepath.Abs performs on a relative path) fails.
// The original working directory is restored when the test ends.
func withRemovedCwd(t *testing.T) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	gone := filepath.Join(t.TempDir(), "gone")
	if err := os.MkdirAll(gone, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(gone); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(gone); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}
