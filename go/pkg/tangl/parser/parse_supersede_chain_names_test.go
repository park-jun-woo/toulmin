//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseSupersedeChainNames — names split on '>'; missing backticks, dangling '>', bad separators, duplicates and short chains fail
package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseSupersedeChainNames(t *testing.T) {
	got, err := parseSupersedeChainNames(" `A` > `B`>`C` ", "t.md", 2)
	if err != nil || !reflect.DeepEqual(got, []string{"A", "B", "C"}) {
		t.Fatalf("valid = %v, %v", got, err)
	}
	for _, c := range []struct{ in, want string }{
		{"A > B > C", "expected backtick-quoted name"},
		{"`A` > `B` >", "expected backtick-quoted name"},
		{"`A` < `B` > `C`", "expected '>'"},
		{"`A` > `B` > `A`", "duplicate name"},
		{"`A` > `B`", "at least 3 names"},
	} {
		if _, err := parseSupersedeChainNames(c.in, "t.md", 2); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: err=%v, want %q", c.in, err, c.want)
		}
	}
}
