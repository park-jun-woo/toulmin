//ff:func feature=engine type=util control=sequence
//ff:what isolatedT — returns a *testing.T disconnected from the running test tree, so that calling Fatalf/Errorf on it does not propagate failure to the real enclosing test
package toulmin

import "testing"

// isolatedT returns a *testing.T disconnected from the running test tree, so that
// calling Fatalf/Errorf on it (which mark it failed) does not propagate failure to
// the real enclosing test. Fatalf/FailNow calls runtime.Goexit, so it must run in
// its own goroutine.
func isolatedT() *testing.T {
	return &testing.T{}
}
