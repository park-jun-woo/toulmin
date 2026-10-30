//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDefinitionKoFileChildren — a file definition with child fields is rejected
package ko

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseDefinitionKoFileChildren(t *testing.T) {
	it := scan.Item{
		Text:     "`도메인` 은 `JSON 목록` 을 \"./d.json\" 파일에서 읽는다",
		No:       2,
		Children: []scan.Item{{Text: "`필드` 는 1 을 뜻한다", No: 3}},
	}
	if _, err := parseDefinitionKo(it, "t.md"); err == nil || !strings.Contains(err.Error(), "file definition cannot have child fields") {
		t.Fatalf("err = %v", err)
	}
}
