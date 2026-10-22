//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseCertaintyKo — parse a '확신 N% <문턱>' confidence gate into an ast.Certainty
package ko

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseCertaintyKo parses the confidence gate "확신 <N>% <문턱>" (e.g.
// "확신 75% 이상이면"), producing an ast.Certainty with the integer percent and the
// canonical English threshold operator (certaintyOpKo). The trailing 문턱 token is a
// compound keyword (이상이면/초과면/미만이면/이하면) whose 면/이면 is not a standalone particle.
func parseCertaintyKo(cond string, line int, path string) (*ast.Certainty, error) {
	rest, ok := takeKeywordKo(strings.TrimSpace(cond), "확신")
	if !ok {
		return nil, errAtKo(path, line, "expected '확신' in certainty gate, got %q", cond)
	}
	rest = strings.TrimLeft(rest, " \t")
	i := 0
	for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
		i++
	}
	if i == 0 {
		return nil, errAtKo(path, line, "expected integer percent in certainty gate, got %q", rest)
	}
	pct, err := strconv.Atoi(rest[:i])
	if err != nil {
		return nil, errAtKo(path, line, "invalid percent %q: %v", rest[:i], err)
	}
	rest = strings.TrimLeft(rest[i:], " \t")
	if !strings.HasPrefix(rest, "%") {
		return nil, errAtKo(path, line, "expected '%%' after percent, got %q", rest)
	}
	op, ok := certaintyOpKo(strings.TrimSpace(rest[1:]))
	if !ok {
		return nil, errAtKo(path, line, "expected certainty threshold operator, got %q", rest[1:])
	}
	return &ast.Certainty{Op: op, Percent: pct}, nil
}
