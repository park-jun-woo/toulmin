//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseBinding — parse a rule binding statement into RuleBinding AST node
package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// parseBinding parses: name is a role [of graph] using func [with specs] [qualified float]
func parseBinding(text string, lineNum int, parentGraph string, index int) (RuleBinding, error) {
	rb := RuleBinding{Line: lineNum, Index: index}

	var role string
	var rest string
	for _, r := range []string{" is a rule ", " is a counter ", " is an except "} {
		idx := strings.Index(text, r)
		if idx >= 0 {
			rb.Name = strings.TrimSpace(text[:idx])
			role = strings.TrimSpace(r)
			role = strings.TrimPrefix(role, "is a ")
			role = strings.TrimPrefix(role, "is an ")
			rest = text[idx+len(r):]
			break
		}
	}
	if role == "" {
		return RuleBinding{}, fmt.Errorf("invalid binding: no role found in %s", text)
	}
	rb.Role = strings.TrimSpace(role)

	qualIdx := strings.Index(rest, " qualified ")
	if qualIdx >= 0 {
		qualStr := strings.TrimSpace(rest[qualIdx+len(" qualified "):])
		q, err := strconv.ParseFloat(qualStr, 64)
		if err != nil {
			return RuleBinding{}, fmt.Errorf("invalid qualifier: %s", qualStr)
		}
		rb.Qualifier = q
		rest = rest[:qualIdx]
	}

	if strings.HasPrefix(rest, "of ") {
		usingIdx := strings.Index(rest, " using ")
		if usingIdx < 0 {
			return RuleBinding{}, fmt.Errorf("invalid binding: missing 'using' in %s", text)
		}
		rb.Graph = strings.TrimSpace(rest[3:usingIdx])
		rest = rest[usingIdx:]
	} else {
		rb.Graph = parentGraph
	}

	usingIdx := strings.Index(rest, "using ")
	if usingIdx < 0 {
		return RuleBinding{}, fmt.Errorf("invalid binding: missing 'using' in %s", text)
	}
	rest = rest[usingIdx+len("using "):]

	withIdx := strings.Index(rest, " with ")
	if withIdx >= 0 {
		funcRef := strings.TrimSpace(rest[:withIdx])
		rb.Func = funcRef
		specText := strings.TrimSpace(rest[withIdx+len(" with "):])
		specs, err := parseSpecCalls(specText)
		if err != nil {
			return RuleBinding{}, fmt.Errorf("invalid spec calls: %w", err)
		}
		rb.Specs = specs
	} else {
		rb.Func = strings.TrimSpace(rest)
	}

	if rb.Name == "" || rb.Func == "" {
		return RuleBinding{}, fmt.Errorf("invalid binding: empty name or func in %s", text)
	}
	return rb, nil
}
