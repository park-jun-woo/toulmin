//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestEncodeErrorsKo — unknown section, out-of-order, missing 주제, missing 사안 are rejected
package ko

import "testing"

func TestEncodeErrorsKo(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{"UnknownSection", "## tangl:엉뚱\n- 이 문서는 `x` 이다\n"},
		{"OutOfOrder", "## tangl:사안\n- `c` 의 경우에\n\n## tangl:주제\n- 이 문서는 `x` 이다\n"},
		{"MissingSubject", "## tangl:사안\n- `c` 의 경우에\n"},
		{"MissingCases", "## tangl:주제\n- 이 문서는 `x` 이다\n"},
		{"BadStatement", "## tangl:주제\n- 이 문서는 `x` 이다\n\n## tangl:사안\n- `c` 의 경우에\n  - 알 수 없는 문장\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := (KoCodec{}).Encode([]byte(c.src), "t.md"); err == nil {
				t.Fatalf("%s: expected an error", c.name)
			}
		})
	}
}
