# toulmin

<p align="center">
  <img src="images/toulmin.webp" alt="toulmin" width="100%">
</p>

[![version](https://img.shields.io/badge/version-0.8.1-blue)](https://github.com/park-jun-woo/toulmin/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow)](https://opensource.org/licenses/MIT)

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
go get github.com/park-jun-woo/toulmin/go/pkg/toulmin  # Go
```

## 핵심 개념

### 규칙은 함수다

```typescript
// TypeScript — 수 v ∈ [-1, +1] 또는 boolean (true → +1, false → -1)
const fn: RuleFunc = (ctx, specs) => [number | boolean, unknown]

// Python — 실수 v ∈ [-1, +1] 또는 bool (True → +1, False → -1)
def fn(ctx: Context, specs: list[Spec]) -> tuple[float | bool, Any]: ...

// Go
func fn(ctx Context, specs Specs) (float64, any) // 등급 판정: v ∈ [-1, +1]
func fn(ctx Context, specs Specs) (bool, any)    // 그대로 허용: true → +1, false → -1

// C
bool rc_rule_func(...);     /* rc_graph_rule / _counter / _except로 등록 */
rc_verdict_func(... double *out_verdict); /* rc_graph_rule_verdict / _counter_verdict / _except_verdict로 등록 */
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

세 종류의 노드, 두 종류의 엣지:

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

**Supersedes는 크리스프한 우선권이다.** `Attacks`는 h-Categoriser에 가중치로 반영되지만 `winner.Supersedes(loser)`는 그렇지 않다. `winner`의 적용도 `a_eff > 0`이면 `loser`는 verdict가 낮아지는 게 아니라 집계에서 통째로 배제된다 — trace에 `Excluded`로 표시되고 순회 자체에 들어가지 않는다. 한 규칙이 다른 규칙과 저울질하는 게 아니라 그냥 우선해야 할 때 쓴다.

`a_eff`는 `winner` 자신의 `(v+1)/2`에 `Except` 게이트를 곱한 값이다(아래 참조). `Counter` 반박은 반영하지 않고, `Strict`인 winner는 `(v+1)/2` 그대로다. `a_eff`가 0.3처럼 일부만 남아도 loser 전체를 배제한다. 예: VIP 할인이 일반 할인을 배제한다. "정지된 VIP" 예외가 VIP 할인을 완전히 막으면 일반 할인은 배제되지 않아 고객은 일반 할인을 받는다. 예외 force가 0.5이거나 VIP 할인에 반박이 붙으면 배제는 유지된다. 배제된 규칙도 자기 배제권은 유지한다(`X > Y > Z`에서 `Y`가 배제돼도 `Z`는 배제된다). 네 포트 모두 지원한다(TS·Python `supersedes`, C `rc_rule_supersedes`).

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

### 등급 규칙, 예외, 공격 강도

규칙 함수는 예/아니오 대신 얼마나 해당하는지를 반환할 수 있다. 네 포트가 아래 의미론을 똑같이 따른다. 예시는 Go 이름으로 쓴다.

| | Go | TypeScript | Python | C |
|---|---|---|---|---|
| 등급 규칙 | `(float64, any)` | `[number, unknown]` | `(float, Any)` | `rc_verdict_func` + `rc_graph_*_verdict` |
| 공격 강도 | `Attacks(t, f)` | `attacks(t, f)` | `attacks(t, force=f)` | `rc_rule_attacks_with_force(r, t, f)` |
| 적용도 | `EvalResult.Applicability` | `applicability` | `applicability` | `rc_eval_result.applicability` |
| 규칙 자체의 v | `TraceEntry.RuleVerdict` | `ruleVerdict` | `rule_verdict` | `rc_trace_entry.rule_verdict` |
| 잘못된 force·qualifier | panic | 예외 | `ValueError` | `RC_ERR_INVALID_ARG` |
| 잘못된 규칙 verdict | `Evaluate` 오류 | 오류 | `RuntimeError` | `RC_ERR_RULE` |

TypeScript와 Python에서 수도 boolean도 아닌 반환값(`undefined`, `None`)은 "거짓"이 아니라 오류다.

```go
func fn(ctx Context, specs Specs) (float64, any) // v ∈ [-1, +1]
```

| `v` | 의미 |
|---|---|
| `+1` | 완전히 해당 (옛 `true`) |
| `0` | **반쯤 해당** |
| `-1` | 해당 없음 (옛 `false`) |

**해당 없음은 0이 아니라 −1이다. 0은 '반쯤 해당'이다.** bool에 익숙한 대로 "아니오"를 `return 0, nil`로 쓰면 규칙이 꺼지지 않는다. 적용도 0.5로 Active인 채 남고, 공격자도 평가되고, 그래프에 계속 반영된다. 판정이 조용히 틀어지며, **엔진은 이를 검출하지 못한다 — 0은 합법적인 값이므로 사용자가 알아서 조심해야 한다.** 해당하지 않으면 `-1`을 반환한다.

- `v`는 `[-1, +1]` 안이어야 한다. NaN이나 범위를 벗어난 값은 오류다. `Evaluate`/`Run`이 멈추고 `rule "<id>": rule verdict <v> out of range [-1, 1]`를 반환한다(규칙 panic과 같은 경로).
- 옛 `func(Context, Specs) (bool, any)` 시그니처도 그대로 쓴다. `true` → `+1`, `false` → `-1`이고, bool 그래프의 판정값은 이전과 정확히 같다(`Except`가 공격하는 경우만 다르다 — 아래).
- `v > -1`이면 Active다. `v == -1`일 때만 공격자 평가를 건너뛴다(지연 평가).
- `TraceEntry.RuleVerdict`에 규칙 자체의 `v`가 남는다. `TraceEntry.Activated`는 `v > -1`로 유지된다.

반박(Counter)과 예외(Except)는 이제 다르게 계산한다:

```
a      = (v + 1) / 2                               자기 적용도
S      = Σ_반박  F · (v_j + 1) / 2                 반박: 분모
a_eff  = a · Π_예외 (1 − F · (v_j + 1) / 2)        예외: 곱셈 게이트
raw    = q · a_eff / (1 + S)
verdict = 2 · raw − 1
```

- **Counter**(반박)는 결론에 맞서는 주장이다. 지금처럼 분모에 들어간다.
- **Except**(예외)는 "이 경우엔 규칙이 적용되지 않는다"는 말이다. 분모에 더하지 않고 대상의 적용도를 깎는다. 완전히 해당하는 예외 하나면 대상은 `-1`이 된다. `v = 0.5`인 예외 셋이면 적용도가 `(1 − 0.75)³ ≈ 0.016`만 남는다. (이전에는 예외 셋이 붙어도 `-0.5`까지만 내려갔다.)
- `Strict` 규칙은 둘 다 무시한다. `Supersedes`(크리스프 배제)는 배제자의 `a_eff > 0`일 때만 배제하므로, 예외가 배제자를 완전히 막으면 배제도 풀린다.
- `EvalResult.Applicability`가 각 warrant의 `a_eff`를 알려 준다.

`Attacks`는 선택 인자로 공격 강도 `F ∈ [0, 1]`을 받는다. 기본값은 `1`이다:

```go
blocked.Attacks(auth)       // F = 1 (기존 동작 그대로)
hint.Attacks(auth, 0.3)     // 약한 반박
partial.Attacks(rule, 0.5)  // 적용도를 최대 절반만 깎는 예외
```

`F`는 반박 항과 예외 게이트 양쪽에 곱해진다. 값을 둘 이상 주거나 NaN이거나 `[0, 1]`을 벗어나면 선언 시점에 panic한다(`Qualifier`와 같은 관례).

### Level

연속값 verdict와 직교하는 축으로, 각 규칙은 크리스프한 리포팅 `Level`을 가진다 — `Unset` / `Fail` / `Review` — 하류 집계용이다(예: `Fail`이 하나라도 있으면 FAIL, 없고 `Review`가 있으면 REVIEW, 아니면 PASS). 기본값으로 `Counter`는 `Fail`, `Except`는 `Review`를 보고하며, `rule.Level(l)`로 이를 재정의한다(`Unset`으로 되돌리는 것도 가능). `TraceEntry.Level`로 노출된다.

```go
overdrawn := g.Counter(isOverdrawn).Level(toulmin.Fail)
```

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
    fmt.Printf("%s role=%s activated=%v level=%v excluded=%v evidence=%v\n",
        t.Name, t.Role, t.Activated, t.Level, t.Excluded, t.Evidence)
}
```

`Excluded`는 크리스프하게 배제(superseded)된 노드를 표시하고, `Level`은 그 크리스프 리포팅 레벨이다(위 참조). `Evidence`는 규칙 함수의 두 번째 반환값이다 — 위치·수량이 명시된 위반을 담으려면 `Fact{Where, Expected, Actual}`를 반환하고 하류에서 `t.Evidence.(toulmin.Fact)`로 회수한다.

모더레이션 로그, 감사 추적, 디버깅 — 별도 로깅 없이 엔진이 제공한다.

## Run

`Evaluate`는 판정하고 반환만 한다 — 순수하고 멱등하다. `Run`은 먼저 판정한 뒤 **행동한다**:
전체 패스를 수행하고, **Active** 노드마다 핸들러를 하나 발화한다. 노드는 규칙이 적용되고
**우세**할 때(`activated && verdict > 0`) Active이며, 이벤트는 이 하나뿐이다. 다른 노드의
결과가 궁금하면 반환된 `trace`를 읽는다: `verdict > 0` 우세, `verdict <= 0` 패배,
`activated == false` 미적용.

```go
g.Rule(isAuthenticated).
    RunOn(func(self toulmin.TraceEntry, t toulmin.Trace) error {
        return audit(self)                  // self = 이 핸들러가 달린 자기 노드 (적용되어 우세)
    })

results, trace, err := g.Run(ctx)  // []EvalResult, Trace, error
```

핸들러는 등록 순서로 발화하며, 첫 핸들러 에러에서 `Run`이 멈춘다. 핸들러는 자기 노드의 entry인
`self`와 읽기 전용 `Trace` `t`를 받는다 — `t`는 `Run`이 반환하는 값과 동일하다 — `Trace` 메서드는
정확히 셋이다: `t.All()`(전 노드의 `TraceEntry`, 등록 순서), `t.Get(name)`(짧은 이름으로 한 노드
조회), `t.Ctx()`(이 Run의 컨텍스트). `self`는 첫 인자로 엔진이 직접 건넨다 — 자기 줄을 이름으로
다시 찾을 필요가 없다(entry 이름에 동적 접미사가 붙으면 `t.Get("X")` 자기-조회가 취약했다). 각 `TraceEntry`는 **Claim**(`Name`), **Ground**(`Ground` = ctx
그대로), **Backing**(`Specs`), **Verdict**를 노출하므로 별도 뷰 없이 감사·설명·gradient 임계값
판정이 가능하다. 공격자 직접 조회 API는 없다 — `t.All()[i].Verdict`로 추론한다. `ctx`는
가변(부수효과 전파)이고, trace가 판정 기록이다. trace를 JSON으로 내보낸다면 `ctx`(= `Ground`)에
직렬화 가능한 값만 담아야 한다.

### 실행 합성 (execution composition)

`rule.Run(g)`는 노드가 **Active**일 때 같은 ctx로 하위 그래프 `g`를 Run하도록 선언한다 —
그래프의 그래프. 판정은 `Attacks`를 따라 *위로* 합성되고(verdict가 위로 흐른다), 실행은
`Run`을 따라 *아래로* 합성된다(하위 그래프의 verdict는 격리되고 에러만 전파된다). 실행 합성은
DAG여야 한다 — 사이클은 사전에 거부되고, 깊이는 64로 제한된다.

```go
order := g.Rule(orderPlaced)
order.RunOn(logOrder).Run(notifyGraph)   // Active order → notify 그래프 Run
```

`RunOn` 핸들러 + `Run(g)` 실행 합성은 Go, TypeScript, Python 포트 모두에서 제공된다.

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
| `pkg/tangl` | TANGEUL 마크다운 정책 언어: 한·영 코덱, 검증기, 평가기, Go 코드 생성, effect 분석, saga 런타임 | `parser.Parse`, `eval.CompileDocument`, `gen.GenerateWithOptions` |

프레임워크 없이 코어만 써도 된다. 위 킬러 예제처럼 직접 규칙 함수를 작성하는 게 가장 유연하다.

## TANGEUL

TANGEUL은 코어 위에 바로 컴파일되는 마크다운 정책 언어다 — 별도 런타임도, 수작업 접착 코드도 없다.

- **마크다운이 곧 실행 코드다.** `tangl:Cases` 블록(노드, `don't` 공격 엣지, `do`/`undo`/`run` 실행 엣지)이 `toulmin.Graph`로 그대로 컴파일된다. 손으로 옮겨 적는 부분이 없다.
- **규칙은 판정을 선언하고, 사안은 그래프 역할을 바인딩한다.** `tangl:규칙`은 판정식 또는 호스트 함수와 기준·근거를 소유한다. `tangl:사안`은 동명 규칙을 인용해 일반·반박·예외 역할, 가중치, 보고 레벨만 부여한다.
- **목록 기준도 문서 안에 둔다.** `` `도메인` 은 ["a.com", "b.com"] 을 뜻한다 ``는 정확 일치 목록 Spec을 만든다. 파일 목록은 `` `도메인` 은 `JSON 목록` 을 "./domains.json" 파일에서 읽는다 ``로 선언하며, 소비자 loader가 문서 기준 경로에서 eval/생성 전에 한 번 스냅샷한다.
- **판정과 실행이 분리된다.** `Provides`/`Internal` 항목은 케이스를 `check`(순수 판정 — `Graph.Evaluate`로 컴파일)하거나 `run`(판정 후 `RunOn`을 통한 `do`/`undo` 부수효과 — `Graph.Run`으로 컴파일)한다.
- **saga 보상이 내장되어 있다.** `do`가 성공하는 순간 짝을 이루는 `undo`가 보상 스택에 적재되고, 같은 패스 안에서 이후 단계가 실패하면 스택을 LIFO로 되돌린다. 복구 불가능한 실패는 원인과 보상 실패를 함께 담은 `*tangl.ReviewError`로 승격된다.
- **effect summary로 감사한다.** `tangl effects`는 실행하지 않고도 한 엔드포인트가 일으킬 수 있는 모든 `do`/`undo`를 정적으로 나열한다.
- **공격 세기도 표면에서.** `` `Y` 면 `X` 무효다 ``는 세기 1의 공격이고, `` `Y` 면 `X` 를 0.5의 세기로 공격한다 ``(영어 `` `Y` attacks `X` with force 0.5 ``)는 `Y.Attacks(X, 0.5)`로 컴파일된다. 세기는 `[0, 1]`(소수 넷째 자리까지) 값이며 `.geul`에는 ATTACKS의 `F_FORCE` 자리로 실린다.
- **크리스프 우선권과 레벨도 표면에서.** `` `X` 면 `Y` 배제한다 ``는 `X`의 실효 적용도가 `a_eff > 0`인 동안 `Y`를 집계에서 배제한다. `` `A` > `B` > `C` 순서로 배제한다 ``는 모든 앞→뒤 우선권 간선으로 전개된다. 관형절 `실패를 보고하는` / `검토를 보고하는`은 크리스프 리포팅 레벨을 표기한다.
- **하나의 이진 단일 진실(TANGEUL).** `tangl encode`는 저작한 `.md` 표면을 로케일 독립적인 컴팩트 `.geul` 파일(GEUL t002 DOC 하나, TOULMIN 코덱: AST + 원문 prose span)로 흡수하며, 이후 모든 명령은 그 스트림 위에서 동작한다. `tangl decode`는 스트림을 다시 텍스트 표면으로 렌더링한다 — 모드 1에서는 바이트 항등, 모드 2에서는 다른 로케일로 재생성.
- **영어·한국어 표면, 같은 AST.** `tangl encode`는 문서의 섹션 헤더로 en/ko를 자동 판별하고, `tangl decode --locale ko`는 임의의 `.geul`을 비개발자 감사자를 위한 한국어 감사 표면으로 렌더링한다(`TANGEUL.ko.md` 참조).

```bash
go install github.com/park-jun-woo/toulmin/go/cmd/tangl@latest

tangl encode policy.md -o policy.geul   # .md 표면을 .geul 스트림으로 흡수 (로케일 자동 판별)
tangl check policy.geul                 # 검증 + lint
tangl effects policy.geul pay           # 엔드포인트 "pay"에서 도달 가능한 정적 do/undo 요약
tangl gen policy.geul -o pay_gen.go     # 문서를 Go 소스로 컴파일
tangl decode policy.geul                # 스트림을 다시 텍스트로 렌더링 (모드 1, 바이트 항등)
tangl decode policy.geul --locale ko    # 같은 스트림을 한국어 감사 표면으로 렌더링 (모드 2)
```

```
$ tangl encode pkg/tangl/parser/testdata/transfer.md -o transfer.geul

$ tangl check transfer.geul
ok

$ tangl effects transfer.geul transfer
do   bank.withdraw        once (can withdraw / balance sufficient)
undo bank.refund               (can withdraw / balance sufficient)
do   bank.deposit         once (can deposit / recipient valid)
do   log.TransferComplete      (can deposit / recipient valid)
```

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
