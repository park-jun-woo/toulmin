# toulmin

**if-else를 그만 중첩하라. 규칙을 선언하고, 관계를 선언하라.**

TypeScript, Python, Go 룰 엔진. 규칙은 함수다. 예외는 그래프 엣지다. DSL 없음. 사이드카 없음. 새 언어 없음.

### TypeScript

```typescript
const isAuthenticated = (ctx, specs) => [ctx.get("user") != null, null]
const isIPBlocked     = (ctx, specs) => [blockedIPs.has(ctx.get("ip")), null]
const isInternalIP    = (ctx, specs) => [ctx.get("ip")?.startsWith("10."), null]
const isRateLimited   = (ctx, specs) => [/* ... */]
const isPremiumUser   = (ctx, specs) => [/* ... */]
const isIncidentMode  = (ctx, specs) => [/* ... */]

const g = new Graph("api:access")
const auth    = g.rule(isAuthenticated)
const blocked = g.counter(isIPBlocked)
const exempt  = g.except(isInternalIP)
blocked.attacks(auth)
exempt.attacks(blocked)

const limited  = g.counter(isRateLimited)    // 화요일: Rate limiting
limited.attacks(auth)
const premium  = g.except(isPremiumUser)     // 수요일: 프리미엄 면제
premium.attacks(limited)
const incident = g.counter(isIncidentMode)   // 목요일: 장애 대응 제한
incident.attacks(premium)

const results = g.evaluate(newContext())
// results[0].verdict > 0: 허용
```

### Python (planned)

```python
g = Graph("api:access")
auth    = g.rule(is_authenticated)
blocked = g.counter(is_ip_blocked)
exempt  = g.except_(is_internal_ip)
blocked.attacks(auth)
exempt.attacks(blocked)

results = g.evaluate(MapContext())
```

### Go

```go
g := toulmin.NewGraph("api:access")
auth    := g.Rule(isAuthenticated)
blocked := g.Counter(isIPBlocked)
exempt  := g.Except(isInternalIP)
blocked.Attacks(auth)
exempt.Attacks(blocked)

results, _ := g.Evaluate(ctx)
```

요구사항은 진화한다. 매일 2줄 추가, 기존 코드 변경 없음. 같은 진화를 if-else로:

```go
// 월요일
func isAuthenticated(ctx Context, specs Specs) (bool, any) {
    req, _ := ctx.Get("req")
    return req.(*Req).User != nil, nil
}

// 목요일 — 4단계 중첩
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") { allow = true }
    } else if isRateLimited(ip) {
        if isPremium(user) {
            if !incidentMode { allow = true }  // 추적 불가
        }
    } else { allow = true }
}
```

toulmin: **요구사항마다 2줄, 기존 코드 변경 없음.** if-else: **매번 전체 구조 재작성.**

## 설치

```bash
npm install rulecat          # TypeScript
pip install rulecat          # Python (planned)
go get github.com/park-jun-woo/toulmin/pkg/toulmin  # Go
```

## 핵심 개념

### 규칙은 함수다

```typescript
// TypeScript
const fn: RuleFunc = (ctx, specs) => [boolean, unknown]

// Python
def fn(ctx: Context, specs: list[Spec]) -> tuple[bool, Any]: ...

// Go
func fn(ctx Context, specs Specs) (bool, any)
```

- `ctx` = get/set 기반의 요청별 컨텍스트 (사용자, IP, 컨텍스트)
- `specs` = 그래프 선언 시 `.with()` / `.with_spec()` / `.With()`로 설정되는 판정 기준
- 반환 = `(판정 결과, 증거)`. 증거는 도메인별 자유 타입.

```typescript
const isInRole: RuleFunc = (ctx, specs) => {
    const user = ctx.get("user")
    if (!user) return [false, null]
    const role = (specs[0] as RoleSpec).role
    return [user.role === role, user.role]
}
```

### Defeats Graph

세 종류의 노드, 한 종류의 엣지:

| 노드 | 역할 |
|---|---|
| **Rule** | 주장 — 공격받을 수 있다 |
| **Counter** | 반박 — 주장을 공격한다 (자기 결론 있음) |
| **Except** | 예외 — 공격만 한다 (자기 결론 없음) |

```go
g := toulmin.NewGraph("voting")
auth     := g.Rule(isAdult)
criminal := g.Counter(hasCriminalRecord)
expunged := g.Except(isExpunged)
criminal.Attacks(auth)      // 전과가 성인을 공격
expunged.Attacks(criminal)  // 말소가 전과를 무력화
```

**Except가 차별점이다.** "예외의 예외"를 if-else로는 구조적으로 표현할 수 없다. defeats graph에서는 엣지 한 줄이다.

### Verdict

h-Categoriser가 연속값 `[-1, +1]`을 계산한다:

