//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCompareKoRhsError — a right operand followed by stray text is rejected
package ko

import (
	"strings"
	"testing"
)

func TestParseCompareKoRhsError(t *testing.T) {
	expr, err := parseCompareKo("`등급` 이 `기대 등급` 덤 과 같으면", 4, "t.md")
	if err == nil || expr != nil || !strings.Contains(err.Error(), "unexpected trailing text after operand") {
		t.Fatalf("expr=%v err=%v, want the trailing-text error", expr, err)
	}
}
