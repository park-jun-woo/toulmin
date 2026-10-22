//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeCondPrefixKo — 그리고→and, 또는→or, no prefix
package ko

import "testing"

func TestTakeCondPrefixKo(t *testing.T) {
	t.Run("And", func(t *testing.T) {
		op, rest, ok := takeCondPrefixKo("그리고 `b` 가 비어 있으면")
		if !ok || op != "and" || rest != " `b` 가 비어 있으면" {
			t.Fatalf("op=%q rest=%q ok=%v", op, rest, ok)
		}
	})

	t.Run("Or", func(t *testing.T) {
		op, rest, ok := takeCondPrefixKo("또는 `직원` 이면")
		if !ok || op != "or" || rest != " `직원` 이면" {
			t.Fatalf("op=%q rest=%q ok=%v", op, rest, ok)
		}
	})

	t.Run("None", func(t *testing.T) {
		if _, _, ok := takeCondPrefixKo("`a` 가 비어 있으면"); ok {
			t.Fatal("expected ok=false when there is no 그리고/또는 prefix")
		}
	})
}