```
raw(a) = w(a) / (1 + Σ raw(attackers))
verdict = 2 × raw - 1
```

| Verdict | 의미 |
|---|---|
| +1.0 | 주장 성립 |
| 0.0 | 판정불가 |
| -1.0 | 완전 반박 |

+쪽이 긍정, -쪽이 부정. 이진 판정이 아니다. **도메인이 해석**한다:

- 접근 제어: `verdict > 0` → 허용, `verdict ≤ 0` → 거부
- 모더레이션: `verdict ≤ 0` → 차단, `0 < v ≤ 0.3` → 검토, `> 0.3` → 허용
- 피처 플래그: `verdict > 0` → 활성

규칙의 주장을 긍정형으로 설계하면 verdict가 직관에 부합한다 — "접근 허용"(+1 = 허용), "파일 구조 준수"(+1 = 준수).

### 평가 옵션

```go
ctx := toulmin.NewContext()
ctx.Set("req", req)
results, _ := g.Evaluate(ctx)                                                         // 기본 (h-Categoriser)
results, _ = g.Evaluate(ctx, toulmin.EvalOption{Trace: true})                          // trace 포함
results, _ = g.Evaluate(ctx, toulmin.EvalOption{Duration: true})                       // 소요시간 측정 (trace 자동 활성화)
```

`EvalOption`으로 평가 동작을 제어한다: `Method` (Matrix/Recursive (planned)), `Trace` (TraceEntry 수집), `Duration` (규칙별 실행 시간 측정).

### Spec

spec 값은 `Spec` 인터페이스를 구현해야 한다:

```go
type Spec interface {
    SpecName() string
    Validate() error
}
```

같은 함수 + 다른 spec = 다른 규칙. 클로저 팩토리 없이 규칙을 재사용한다:

```go
g := toulmin.NewGraph("access")
admin  := g.Rule(isInRole).With(&RoleSpec{Role: "admin"})
editor := g.Rule(isInRole).With(&RoleSpec{Role: "editor"}).Qualifier(0.8)
```

spec이 `nil`이면 규칙에 판정 기준이 필요 없다는 뜻이다. Spec 구조체에 func 필드는 금지된다 — `Validate()`가 이를 거부한다.

## Trace

`EvalOption{Trace: true}`로 각 규칙의 판정 근거를 추적한다:

```go
ctx := toulmin.NewContext()
results, _ := g.Evaluate(ctx, toulmin.EvalOption{Trace: true})
for _, t := range results[0].Trace {
    fmt.Printf("%s role=%s activated=%v evidence=%v\n",
        t.Name, t.Role, t.Activated, t.Evidence)
}
```

모더레이션 로그, 감사 추적, 디버깅 — 별도 로깅 없이 엔진이 제공한다.

## Run

`Evaluate`는 판정하고 반환만 한다 — 순수하고 멱등하다. `Run`은 먼저 판정한 뒤 **행동한다**:
전체 패스를 수행하고, 각 노드의 결과에 따라 정확히 하나의 핸들러를 발화한다.

| 이벤트 | 조건 | 의미 |
|---|---|---|
| `Active` | func true && verdict > 0 | 적용되어 우세 |
| `Defeated` | func true && verdict <= 0 | 적용됐으나 패배 |
| `Inactive` | func false | 규칙 미적용 |

```go
g.Rule(isAuthenticated).
    OnActive(func(ctx toulmin.Context, ev toulmin.NodeEvent, view toulmin.RunView) error {
        return audit(ev)          // 적용되어 우세
    }).
    OnDefeated(func(ctx toulmin.Context, ev toulmin.NodeEvent, view toulmin.RunView) error {
        return deny(ev)           // 적용됐으나 패배
    })

results, view, err := g.Run(ctx)  // []EvalResult, RunView, error
```

핸들러는 등록 순서로 발화하며, 첫 핸들러 에러에서 `Run`이 멈춘다. 어떤 핸들러가 발화하기
전에 `Run`은 모든 노드의 최종 이벤트를 담은 불변 스냅샷 **RunView**를 한 번 빌드해 모든
핸들러에 전달하고(반환값으로도 돌려준다), 핸들러는 자신의 이벤트와 함께 `view.All()`,
`view.Get(name)`, `view.Attackers(name)`로 그래프 전체의 최종 상태를 읽는다. `ctx`는 가변
(부수효과)이고, `view`는 불변 판정 스냅샷이다.

### 실행 합성 (execution composition)

`rule.Run(g)`는 노드가 **Active**일 때 같은 ctx로 하위 그래프 `g`를 Run하도록 선언한다 —
그래프의 그래프. 판정은 `Attacks`를 따라 *위로* 합성되고(verdict가 위로 흐른다), 실행은
`Run`을 따라 *아래로* 합성된다(하위 그래프의 verdict는 격리되고 에러만 전파된다). 실행 합성은
DAG여야 한다 — 사이클은 사전에 거부되고, 깊이는 64로 제한된다.

