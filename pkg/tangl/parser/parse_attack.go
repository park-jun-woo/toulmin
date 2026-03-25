//ff:func feature=tangl type=parser control=sequence
//ff:what parseAttack — parse an attack declaration into AttackDecl AST node
package parser

import (
	"fmt"
	"strings"
)

// parseAttack parses: attacker attacks target
func parseAttack(text string, lineNum int, parentGraph string) (AttackDecl, error) {
	parts := strings.SplitN(text, " attacks ", 2)
	if len(parts) != 2 {
		return AttackDecl{}, fmt.Errorf("invalid attack: %s", text)
	}
	attacker := strings.TrimSpace(parts[0])
	target := strings.TrimSpace(parts[1])
	if attacker == "" || target == "" {
		return AttackDecl{}, fmt.Errorf("invalid attack: empty attacker or target in %s", text)
	}
	return AttackDecl{Attacker: attacker, Target: target, Graph: parentGraph, Line: lineNum}, nil
}
