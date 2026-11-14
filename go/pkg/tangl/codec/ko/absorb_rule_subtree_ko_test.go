//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestAbsorbRuleSubtreeKo — leaf, a child exceeding the head, a child within it, and nested recursion
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestAbsorbRuleSubtreeKo(t *testing.T) {
	cases := []struct {
		name string
		item scan.Item
		want uint32
	}{
		{
			name: "Leaf",
			item: scan.Item{Span: scan.Span{Off: 5, Len: 10}},
			want: 15,
		},
		{
			name: "ChildExceedsHead",
			item: scan.Item{
				Span:     scan.Span{Off: 0, Len: 5},
				Children: []scan.Item{{Span: scan.Span{Off: 2, Len: 20}}},
			},
			want: 22,
		},
		{
			name: "ChildWithinHead",
			item: scan.Item{
				Span:     scan.Span{Off: 0, Len: 100},
				Children: []scan.Item{{Span: scan.Span{Off: 0, Len: 5}}},
			},
			want: 100,
		},
		{
			name: "NestedGrandchildExceeds",
			item: scan.Item{
				Span: scan.Span{Off: 0, Len: 5},
				Children: []scan.Item{
					{
						Span:     scan.Span{Off: 0, Len: 5},
						Children: []scan.Item{{Span: scan.Span{Off: 10, Len: 50}}},
					},
				},
			},
			want: 60,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := absorbRuleSubtreeKo(c.item); got != c.want {
				t.Errorf("absorbRuleSubtreeKo() = %d, want %d", got, c.want)
			}
		})
	}
}