```go
order := g.Rule(orderPlaced)
order.OnActive(logOrder).Run(notifyGraph)   // Active order → notify 그래프 Run
```

Run 핸들러 + RunView 계열은 Go, TypeScript, Python 포트 모두에서 제공된다. 실행
합성(`rule.Run(g)`)은 현재 Go 전용이다.

## 프레임워크 패키지

코어 위에 도메인별 프레임워크를 제공한다. 규칙 함수와 래퍼가 미리 구현되어 있다.

| 패키지 | 도메인 | 핵심 API |
|---|---|---|
| `pkg/toulmin` | 코어 엔진 | `Graph`, `EvalOption` |
| `pkg/policy` | 접근 제어 (인증, 인가, IP, Rate limit) | `Guard` (net/http 미들웨어) |
| `pkg/state` | 상태 전이 (FSM) | `Machine.Can`, `Mermaid()` |
| `pkg/approve` | 다단계 결재 | `Flow.Evaluate` |
| `pkg/price` | 할인 판정 (쿠폰, 멤버십) | `Pricer.Evaluate` |
| `pkg/feature` | 피처 플래그 (롤아웃, 토글) | `Flags.IsEnabled` |
| `pkg/moderate` | 콘텐츠 모더레이션 (혐오, 스팸) | `Moderator.Review` |
| `pkg/tangl` | 마크다운 기반 정책 언어 | `parser.Parse`, `validate.Validate` |

프레임워크 없이 코어만 써도 된다. 위 킬러 예제처럼 직접 규칙 함수를 작성하는 게 가장 유연하다.

## 왜 toulmin인가

### vs if-else

- 규칙 추가: `new.Attacks(existing)` 한 줄 vs 중첩 구조 전체 리팩터링
- 예외 처리: 엣지 선언 vs 조건문 안에 조건문
- 판정 근거: `Evaluate(EvalOption{Trace: true})` 내장 vs 별도 로깅 구축
- 테스트: 규칙 함수 단위 테스트 vs 조합 폭발

### vs OPA/Casbin/Cedar

| | toulmin | OPA | Casbin | Cedar |
|---|---|---|---|---|
| 규칙 언어 | **TS/Python/Go 함수** | Rego (DSL) | PERM 모델 (설정) | Cedar (DSL) |
| 예외 처리 | **defeats graph** | 규칙 우선순위 | 정책 우선순위 | forbid/permit |
| 예외의 예외 | **Except** | 없음 | 없음 | 없음 |
| 판정 | **연속값 [-1,1]** | allow/deny | allow/deny | allow/deny |
| 판정 근거 | **Trace 내장** | Decision log | 없음 | 없음 |
| 의존성 | **없음** | Go + Rego 런타임 | Go | Rust + FFI |
| 학습 곡선 | 쓰는 언어만 알면 됨 | Rego 학습 필요 | PERM 모델 학습 | Cedar 문법 학습 |

### 학술적 기반

| 구성 요소 | 출처 |
|---|---|
| 6요소 논증 구조 | Toulmin (1958) |
| strict/defeasible/defeater | Nute (1994) |
| h-Categoriser | Amgoud & Ben-Naim (2013, 2017) |

## 테스트

`RunCases`로 테이블 드리븐 정책 테스트를 보일러플레이트 없이 작성한다:

```go
func TestAccessPolicy(t *testing.T) {
    g := buildAccessGraph()
    toulmin.RunCases(t, g, []toulmin.TestCase{
        {Name: "admin 허용",      Context: &Ctx{Role: "admin"},  Expect: toulmin.VerdictAbove(0)},
        {Name: "차단 IP 거부",     Context: &Ctx{IP: "blocked"},  Expect: toulmin.VerdictAtMost(0)},
        {Name: "미인증 거부",      Context: &Ctx{User: nil},      Expect: toulmin.NoResult},
        {Name: "부분 재정의",      Context: &Ctx{Role: "editor"}, Expect: toulmin.VerdictBetween(0, 0.5)},
    })
}
```

| Expectation | 조건 |
|---|---|
| `VerdictAbove(v)` | verdict > v |
| `VerdictAtMost(v)` | verdict <= v |
| `VerdictBetween(lo, hi)` | lo < verdict <= hi |
| `NoResult` | 활성 warrant 없음 |

## CLI

```bash
toulmin evaluate    # 예시 실행
```

## 사용 사례

- **[filefunc](https://github.com/park-jun-woo/filefunc)** — LLM 네이티브 Go 코드 구조 도구. `validate` 명령이 toulmin defeats graph로 규칙 예외(F5, F6 등)를 처리한다.

## 라이선스

MIT
