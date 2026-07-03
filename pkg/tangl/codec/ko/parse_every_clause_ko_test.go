//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseEveryClauseKo — interval-only, interval+until, missing 까지, bad interval
package ko

import "testing"

func TestParseEveryClauseKo(t *testing.T) {
	t.Run("IntervalOnly", func(t *testing.T) {
		iv, until, err := parseEveryClauseKo("15분마다", 1, "t.md")
		if err != nil || iv != "15m" || until != "" {
			t.Fatalf("iv=%q until=%q err=%v", iv, until, err)
		}
	})

	t.Run("IntervalAndUntil", func(t *testing.T) {
		iv, until, err := parseEveryClauseKo("1초마다, `서빙 가능` 까지", 1, "t.md")
		if err != nil || iv != "1s" || until != "서빙 가능" {
			t.Fatalf("iv=%q until=%q err=%v", iv, until, err)
		}
	})

	t.Run("MissingUntilKeyword", func(t *testing.T) {
		if _, _, err := parseEveryClauseKo("1초마다, `서빙 가능`", 1, "t.md"); err == nil {
			t.Fatal("expected error when 까지 is absent")
		}
	})

	t.Run("BadInterval", func(t *testing.T) {
		if _, _, err := parseEveryClauseKo("마다", 1, "t.md"); err == nil {
			t.Fatal("expected error for an empty interval")
		}
	})

	t.Run("UntilMissingBacktick", func(t *testing.T) {
		if _, _, err := parseEveryClauseKo("1초마다, 서빙 까지", 1, "t.md"); err == nil {
			t.Fatal("expected error for a missing until case name")
		}
	})

	t.Run("UntilTrailingText", func(t *testing.T) {
		if _, _, err := parseEveryClauseKo("1초마다, `서빙 가능` 여분텍스트 까지", 1, "t.md"); err == nil {
			t.Fatal("expected error for trailing text between the until case name and '까지'")
		}
	})
}
