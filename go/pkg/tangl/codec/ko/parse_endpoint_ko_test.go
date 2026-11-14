//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseEndpointKo — endpoint head name, missing 한다면, missing backtick, trailing text
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseEndpointKo(t *testing.T) {
	t.Run("Head", func(t *testing.T) {
		ep, err := parseEndpointKo(scan.Item{Text: "`커피 만들기` 를 한다면", No: 30}, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ep.Name != "커피 만들기" || ep.Line != 30 {
			t.Fatalf("got %+v", ep)
		}
	})

	t.Run("MissingKeyword", func(t *testing.T) {
		if _, err := parseEndpointKo(scan.Item{Text: "`X` 진행한다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error when 한다면 is absent")
		}
	})

	t.Run("MissingBacktick", func(t *testing.T) {
		if _, err := parseEndpointKo(scan.Item{Text: "커피 를 한다면", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for a missing backtick name")
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		if _, err := parseEndpointKo(scan.Item{Text: "`X` `Y` 를 한다면", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for unexpected trailing text")
		}
	})
}
