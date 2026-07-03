//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestHasRuleNode — an empty slice, no rule node present, and a rule node found mid-scan
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestHasRuleNode covers hasRuleNode's loop: zero iterations, a full scan
// that finds no NodeRule, and a match found after a non-matching entry.
func TestHasRuleNode(t *testing.T) {
	cases := []struct {
		name      string
		primaries []packet.NodePacket
		want      bool
	}{
		{"empty", nil, false},
		{"noneMatch", []packet.NodePacket{{Type: codebook.NodeNamespace}, {Type: codebook.NodeAlias}}, false},
		{"matchSecond", []packet.NodePacket{{Type: codebook.NodeCase}, {Type: codebook.NodeRule}}, true},
	}
	for _, c := range cases {
		if got := hasRuleNode(c.primaries); got != c.want {
			t.Errorf("%s: hasRuleNode = %v, want %v", c.name, got, c.want)
		}
	}
}
