# Phase 010: Backing을 일급 시민으로 — 코어 확장 (구현 완료)

## 목표

Toulmin 6-element 구조의 Backing을 엔진 코어에 일급 시민으로 포함시킨다.
판정 기준(화이트리스트, 임계값, role명 등)을 클로저에 숨기지 않고, 엔진이 관리하는 명시적 값으로 선언한다.
EvaluateTrace에서 backing 값까지 추적 가능하게 한다.

## 배경

### 제1원칙: Toulmin 6-element 구조

```
1. Claim      — 주장 (뭘 판정하나)
2. Ground     — 근거 (판정 대상의 사실)
3. Warrant    — 보증 (판정 규칙)
4. Backing    — 뒷받침 (규칙의 판정 기준)
5. Qualifier  — 한정어 (확신도)
6. Rebuttal   — 반박
```

### 현재 문제

1. **Backing이 주석이다**: `//tm:backing "Böhm-Jacopini theorem"`은 사람이 읽는 메타데이터일 뿐, 엔진이 모른다
2. **판정 기준이 클로저에 숨어있다**: `HasRole("admin")`, `IsWhitelisted(whitelist)` 등 클로저가 backing을 캡처하지만, 엔진은 이 값을 볼 수 없다
3. **EvaluateTrace가 "왜"를 절반만 설명한다**: "IsWhitelisted=true"는 보여주지만, "화이트리스트가 `[10.0.0.1]`이었기 때문"은 보여주지 못한다
4. **클로저 funcID 문제**: 클로저를 매번 생성하면 funcID가 달라져서 변수에 저장 후 재사용해야 한다. backing이 일급 시민이면 이 문제가 사라진다

### ground와 backing의 구분

| | ground | backing |
|---|---|---|
| Toulmin 역할 | 판정 대상의 사실 | 규칙의 판정 기준 |
| 언제 결정 | 평가 시점 (요청마다 다름) | 선언 시점 (graph 구성 시 고정) |
| 예시 | User, ClientIP, 현재 상태 | 화이트리스트, role명, 임계값 |
| 코드에서 | `g.Evaluate(claim, ground)` | `Warrant(fn, backing, qualifier)` |

## 핵심 설계

### Rule 함수 시그니처 확장

```go
// 현재: func(claim any, ground any) (bool, any)
// 확장: func(claim any, ground any, backing any) (bool, any)
```

backing이 nil이면 기존과 동일하게 동작한다. 하위 호환성을 유지한다.

### GraphBuilder API 변경

backing은 Warrant/Rebuttal/Defeater의 인자로 직접 전달한다. 체이닝이 아니라 하나의 호출에 Toulmin 구조가 완전히 표현된다.

**인자 순서: `fn, backing, qualifier`**

- fn — 무슨 규칙인가
- backing — 이 규칙의 판정 기준은 무엇인가
- qualifier — 이 규칙의 확신도는 얼마인가

backing이 없는 rule은 명시적으로 `nil`을 전달한다. Toulmin 6-element 중 backing이 없다는 것 자체가 의미있는 선언이다.

```go
g := toulmin.NewGraph("admin").
    Warrant(IsInRole, "admin", 1.0).
    Warrant(IsAuthenticated, nil, 1.0).
    Rebuttal(IsIPBlocked, blocklist, 1.0).
    Defeater(IsWhitelisted, whitelist, 1.0).
    Defeat(IsIPBlocked, IsAuthenticated).
    Defeat(IsWhitelisted, IsIPBlocked)
```

### 같은 함수 + 다른 backing = 다른 rule

backing이 일급 시민이면 같은 함수를 다른 backing으로 재사용할 수 있다:

```go
g := toulmin.NewGraph("multi-role").
    Warrant(IsInRole, "admin", 1.0).
    Warrant(IsInRole, "editor", 1.0)
```

이 경우 funcID만으로는 rule을 구분할 수 없으므로, **funcID + backing 조합**으로 rule을 식별한다. 클로저가 필요 없어진다.

### rule 함수 예시: 클로저 → 순수 함수

