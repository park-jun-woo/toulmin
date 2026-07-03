//ff:func feature=tangl type=parser control=sequence
//ff:what applyNodeClauseKo — parse one node tail clause ('기준' with-terms or '가중' qualified)
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// applyNodeClauseKo parses a single trailing node clause and returns the unconsumed
// remainder: '기준 <terms>' fills n.With (terms chained by 그리고), '가중 <float>' fills
// n.Qualified. An unrecognized clause is an error.
func applyNodeClauseKo(n *ast.Node, t string, line int, path string) (string, error) {
	if r, ok := takeKeywordKo(t, "기준"); ok {
		with, rest, err := parseWithTermsKo(r, line, path)
		if err != nil {
			return "", err
		}
		n.With = with
		return rest, nil
	}
	if r, ok := takeKeywordKo(t, "가중"); ok {
		q, rest, err := parseQualifiedKo(r, line, path)
		if err != nil {
			return "", err
		}
		n.Qualified = &q
		return rest, nil
	}
	return "", errAtKo(path, line, "unexpected node clause: %q", t)
}
