//ff:func feature=tangl type=parser control=sequence
//ff:what parseEval — parse an evaluation declaration into EvalDecl AST node
package parser

import (
	"fmt"
	"strings"
)

// parseEval parses: name is results of evaluating graph [with trace|duration]
func parseEval(text string, lineNum int) (EvalDecl, error) {
	parts := strings.SplitN(text, " is results of evaluating ", 2)
	if len(parts) != 2 {
		return EvalDecl{}, fmt.Errorf("invalid eval: %s", text)
	}
	name := strings.TrimSpace(parts[0])
	rest := strings.TrimSpace(parts[1])

	ed := EvalDecl{Name: name, Line: lineNum}

	withIdx := strings.Index(rest, " with ")
	if withIdx >= 0 {
		ed.Graph = strings.TrimSpace(rest[:withIdx])
		option := strings.TrimSpace(rest[withIdx+len(" with "):])
		switch option {
		case "trace":
			ed.Trace = true
		case "duration":
			ed.Duration = true
		default:
			return EvalDecl{}, fmt.Errorf("unknown eval option: %s", option)
		}
	} else {
		ed.Graph = rest
	}

	if ed.Name == "" || ed.Graph == "" {
		return EvalDecl{}, fmt.Errorf("invalid eval: empty name or graph in %s", text)
	}
	return ed, nil
}
