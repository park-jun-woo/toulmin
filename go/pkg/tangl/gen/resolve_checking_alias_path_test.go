//ff:func feature=tangl type=codegen control=sequence
//ff:what TestResolveCheckingAliasPath — covers extension fixup, absolute pass-through, and abs failure
package gen

import "testing"

func TestResolveCheckingAliasPath(t *testing.T) {
	t.Run("relative sibling with extension appended", func(t *testing.T) {
		got, err := resolveCheckingAliasPath("/repo/gate.geul", "shared/policy")
		if err != nil {
			t.Fatal(err)
		}
		if want := "/repo/shared/policy.geul"; got != want {
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

	t.Run("absolute package skips the join", func(t *testing.T) {
		got, err := resolveCheckingAliasPath("/repo/a.geul", "/elsewhere/b")
		if err != nil {
			t.Fatal(err)
		}
		if want := "/elsewhere/b.geul"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("abs failure on removed working directory", func(t *testing.T) {
		withRemovedCwd(t)
		if _, err := resolveCheckingAliasPath("rel/a.geul", "b"); err == nil {
			t.Fatal("expected filepath.Abs to fail on a removed working directory")
		}
	})
}
