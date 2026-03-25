//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateValid — validate passes for correct input
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateValid tests that validation passes for a correct TANGL file.
func TestValidateValid(t *testing.T) {
	input := `## tangl:Import
- policy is from "github.com/example/pkg"

## tangl:Rules
- rule "isAuthenticated" is
    return that user is not nil

## tangl:Graph
- access is a graph "api:access"
  - auth is a rule using isAuthenticated
  - admin is a rule using policy.IsInRole with policy.Role("admin")
  - blocked is a counter using policy.IsIPBlocked with policy.IPList("blocklist")
  - exempt is an except using policy.IsInternalIP
  - blocked attacks auth
  - exempt attacks blocked

## tangl:Evaluate
- acResult is results of evaluating access with trace
`
	f, err := parser.ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}
	if err := Validate(f); err != nil {
		t.Fatalf("Validate should pass, got: %v", err)
	}
}