```go
// 현재 (클로저 방식)
func HasRole(role string) func(any, any) (bool, any) {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*RouteContext)
        return ctx.User.Role == role, nil
    }
}
// 사용: hasAdmin := HasRole("admin")
//       g.Warrant(hasAdmin, 1.0)

// 확장 (backing 방식)
func IsInRole(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RouteContext)
    role := backing.(string)
    return ctx.User.Role == role, nil
}
// 사용: g.Warrant(IsInRole, "admin", 1.0)
```

클로저가 사라지고, backing이 명시적으로 드러난다.

### Warrant/Rebuttal/Defeater 시그니처

```go
func (b *GraphBuilder) Warrant(fn func(any, any, any) (bool, any), backing any, qualifier float64) *GraphBuilder

func (b *GraphBuilder) Rebuttal(fn func(any, any, any) (bool, any), backing any, qualifier float64) *GraphBuilder

func (b *GraphBuilder) Defeater(fn func(any, any, any) (bool, any), backing any, qualifier float64) *GraphBuilder
```

### RuleMeta 확장

```go
type RuleMeta struct {
    Name      string
    Qualifier float64
    Strength  Strength
    Defeats   []string
    Backing   any            // 신규: backing 값
    Fn        func(claim any, ground any, backing any) (bool, any)
}
```

### TraceEntry 확장

```go
type TraceEntry struct {
    Name      string  `json:"name"`
    Role      string  `json:"role"`
    Activated bool    `json:"activated"`
    Qualifier float64 `json:"qualifier"`
    Evidence  any     `json:"evidence,omitempty"`
    Backing   any     `json:"backing,omitempty"`   // 신규: 판정 기준
}
```

EvaluateTrace 결과에 backing이 포함되므로, "왜 이 규칙이 이 결과를 냈는가"를 완전히 설명할 수 있다:

```json
{
    "name": "IsInRole",
    "role": "warrant",
    "activated": true,
    "qualifier": 1.0,
    "backing": "admin",
    "evidence": null
}
```

### Rule 식별: funcID + backing

같은 함수를 다른 backing으로 등록할 수 있으므로, rule 식별자를 변경한다:

```go
// 현재: funcID(fn) → "github.com/example/pkg.IsInRole"
// 확장: funcID(fn) + "#" + fmt.Sprint(backing) → "github.com/example/pkg.IsInRole#admin"
```

backing이 nil이면 기존과 동일하게 funcID만 사용한다.

### Defeat에서의 함수 참조

Defeat는 함수 포인터로 참조하므로, backing이 다른 같은 함수를 구분해야 한다. 두 가지 방법:

**방법 A: Defeat도 backing 명시**

```go
g.Defeat(IsIPBlocked, IsAuthenticated)  // backing nil끼리 매칭
```

backing이 nil인 경우는 기존과 동일. backing이 있는 rule을 Defeat에서 참조할 때는 별도 API가 필요할 수 있다:

```go
// 같은 함수 + 다른 backing을 가진 rule을 Defeat에서 구분
g.DefeatWith(IsInRole, "admin", IsInRole, "editor")
```

**방법 B: Warrant/Rebuttal 등록 순서로 내부 식별 (심플)**

Defeat에서는 함수 포인터만 사용하되, backing이 같은 함수가 여러 개면 등록 순서로 구분한다. 단, 같은 함수 + 같은 역할이 아닌 이상 충돌이 거의 없다.

→ **방법 A 채택**: 명시적이 Toulmin 원칙에 부합한다.

### evalContext 확장

```go
type evalContext struct {
    fnMap       map[string]func(any, any, any) (bool, any)  // backing 포함
    qualMap     map[string]float64
    strMap      map[string]Strength
    backingMap  map[string]any                               // 신규
    edges       map[string][]string
    attackerSet map[string]bool
    ran         map[string]bool
    active      map[string]bool
    evidence    map[string]any
    trace       []TraceEntry
    roleMap     map[string]string
}
```

### calc 확장

```go
// 현재
fn(claim, ground)

// 확장
fn(claim, ground, ctx.backingMap[name])
```

### 하위 호환성

기존 `func(claim any, ground any) (bool, any)` 시그니처도 지원한다.
내부에서 래핑하여 새 시그니처로 변환한다:

