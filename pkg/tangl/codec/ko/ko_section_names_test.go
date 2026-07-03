//ff:func feature=tangl type=parser control=sequence
//ff:what TestKoSectionNames — all seven Korean headings map to canonical English names; unknown reports false
package ko

import "testing"

func TestKoSectionNames(t *testing.T) {
	t.Run("SevenSections", func(t *testing.T) {
		want := map[string]string{
			"주제": "Subject",
			"참조": "See",
			"정의": "Definitions",
			"규칙": "Rules",
			"사안": "Cases",
			"업무": "Provides",
			"내부": "Internal",
		}
		for ko, en := range want {
			got, ok := koSectionNames(ko)
			if !ok {
				t.Errorf("koSectionNames(%q): ok = false, want true", ko)
			}
			if got != en {
				t.Errorf("koSectionNames(%q) = %q, want %q", ko, got, en)
			}
		}
	})

	t.Run("Unknown", func(t *testing.T) {
		if got, ok := koSectionNames("Subject"); ok || got != "" {
			t.Fatalf("koSectionNames(%q) = (%q, %v), want (\"\", false)", "Subject", got, ok)
		}
	})
}
