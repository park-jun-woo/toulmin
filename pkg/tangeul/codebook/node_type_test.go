//ff:func feature=tangl type=model control=iteration dimension=1
//ff:what TestNodeType — each NodeType constant equals the §2.1 spec table code
package codebook

import "testing"

func TestNodeType(t *testing.T) {
	cases := []struct {
		got  NodeType
		want uint16
	}{
		{NodeNamespace, 0x0001},
		{NodeCase, 0x0002},
		{NodeRule, 0x0003},
		{NodeFunc, 0x0004},
		{NodeEndpoint, 0x0005},
		{NodeTrigger, 0x0006},
		{NodeTerm, 0x0007},
		{NodeField, 0x0008},
		{NodeAlias, 0x0009},
		{NodeGroup, 0x000A},
	}
	for _, c := range cases {
		if uint16(c.got) != c.want {
			t.Errorf("NodeType: got 0x%04X, want 0x%04X", uint16(c.got), c.want)
		}
	}
}
