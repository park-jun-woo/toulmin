//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parsePath — parse a backtick dot-chain data path (a.b.c)
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parsePath parses a backtick dot-chain data path (one or more backtick names
// joined by dots) and returns the Path, the remaining text, and whether a
// well-formed path was found. Unlike parseRef, it accepts an unbounded number
// of steps.
func parsePath(s string) (ast.Path, string, bool) {
	first, rest, ok := takeBacktick(s)
	if !ok {
		return ast.Path{}, s, false
	}
	steps := []string{first}
	for strings.HasPrefix(rest, ".") {
		next, r2, ok2 := takeBacktick(rest[1:])
		if !ok2 {
			return ast.Path{}, s, false
		}
		steps = append(steps, next)
		rest = r2
	}
	return ast.Path{Steps: steps}, rest, true
}
