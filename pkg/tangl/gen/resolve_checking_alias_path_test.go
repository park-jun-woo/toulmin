//ff:func feature=tangl type=codegen control=sequence
//ff:what TestResolveCheckingAliasPath — covers extension fixup, absolute pass-through, and abs failure
package gen

import "testing"

func TestResolveCheckingAliasPath(t *testing.T) {
	t.Run("relative sibling with extension appended", func(t *testing.T) {
		got, err := resolveCheckingAliasPath("/repo/gate.tangeul", "shared/policy")
		if err != nil {
			t.Fatal(err)
		}
		if want := "/repo/shared/policy.tangeul"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("explicit extension preserved", func(t *testing.T) {
		got, err := resolveCheckingAliasPath("/repo/a.tangeul", "b.tangeul")
		if err != nil {
			t.Fatal(err)
		}
		if want := "/repo/b.tangeul"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("absolute package skips the join", func(t *testing.T) {
		got, err := resolveCheckingAliasPath("/repo/a.tangeul", "/elsewhere/b")
		if err != nil {
			t.Fatal(err)
		}
		if want := "/elsewhere/b.tangeul"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("abs failure on removed working directory", func(t *testing.T) {
		withRemovedCwd(t)
		if _, err := resolveCheckingAliasPath("rel/a.tangeul", "b"); err == nil {
			t.Fatal("expected filepath.Abs to fail on a removed working directory")
		}
	})
}
