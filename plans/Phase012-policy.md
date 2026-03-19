# Phase 012: 정책 프레임워크 — pkg/policy (구현 완료)

## 목표

toulmin 기반 정책 판정 프레임워크를 `pkg/policy`에 구현한다.
OPA/Rego 없이 Go 함수 + defeats graph로 인증/인가/접근 제어를 선언적으로 처리한다.
Gin 어댑터를 제공하여 미들웨어로 즉시 사용 가능하게 한다.

## 배경

### 현재 문제

1. **OPA는 무겁다**: Rego 언어 학습 + OPA 런타임 + 번들 관리가 필요하다. 단순 정책에도 인프라 오버헤드가 크다
2. **미들웨어 순서 의존성**: Gin 미들웨어 체인에서 인증 → 인가 → rate limit → IP 차단 순서가 암묵적이다
3. **예외 처리가 if-else**: 화이트리스트, 관리자 오버라이드 같은 예외가 if-else로 박힌다

### toulmin이 해결하는 것

- 정책 규칙 = Go 함수 (1-2 depth)
- 규칙 간 관계 = defeats graph (명시적 선언)
- 예외 = defeat edge (if-else 대체)
- 판정 근거 = EvaluateTrace (감사 로그, backing 값까지 추적)

### claim/ground/backing 분리

| 역할 | 정책 프레임워크에서 |
|---|---|
| claim | nil (라우트 매칭으로 확정) |
| ground | RequestContext (User, ClientIP, Headers) — 요청마다 다름 |
| backing | 판정 기준 (role명, IPListBacking, 헤더명) — 선언 시 고정 |
| verdict | > 0 허용, <= 0 거부 (undecided = 거부) |

## 구현 결과

### rule 함수 세트

| 함수 | backing | 설명 |
|---|---|---|
| IsAuthenticated | nil | User가 nil이 아닌지 |
| IsInRole | string | User.Role == backing |
| IsOwner | nil | User.ID == ResourceOwnerID |
| IsIPInList | *IPListBacking | backing.Check(ClientIP) |
| IsRateLimited | nil | RateLimiter.IsLimited(ClientIP) |
| HasHeader | string | Headers[backing] != "" |

IsAdminOverride는 `IsInRole` + backing "admin"으로 대체하여 제거.
IsInternalService는 `HasHeader` + backing "X-Internal-Token"으로 일반화.

### IPListBacking

```go
type IPListBacking struct {
    Purpose string           // "blocklist", "whitelist"
    Check   func(string) bool
}
```

함수는 개념(`IsIPInList` = "IP가 목록에 있는가"), backing은 구체적 기준(`Purpose` + `Check`).
포인터이므로 ruleID 충돌 없음. EvaluateTrace에서 Purpose가 보임.

### RequestContext

```go
type RequestContext struct {
    User            *User
    ClientIP        string
    ResourceOwnerID string
    Headers         map[string]string
    RateLimiter     RateLimiter
    Metadata        map[string]any
}
```

런타임 데이터만 포함. 판정 기준(IP 목록, 역할명 등)은 backing으로 분리.

### Guard / GuardDebug

```go
// Guard — Evaluate (경량), verdict <= 0 거부
func Guard(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc

// GuardDebug — EvaluateTrace, X-Policy-Verdict/X-Policy-Trace 헤더 + 응답 body에 trace
func GuardDebug(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc
```

### 사용 예시

```go
blocklist := &policy.IPListBacking{Purpose: "blocklist", Check: isBlocked}
whitelist := &policy.IPListBacking{Purpose: "whitelist", Check: isWhitelisted}

g := toulmin.NewGraph("admin:users")
auth := g.Warrant(policy.IsAuthenticated, nil, 1.0)
admin := g.Warrant(policy.IsInRole, "admin", 1.0)
blocked := g.Rebuttal(policy.IsIPInList, blocklist, 1.0)
allowed := g.Defeater(policy.IsIPInList, whitelist, 1.0)
internal := g.Defeater(policy.HasHeader, "X-Internal-Token", 1.0)
g.Defeat(blocked, auth)
g.Defeat(allowed, blocked)

r := gin.Default()
r.GET("/admin/users", policy.Guard(g, buildCtx), handler)
```

## 구현 과정 발견 사항

1. **IsAdminOverride는 IsInRole의 중복** — backing "admin"으로 대체. 별도 함수 불필요
2. **IsInternalService → HasHeader 일반화** — 헤더명을 backing으로 받으면 모든 헤더 체크에 재사용
3. **IPListBacking 구조체** — `func(string) bool`만 backing으로 넣으면 목적이 불명. Purpose 필드로 의미 명시, 포인터로 ruleID 충돌 해소
4. **ruleID 충돌 주의** — `fmt.Sprint(backing)` 결과가 같으면 같은 ruleID. 같은 backing 내용 = 같은 규칙이 맞음. 다른 목적이면 다른 backing 값을 사용

## 산출물

```
pkg/
  policy/
    user.go                       — User 구조체
    request_context.go            — RequestContext (런타임 데이터 + Headers + Metadata)
    rate_limiter.go               — RateLimiter 인터페이스
    ip_list_backing.go            — IPListBacking (Purpose + Check)
    rule_is_authenticated.go      — IsAuthenticated (backing: nil)
    rule_is_in_role.go            — IsInRole (backing: string)
    rule_is_owner.go              — IsOwner (backing: nil)
    rule_is_ip_in_list.go         — IsIPInList (backing: *IPListBacking)
    rule_is_rate_limited.go       — IsRateLimited (backing: nil)
    rule_has_header.go            — HasHeader (backing: string)
    gin_guard.go                  — Guard + GuardDebug
    format_trace.go               — formatTrace + formatVerdict
    rule_test.go                  — rule 함수 단위 테스트
    gin_guard_test.go             — Guard/GuardDebug 통합 테스트
```

## 검증 기준 (달성)

1. 모든 rule 함수가 `(claim, ground, backing)` 시그니처로 동작한다
2. IsInRole이 backing으로 역할명을 받아 재사용 가능하다
3. IsIPInList이 IPListBacking으로 blocklist/whitelist를 구분한다
4. Guard가 verdict <= 0이면 403을 반환한다
5. GuardDebug가 판정 근거를 헤더 + body에 포함한다
6. Defeat로 *Rule 참조를 전달하여 동작한다
7. 체이닝 없이 개별 문장으로 사용된다
8. `go test ./...` 전체 PASS

## 의존성

- Phase 010: backing 일급 시민
- Phase 011: Rule 참조 반환, Defeat(*Rule, *Rule), Graph 타입명
