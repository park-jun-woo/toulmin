//ff:func feature=tangl type=parser control=selection
//ff:what renderCertaintyKo — render a Certainty gate as `확신 <N>% <문턱>이면` (the inverse of parseCertaintyKo)
package ko

import (
	"fmt"
	"strconv"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderCertaintyKo renders a confidence gate back to "확신 <N>% <문턱>", where the
// compound threshold token is 이상이면/초과면/미만이면/이하면 for the canonical Certainty.Op
// at least/above/less than/at most (the reverse of certaintyOpKo). An unknown Op is
// an error.
func renderCertaintyKo(c ast.Certainty) (string, error) {
	var token string
	switch c.Op {
	case "at least":
		token = "이상이면"
	case "above":
		token = "초과면"
	case "less than":
		token = "미만이면"
	case "at most":
		token = "이하면"
	default:
		return "", fmt.Errorf("ko: unknown certainty operator %q", c.Op)
	}
	return "확신 " + strconv.Itoa(c.Percent) + "% " + token, nil
}