```go
// 기존 시그니처를 내부에서 래핑
func wrapLegacy(fn func(any, any) (bool, any)) func(any, any, any) (bool, any) {
    return func(claim any, ground any, backing any) (bool, any) {
        return fn(claim, ground)
    }
}
```

기존 코드가 깨지지 않도록 Warrant/Rebuttal/Defeater가 두 시그니처를 모두 받을 수 있어야 한다. Go의 `any` 타입으로 함수 인자를 받아 내부에서 시그니처를 판별한다:

```go
func (b *GraphBuilder) Warrant(fn any, backing any, qualifier float64) *GraphBuilder {
    switch f := fn.(type) {
    case func(any, any, any) (bool, any):
        // 새 시그니처 — 그대로 사용
    case func(any, any) (bool, any):
        // 기존 시그니처 — wrapLegacy
    }
}
```

### 사용 예시 비교

```go
// 현재 (클로저 방식) — Phase 011 route에서
blocklist := func(ip string) bool { return ip == "1.2.3.4" }
whitelist := func(ip string) bool { return ip == "10.0.0.1" }
hasAdmin := route.HasRole("admin")
ipBlocked := route.IsIPBlocked(blocklist)
whitelisted := route.IsWhitelisted(whitelist)

g := toulmin.NewGraph("route:admin").
    Warrant(route.IsAuthenticated, 1.0).
    Warrant(hasAdmin, 1.0).
    Rebuttal(ipBlocked, 1.0).
    Defeater(whitelisted, 1.0).
    Defeat(ipBlocked, route.IsAuthenticated).
    Defeat(whitelisted, ipBlocked)

// 확장 (backing 방식)
blocklist := func(ip string) bool { return ip == "1.2.3.4" }
whitelist := func(ip string) bool { return ip == "10.0.0.1" }

g := toulmin.NewGraph("route:admin").
    Warrant(route.IsAuthenticated, nil, 1.0).
    Warrant(route.IsInRole, "admin", 1.0).
    Rebuttal(route.IsIPInList, blocklist, 1.0).
    Defeater(route.IsIPInList, whitelist, 1.0).
    Defeat(route.IsIPInList, route.IsAuthenticated).   // blocklist backing
    Defeat(route.IsIPInList, route.IsIPInList)          // whitelist defeats blocklist — DefeatWith 필요
```

위 예시에서 같은 `IsIPInList` 함수가 blocklist/whitelist 두 backing으로 사용되므로 `DefeatWith`가 필요하다:

```go
g.DefeatWith(route.IsIPInList, whitelist, route.IsIPInList, blocklist)
```

## 범위

### 포함

1. **RuleMeta 확장**: Backing 필드, Fn 시그니처 `(claim, ground, backing)` 변경
2. **GraphBuilder 변경**: Warrant/Rebuttal/Defeater 시그니처 `(fn, backing, qualifier)`
3. **DefeatWith 추가**: 같은 함수 + 다른 backing 구분
4. **Rule 식별자 변경**: funcID + backing 조합
5. **evalContext 확장**: backingMap 추가, calc에서 backing 전달
6. **TraceEntry 확장**: Backing 필드 추가
7. **하위 호환성**: 기존 `(claim, ground)` 시그니처 래핑
8. **기존 테스트 수정**: 새 시그니처 적용
9. **신규 테스트**: backing 기반 rule 테스트

### 제외

- 프레임워크 패키지 (pkg/route 등) — Phase 011에서 backing 적용
- CLI/YAML 코드젠 — 별도 Phase에서 backing 지원

## 산출물

```
pkg/
  toulmin/
    rule_meta.go                  — Backing 필드, Fn 시그니처 변경
    graph_builder.go              — GraphBuilder 구조체 (변경 없음)
    graph_builder_warrant.go      — Warrant(fn, backing, qualifier) 변경
    graph_builder_rebuttal.go     — Rebuttal(fn, backing, qualifier) 변경
    graph_builder_defeater.go     — Defeater(fn, backing, qualifier) 변경
    graph_builder_defeat.go       — Defeat 기존 유지
    graph_builder_defeat_with.go  — DefeatWith(fn1, b1, fn2, b2) 신규
    func_id.go                    — funcID + backing 조합 식별자
    rule_id.go                    — ruleID(fn, backing) 헬퍼 (신규)
    trace_entry.go                — Backing 필드 추가
    eval_context.go               — backingMap 추가
    new_eval_context.go           — backingMap 초기화
    eval_context_calc.go          — backing 전달
    eval_context_calc_trace.go    — backing 전달 + trace 기록
    wrap_legacy.go                — 기존 시그니처 래핑 (신규)
    eval_result.go                — 변경 없음
    *_test.go                     — 기존 테스트 수정 + 신규 테스트
```

