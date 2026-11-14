//ff:func feature=tangl type=util control=sequence
//ff:what TestResolveCheckingAliasPath — tests relative resolution and .geul extension fixup
package registry

import (
	"path/filepath"
	"testing"
)

func TestResolveCheckingAliasPath(t *testing.T) {
	t.Run("relative sibling with extension appended", func(t *testing.T) {
		got, err := resolveCheckingAliasPath("/repo/reins/gate.geul", "shared/policy")
		if err != nil {
			t.Fatal(err)
		}
		if want := "/repo/reins/shared/policy.geul"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("parent-relative path is cleaned", func(t *testing.T) {
		got, err := resolveCheckingAliasPath("/repo/turn/turn.geul", "../gate/gate")
		if err != nil {
			t.Fatal(err)
		}
		if want := "/repo/gate/gate.geul"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("explicit extension preserved", func(t *testing.T) {
		got, err := resolveCheckingAliasPath("/repo/a.geul", "b.geul")
		if err != nil {
			t.Fatal(err)
		}
		if want := "/repo/b.geul"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("result is absolute", func(t *testing.T) {
		got, err := resolveCheckingAliasPath("/repo/a.geul", "b")
		if err != nil {
			t.Fatal(err)
		}
		if !filepath.IsAbs(got) {
			t.Errorf("expected absolute path, got %q", got)
		}
	})

	t.Run("abs failure on removed working directory", func(t *testing.T) {
		// A relative referencing path keeps the joined path relative, so
		// filepath.Abs must consult the working directory — which no longer
		// exists — and returns an error.
		withRemovedCwd(t)
		if _, err := resolveCheckingAliasPath("rel/a.geul", "b"); err == nil {
			t.Fatal("expected filepath.Abs to fail on a removed working directory")
		}
	})
}
