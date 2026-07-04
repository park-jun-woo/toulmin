//ff:func feature=tangl type=util control=sequence
//ff:what TestLoad — exercises the production Load wrapper over loadRegistry+loadDocAST
package registry

import "testing"

func TestLoad(t *testing.T) {
	// Load wires the production tangeul-backed loader into loadRegistry; a
	// path that cannot be read makes that loader fail, exercising the wrapper's
	// single delegating statement and its error return.
	if _, err := Load("/nonexistent/does-not-exist.tangeul"); err == nil {
		t.Fatal("expected an error loading a nonexistent document")
	}
}
