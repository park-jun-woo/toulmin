//ff:func feature=tangl type=parser control=selection
//ff:what parseLine — classify and parse a single TANGL statement line
package parser

import (
	"fmt"
	"strings"
)

// parseLine classifies a statement line and delegates to the appropriate parser.
func parseLine(line string, lineNum int, parentGraph string) (any, error) {
	switch {
	case strings.Contains(line, " is from \""):
		return parseImport(line, lineNum)
	case strings.Contains(line, " is a graph \""):
		return parseGraph(line, lineNum)
	case strings.Contains(line, " attacks "):
		return parseAttack(line, lineNum, parentGraph)
	case strings.Contains(line, " is results of evaluating "):
		return parseEval(line, lineNum)
	case isBinding(line):
		return parseBinding(line, lineNum, parentGraph, 0)
	default:
		return nil, fmt.Errorf("unrecognized statement: %s", line)
	}
}
