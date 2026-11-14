//ff:func feature=tangl type=util control=sequence
//ff:what TestLoadDocAST — covers successful parse and load-error propagation
package registry

import "testing"

func TestLoadDocAST(t *testing.T) {
	t.Run("valid .geul returns its AST", func(t *testing.T) {
		d, err := loadDocAST("../testdata/access.geul")
		if err != nil {
			t.Fatal(err)
		}
		if d == nil {
			t.Fatal("expected a non-nil AST")
		}
	})

	t.Run("unreadable path propagates the load error", func(t *testing.T) {
		if _, err := loadDocAST("/nonexistent/does-not-exist.geul"); err == nil {
			t.Fatal("expected an error loading a nonexistent document")
		}
	})
}
