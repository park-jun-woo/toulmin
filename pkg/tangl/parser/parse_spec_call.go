//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseSpecCalls — parse with clause spec calls into SpecCall list
package parser

import (
	"fmt"
	"strings"
)

// parseSpecCalls parses: Role("admin") and IPList("block") or Role("admin"), IPList("block") and Header("X-Token")
func parseSpecCalls(text string) ([]SpecCall, error) {
	var result []SpecCall
	segments := splitSpecSegments(text)
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		sc, err := parseSingleSpec(seg)
		if err != nil {
			return nil, err
		}
		result = append(result, sc)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no spec calls found in: %s", text)
	}
	return result, nil
}
