//ff:func feature=tangl type=parser control=sequence
//ff:what templateFullExample — TANGL full example template for testing
package parser

// templateFullExample is the full TANGL example from the grammar spec.
var templateFullExample = `# API 접근 제어 정책

인증된 사용자만 접근 가능하다.

## tangl:Import
- policy is from "github.com/park-jun-woo/toulmin/pkg/policy"

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

## 판정 기준

- verdict > 0이면 허용
- verdict <= 0이면 거부
`