## 단계

### Step 1: 시그니처 변경

- `RuleMeta.Fn`: `func(any, any) (bool, any)` → `func(any, any, any) (bool, any)`
- `RuleMeta.Backing`: `any` 필드 추가
- `wrap_legacy.go` 신규: 기존 시그니처를 새 시그니처로 래핑

### Step 2: Rule 식별자 변경

- `rule_id.go` 신규: `ruleID(fn, backing)` — funcID + "#" + fmt.Sprint(backing)
- backing이 nil이면 funcID만 사용
- 기존 `funcID` 호출부를 `ruleID`로 변경

### Step 3: GraphBuilder 변경

- Warrant/Rebuttal/Defeater: `(fn any, backing any, qualifier float64)` 시그니처
- fn을 `any`로 받아 내부에서 시그니처 판별 (하위 호환)
- backing + funcID로 ruleID 생성하여 등록

### Step 4: DefeatWith 추가

- `DefeatWith(fromFn any, fromBacking any, toFn any, toBacking any) *GraphBuilder`
- 기존 `Defeat(from, to)`는 backing nil인 rule 간 매칭 (하위 호환)

### Step 5: evalContext 확장

- `backingMap map[string]any` 추가
- `newEvalContext`에서 backingMap 초기화
- `calc`/`calcTrace`에서 `fn(claim, ground, backingMap[name])` 호출

### Step 6: TraceEntry 확장

- `Backing any` 필드 추가
- `calcTrace`에서 trace 기록 시 backing 포함

### Step 7: 기존 테스트 수정

- 기존 rule 함수 시그니처 변경 또는 wrapLegacy 적용
- Warrant/Rebuttal/Defeater 호출에 backing 인자 추가 (nil)
- 기존 테스트 전부 PASS 확인

### Step 8: 신규 테스트

- 같은 함수 + 다른 backing → 다른 rule로 식별
- backing이 EvaluateTrace 결과에 포함
- backing nil → 기존 동작과 동일
- 기존 `(claim, ground)` 시그니처 → wrapLegacy → 정상 동작
- DefeatWith로 같은 함수 + 다른 backing 간 defeat

### Step 9: pkg/route 테스트 수정

- pkg/route의 기존 rule 함수와 테스트를 새 시그니처에 맞춰 수정
- `go test ./...` 전체 PASS 확인

## 검증 기준

1. `Warrant(fn, "admin", 1.0)` 으로 backing을 인자로 전달할 수 있다
2. `Warrant(fn, nil, 1.0)` 으로 backing 없음을 명시적으로 선언할 수 있다
3. 같은 함수 + 다른 backing이 다른 rule로 식별된다
4. rule 함수가 `(claim, ground, backing)` 으로 호출되어 backing 값을 받는다
5. EvaluateTrace 결과의 TraceEntry에 Backing 필드가 포함된다
6. backing이 nil이면 기존과 동일하게 동작한다
7. 기존 `(claim, ground)` 시그니처가 wrapLegacy를 통해 정상 동작한다
8. DefeatWith로 같은 함수 + 다른 backing을 가진 rule 간 defeat가 동작한다
9. 클로저 없이 `IsInRole` + backing `"admin"` 으로 role 판정이 가능하다
10. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (현재 구조)

## 영향 범위

이 Phase 이후 Phase 011-017의 모든 프레임워크가 영향을 받는다:
- 클로저 rule 팩토리 → 순수 함수 + backing 인자 방식으로 전환
- 클로저 funcID 문제 근본 해소
- EvaluateTrace의 설명력 향상 (backing 값까지 추적)
