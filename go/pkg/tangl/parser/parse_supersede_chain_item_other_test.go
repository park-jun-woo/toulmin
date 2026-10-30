//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseSupersedeChainItemOther — statements without the chain keyword are left unmatched without error
package parser

import "testing"

func TestParseSupersedeChainItemOther(t *testing.T) {
	for _, text := range []string{"`A` supersedes `B`", "supersede `A` > `B` > `C`", ""} {
		edges, ok, err := parseSupersedeChainItem(item{Text: text, Line: 1}, "t.md")
		if ok || err != nil || edges != nil {
			t.Errorf("%q: edges=%v ok=%v err=%v, want unmatched", text, edges, ok, err)
		}
	}
}
