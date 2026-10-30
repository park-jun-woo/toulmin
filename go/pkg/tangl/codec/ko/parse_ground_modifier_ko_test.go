//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseGroundModifierKo — '`a` 그리고 `b` 를 근거로 하는' yields names and the tail; anything else is unmatched
package ko

import (
	"reflect"
	"testing"
)

func TestParseGroundModifierKo(t *testing.T) {
	names, tail, ok := parseGroundModifierKo("`시도 기록` 그리고 `로그` 를 근거로 하는 `래칫`.`f` 쓰는")
	if !ok || !reflect.DeepEqual(names, []string{"시도 기록", "로그"}) || tail != " `래칫`.`f` 쓰는" {
		t.Fatalf("names=%v tail=%q ok=%v", names, tail, ok)
	}
	for _, s := range []string{
		"근거 없음",
		"`래칫`.`f` 쓰는",
		"`a` 를 기준으로 하는",
		"`a` 를 근거로 쓰는",
		"`a` 그리고 근거로 하는",
	} {
		names, tail, ok := parseGroundModifierKo(s)
		if ok || names != nil || tail != s {
			t.Errorf("%q: names=%v tail=%q ok=%v, want unmatched with input preserved", s, names, tail, ok)
		}
	}
}
