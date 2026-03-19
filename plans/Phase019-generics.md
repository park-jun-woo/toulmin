# Phase 019: 제네릭 전환 — Graph[C, G, B] 타입 안전 API

## 목표

코어 엔진과 프레임워크 패키지의 `any` 타입을 Go 제네릭으로 전환한다. `Graph[C, G, B]`로 claim/ground/backing을 타입화하고, 프레임워크의 도메인 타입(User, Approver)을 제네릭 파라미터로 노출하여 컴파일 타임 타입 안전성을 확보한다. evidence만 `any`를 유지한다.

## 배경

### 현재 문제

1. **모든 rule 함수가 타입 단언 필수**: `ground.(*RequestContext)`, `backing.(*RoleBacking)` 등 런타임 패닉 위험
2. **Graph API가 `any`로 열려 있음**: `g.Warrant(anyFunc, anyBacking, 1.0)` — 잘못된 함수 시그니처를 넣어도 컴파일 에러 없음, 런타임에 `toRuleFunc`에서 패닉
3. **Evaluate 호출도 `any`**: `g.Evaluate(claim, ground)` — claim/ground 타입 불일치를 컴파일러가 잡지 못함
4. **프레임워크 backing의 추출 함수가 `func(any) string`**: 잘못된 타입을 넘겨도 컴파일 에러 없음

### `any` 사용 현황 (150+ 건)

| 위치 | `any` 사용 | 제네릭화 |
|---|---|---|
| Rule 함수 시그니처 claim, ground | `func(any, any, any)` | C, G로 전환 |
| Rule 함수 시그니처 backing | rule마다 다른 타입 | B로 전환 — 동질 그래프는 concrete 타입, 이종 그래프는 `any` 폴백 |
| Rule 함수 반환 evidence | rule마다 다른 타입 | **any 유지** — EvalResult 슬라이스에 이종 혼재, 타입화 불가 |
| Graph.Evaluate(claim, ground) | claim, ground | C, G로 전환 |
| Graph.Warrant(fn, backing, q) | fn, backing | fn은 `func(C, G, B)`, backing은 B |
| Context 구조체 (User, Approver, Resource) | 도메인 타입 | 프레임워크 제네릭화로 전환 |
| Backing 추출 함수 `func(any) string` | 도메인 타입 | 프레임워크 제네릭화로 전환 |
| EvalResult.Evidence, TraceEntry.Evidence | 이종 | **any 유지** |

### 설계 원칙

- **C, G, B를 타입 파라미터로 노출한다** — claim/ground는 그래프 단위 동질, backing은 동질 그래프에서 완전 타입화, 이종 그래프에서 `any` 폴백
- **evidence만 `any`를 유지한다** — `[]EvalResult`의 각 원소가 서로 다른 evidence 타입을 반환하므로 단일 타입 파라미터 불가
- **프레임워크 패키지는 도메인 타입을 제네릭 파라미터로 노출한다** — User, Approver, Resource 등
- **Engine(레거시)은 변경하지 않는다** — 기존 string 기반 API 호환 유지

### B를 포함하는 이유

**동질 backing 그래프에서 backing 타입 단언이 완전히 제거된다:**

```go
// Graph[C, G, B] — B = *DiscountBacking
g := toulmin.NewGraph[*PurchaseRequest, *PriceContext, *DiscountBacking]("discount")

// rule 함수에서 타입 단언 제로 — backing이 이미 *DiscountBacking
func HasCoupon(claim *PurchaseRequest, ground *PriceContext, backing *DiscountBacking) (bool, any) {
    return claim.CouponCode == backing.Name, backing.Name
}
```

**프레임워크 backing 인터페이스로 등록 시점 검증:**

```go
// policy 패키지가 마커 인터페이스 정의
type Backing interface{ policyBacking() }
func (*RoleBacking) policyBacking()  {}
func (*IPListBacking) policyBacking() {}

// 잘못된 backing 타입을 넘기면 컴파일 에러
g := toulmin.NewGraph[any, *RequestContext, policy.Backing]("admin")
```

**이종 그래프는 B = any 폴백 — 손해 없음:**

```go
// B = any — Graph[C, G]와 동일한 수준, 타입 파라미터만 하나 더 명시
g := toulmin.NewGraph[any, *RequestContext, any]("admin")
```

### E를 제외하는 이유

evidence는 B와 근본적으로 다르다:

- **backing은 선언 시 고정** — graph 구성 시 Warrant/Rebuttal에 넘기는 값. 그래프 소유자가 타입을 안다
- **evidence는 실행 시 생성** — rule 함수가 반환하는 값. `[]EvalResult`의 각 원소가 서로 다른 evidence 타입
- `EvalResult` 슬라이스에 `string`, `float64`, `*struct` evidence가 혼재 — 단일 E 타입 파라미터로 통일 불가
- evidence 소비자는 rule 작성자와 동일인 — 자기가 넣은 타입을 자기가 꺼내므로 타입 단언 위험 낮음

## 핵심 설계

### 코어 엔진 — `Graph[C, G, B]`

```go
// pkg/toulmin/graph.go
type Graph[C, G, B any] struct {
    name    string
    rules   []*Rule[C, G, B]
    edges   []defeatEdge[C, G, B]
}

func NewGraph[C, G, B any](name string) *Graph[C, G, B]
```

```go
// pkg/toulmin/rule.go
type Rule[C, G, B any] struct {
    id        string
    fn        func(C, G, B) (bool, any)   // evidence만 any
    backing   B
    qualifier float64
    strength  Strength
}
```

```go
// pkg/toulmin/graph_warrant.go
func (g *Graph[C, G, B]) Warrant(fn func(C, G, B) (bool, any), backing B, qualifier float64) *Rule[C, G, B]

// graph_rebuttal.go, graph_defeater.go 동일 패턴
func (g *Graph[C, G, B]) Rebuttal(fn func(C, G, B) (bool, any), backing B, qualifier float64) *Rule[C, G, B]
func (g *Graph[C, G, B]) Defeater(fn func(C, G, B) (bool, any), backing B, qualifier float64) *Rule[C, G, B]
```

```go
// pkg/toulmin/graph_defeat.go
func (g *Graph[C, G, B]) Defeat(from *Rule[C, G, B], to *Rule[C, G, B])
```

```go
// pkg/toulmin/graph_evaluate.go
func (g *Graph[C, G, B]) Evaluate(claim C, ground G) ([]EvalResult, error)
func (g *Graph[C, G, B]) EvaluateTrace(claim C, ground G) ([]EvalResult, error)
```

```go
// pkg/toulmin/eval_result.go — evidence만 any 유지
type EvalResult struct {
    Name     string       `json:"name"`
    Verdict  float64      `json:"verdict"`
    Evidence any          `json:"evidence,omitempty"`   // any 유지
    Trace    []TraceEntry `json:"trace"`
}

type TraceEntry struct {
    Name      string  `json:"name"`
    Role      string  `json:"role"`
    Activated bool    `json:"activated"`
    Qualifier float64 `json:"qualifier"`
    Evidence  any     `json:"evidence,omitempty"`   // any 유지
    Backing   any     `json:"backing,omitempty"`    // B → any로 직렬화 (JSON 호환)
}
```

### 프레임워크 — 도메인 타입 제네릭화

#### policy — `U` (User), `B = policy.Backing`

```go
// pkg/policy/backing.go — 마커 인터페이스
type Backing interface{ policyBacking() }

func (*RoleBacking[U]) policyBacking()  {}
func (*OwnerBacking[U]) policyBacking() {}
func (*IPListBacking) policyBacking()   {}
func (NilBacking) policyBacking()       {}

// nil backing 대체용
type NilBacking struct{}
```

```go
// pkg/policy/request_context.go
type RequestContext[U any] struct {
    User     U
    ClientIP string
    Headers  map[string]string
    Metadata map[string]any
}
```

```go
// pkg/policy/role_backing.go
type RoleBacking[U any] struct {
    Role     string
    RoleFunc func(U) string
}
```

```go
// pkg/policy/owner_backing.go
type OwnerBacking[U any] struct {
    UserIDFunc     func(U) string
    ResourceIDFunc func(any) string
}
```

```go
// pkg/policy/rule_is_authenticated.go
func IsAuthenticated[U any](claim any, ground *RequestContext[U], backing Backing) (bool, any) {
    var zero U
    return any(ground.User) != any(zero), nil
}
```

```go
// pkg/policy/rule_is_in_role.go
func IsInRole[U any](claim any, ground *RequestContext[U], backing Backing) (bool, any) {
    rb := backing.(*RoleBacking[U])
    role := rb.RoleFunc(ground.User)
    return role == rb.Role, role
}
```

**사용 예시:**

```go
roleFunc := func(u *MyUser) string { return u.Role }

g := toulmin.NewGraph[any, *policy.RequestContext[*MyUser], policy.Backing]("admin")
auth  := g.Warrant(policy.IsAuthenticated[*MyUser], policy.NilBacking{}, 1.0)
admin := g.Warrant(policy.IsInRole[*MyUser], &policy.RoleBacking[*MyUser]{
    Role: "admin", RoleFunc: roleFunc,
}, 1.0)
blocked := g.Rebuttal(policy.IsIPInList[*MyUser], &policy.IPListBacking{
    Purpose: "blocklist", Check: isBlocked,
}, 1.0)
g.Defeat(blocked, auth)
```

#### state — `U` (User), `B = state.Backing`

```go
type Backing interface{ stateBacking() }

type TransitionContext[U any] struct {
    User     U
    Resource any
    Metadata map[string]any
}

type OwnerBacking[U any] struct {
    OwnerIDFunc func(any) string
    UserIDFunc  func(U) string
}
func (*OwnerBacking[U]) stateBacking() {}
```

#### approve — `A` (Approver), `B = approve.Backing`

```go
type Backing interface{ approveBacking() }

type ApprovalContext[A any] struct {
    Approver A
    Budget   *Budget
    OrgTree  *OrgTree
    Metadata map[string]any
}

type ApproverBacking[A any] struct {
    Role      string
    Level     int
    IDFunc    func(A) string
    RoleFunc  func(A) string
    LevelFunc func(A) int
}
func (*ApproverBacking[A]) approveBacking() {}
```

#### price — `B = *DiscountBacking` (동질 backing)

price는 모든 rule의 backing이 `*DiscountBacking` 또는 `*MemberBacking`이다. 마커 인터페이스 또는 공통 인터페이스로 통합 가능.

```go
type Backing interface{ priceBacking() }

type DiscountBacking struct {
    Name  string
    Rate  float64
    Fixed float64
    Min   float64
    Max   float64
}
func (*DiscountBacking) priceBacking() {}

type MemberBacking[U any] struct {
    Level          string
    MembershipFunc func(U) string
    Discount       *DiscountBacking
}
func (*MemberBacking[U]) priceBacking() {}

type PriceContext[U any] struct {
    User     U
    Coupons  []Coupon
    Metadata map[string]any
}
```

#### feature — `B = any` (이종 backing)

feature의 backing이 `nil`, `string`, `float64`, `func(any)string`, `[2]any` 등 이종이므로 `B = any`.

```go
// Graph[string, *UserContext, any]
```

#### moderate — `B = moderate.Backing`

```go
type Backing interface{ moderateBacking() }

// Classifier implements Backing
// NilBacking for rules without backing
// int, float64 → IntBacking, FloatBacking wrappers
```

### 각 프레임워크별 타입 파라미터 요약

| 패키지 | 타입 파라미터 | Graph 인스턴스 | B 전략 |
|---|---|---|---|
| policy | `U` (User) | `Graph[any, *RequestContext[U], Backing]` | 마커 인터페이스 |
| state | `U` (User) | `Graph[*TransitionRequest, *TransitionContext[U], Backing]` | 마커 인터페이스 |
| approve | `A` (Approver) | `Graph[*ApprovalRequest, *ApprovalContext[A], Backing]` | 마커 인터페이스 |
| price | `U` (User) | `Graph[*PurchaseRequest, *PriceContext[U], Backing]` | 마커 인터페이스 |
| feature | — | `Graph[string, *UserContext, any]` | any 폴백 |
| moderate | — | `Graph[*Content, *ContentContext, Backing]` | 마커 인터페이스 |

### NilBacking 패턴

backing이 nil인 rule은 각 프레임워크의 `NilBacking` 구조체를 사용한다. `NilBacking{}`은 마커 인터페이스를 만족하므로 B 타입 제약을 통과한다.

```go
// 각 프레임워크별 NilBacking
type NilBacking struct{}
func (NilBacking) policyBacking() {}

// 사용
auth := g.Warrant(policy.IsAuthenticated[*MyUser], policy.NilBacking{}, 1.0)
```

rule 함수 안에서 NilBacking 처리:

```go
func IsAuthenticated[U any](claim any, ground *RequestContext[U], backing Backing) (bool, any) {
    // backing 사용 안 함 — NilBacking이 들어옴
    var zero U
    return any(ground.User) != any(zero), nil
}
```

### 사용 예시 — Before/After

**Before:**

```go
roleFunc := func(u any) string { return u.(*MyUser).Role }
g := toulmin.NewGraph("admin")
auth := g.Warrant(policy.IsAuthenticated, nil, 1.0)
admin := g.Warrant(policy.IsInRole, &policy.RoleBacking{
    Role: "admin", RoleFunc: roleFunc,
}, 1.0)

// rule 함수 안
func IsInRole(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RequestContext)             // 런타임 패닉 위험
    rb := backing.(*RoleBacking)               // 런타임 패닉 위험
    role := rb.RoleFunc(ctx.User)              // func(any) — 런타임 패닉 위험
    return role == rb.Role, role
}
```

**After:**

```go
roleFunc := func(u *MyUser) string { return u.Role }
g := toulmin.NewGraph[any, *policy.RequestContext[*MyUser], policy.Backing]("admin")
auth := g.Warrant(policy.IsAuthenticated[*MyUser], policy.NilBacking{}, 1.0)
admin := g.Warrant(policy.IsInRole[*MyUser], &policy.RoleBacking[*MyUser]{
    Role: "admin", RoleFunc: roleFunc,
}, 1.0)

// rule 함수 안
func IsInRole[U any](claim any, ground *RequestContext[U], backing Backing) (bool, any) {
    rb := backing.(*RoleBacking[U])            // backing만 타입 단언 (B 인터페이스 → concrete)
    role := rb.RoleFunc(ground.User)           // func(U) — 컴파일 타임 안전
    return role == rb.Role, role
}
```

**타입 단언 변화:**

| 위치 | Before | After |
|---|---|---|
| claim | `claim.(type)` | 컴파일 타임 |
| ground | `ground.(*RequestContext)` | 컴파일 타임 |
| ground.User | `func(any) string` 내부 | 컴파일 타임 (`func(U) string`) |
| backing | `backing.(*RoleBacking)` | `backing.(*RoleBacking[U])` — B 인터페이스 → concrete |
| evidence | `result.Evidence.(type)` | `result.Evidence.(type)` — 변경 없음 |

backing 타입 단언이 완전히 제거되지는 않지만 (인터페이스 → concrete), **잘못된 backing 타입이 그래프에 등록되는 것 자체를 컴파일러가 차단**한다. 동질 backing 그래프(price 등)에서는 backing 타입 단언도 완전 제거된다.

### 레거시 호환

```go
// Engine은 기존 string 기반 API. any 시그니처 유지. 변경 없음.
type Engine struct { ... }
func (e *Engine) Register(meta RuleMeta) error
func (e *Engine) Evaluate(claim any, ground any) ([]EvalResult, error)
```

### 트레이드오프

| | Before (`any`) | After (`Graph[C, G, B]`) |
|---|---|---|
| 타입 안전성 | 런타임 패닉 | 컴파일 타임 에러 |
| API 간결성 | `NewGraph("name")` | `NewGraph[any, *Ctx, Backing]("name")` |
| rule 함수 | `func(any, any, any)` | `func(any, *Ctx, Backing)` |
| backing 안전성 | `backing.(type)` 런타임 | 동질: 타입 단언 제로, 이종: 인터페이스→concrete |
| nil backing | `nil` | `NilBacking{}` (보일러플레이트 미미) |
| Go 최소 버전 | Go 1.18 (이미 충족) | Go 1.18 (변경 없음) |
| breaking change | — | **전면 API 파괴** |
| 학습 곡선 | 낮음 | `Graph[C, G, B]` + 마커 인터페이스 이해 필요 |

## 영향 범위

### 코어 엔진 수정 파일 (22개)

| 파일 | 변경 |
|---|---|
| `graph.go` | `Graph` → `Graph[C, G, B]` |
| `rule.go` | `Rule` → `Rule[C, G, B]`, fn `func(C, G, B) (bool, any)` |
| `graph_warrant.go` | 시그니처 변경 |
| `graph_rebuttal.go` | 시그니처 변경 |
| `graph_defeater.go` | 시그니처 변경 |
| `graph_defeat.go` | 시그니처 변경 |
| `graph_evaluate.go` | 시그니처 변경 |
| `graph_evaluate_trace.go` | 시그니처 변경 |
| `new_graph.go` | 시그니처 변경 |
| `eval_context.go` | 내부 타입 파라미터 추가 |
| `eval_context_calc.go` | 내부 타입 파라미터 추가 |
| `eval_context_calc_trace.go` | 내부 타입 파라미터 추가 |
| `new_eval_context.go` | 내부 타입 파라미터 추가 |
| `to_rule_func.go` | 삭제 — 타입 안전 시그니처로 불필요 |
| `wrap_legacy.go` | Engine 전용으로 유지 |
| `detect_cycle.go` | 타입 파라미터 추가 |
| `build_attacker_set.go` | 타입 파라미터 추가 |
| `defeat_edge.go` | 타입 파라미터 추가 |
| `infer_role.go` | 타입 파라미터 추가 |
| `is_warrant.go` | 타입 파라미터 추가 |
| `graph_test.go` | 타입 파라미터 적용 |
| `engine_test.go` | 변경 없음 (Engine은 레거시) |

### 프레임워크 수정 — 패키지당

| 패키지 | 수정 파일 수 (예상) | 주요 변경 |
|---|---|---|
| pkg/policy | 17 | Backing 인터페이스, NilBacking, RequestContext[U], RoleBacking[U], rule 함수, Guard |
| pkg/state | 16 | Backing 인터페이스, NilBacking, TransitionContext[U], OwnerBacking[U], rule 함수, Machine |
| pkg/approve | 22 | Backing 인터페이스, NilBacking, ApprovalContext[A], ApproverBacking[A], rule 함수, Flow |
| pkg/price | 17 | Backing 인터페이스, PriceContext[U], MemberBacking[U], rule 함수, Pricer |
| pkg/feature | 11 | Graph[string, *UserContext, any], Flags, rule 함수 |
| pkg/moderate | 14 | Backing 인터페이스, NilBacking, Moderator, rule 함수 |
| **합계** | **~119** | |

### 신규 파일 (프레임워크별 Backing 인터페이스 + NilBacking)

| 파일 | 내용 |
|---|---|
| `pkg/policy/backing.go` | `Backing` 마커 인터페이스 + `NilBacking` |
| `pkg/state/backing.go` | `Backing` 마커 인터페이스 + `NilBacking` |
| `pkg/approve/backing.go` | `Backing` 마커 인터페이스 + `NilBacking` |
| `pkg/price/backing.go` | `Backing` 마커 인터페이스 |
| `pkg/moderate/backing.go` | `Backing` 마커 인터페이스 + `NilBacking` + 래퍼 타입 |

## 범위

### 포함

1. 코어 엔진 `Graph[C, G, B]`, `Rule[C, G, B]` 제네릭 전환
2. 6개 프레임워크 패키지 제네릭 전환
3. 프레임워크별 `Backing` 마커 인터페이스 + `NilBacking` 도입
4. internal 패키지 Graph 타입 파라미터 반영
5. 전체 테스트 갱신
6. Engine(레거시) API 호환 유지
7. README, manual-for-ai.md 갱신

### 제외

- evidence 타입 제네릭화 (`[]EvalResult`에 이종 evidence 혼재 — any 유지)
- Engine API 변경 (레거시 호환)

## 단계

### Step 1: 코어 엔진 제네릭 전환 (pkg/toulmin)

1. `Graph[C, G, B]`, `Rule[C, G, B]`, 내부 타입에 타입 파라미터 추가
2. `Warrant`, `Rebuttal`, `Defeater` 시그니처 변경: `fn func(C, G, B) (bool, any)`
3. `Evaluate`, `EvaluateTrace` 시그니처 변경: `(claim C, ground G)`
4. `toRuleFunc` 삭제 (타입 안전 시그니처로 불필요)
5. evalContext, detect_cycle 등 내부 타입 전파
6. `graph_test.go` 갱신
7. Engine(레거시)은 변경 없음

### Step 2: feature 프레임워크 (B = any 폴백)

8. `Graph[string, *UserContext, any]` 적용
9. Flags에 타입 파라미터 전파
10. rule 함수 시그니처: `func(string, *UserContext, any) (bool, any)`
11. 테스트 갱신

### Step 3: moderate 프레임워크 (Backing 마커 인터페이스)

12. `Backing` 마커 인터페이스 + `NilBacking` + 래퍼 타입 (IntBacking, FloatBacking)
13. `Graph[*Content, *ContentContext, moderate.Backing]` 적용
14. Moderator에 타입 파라미터 전파
15. rule 함수 시그니처 변경
16. 테스트 갱신

### Step 4: policy, state 프레임워크 (User 제네릭 + Backing 인터페이스)

17. `Backing` 마커 인터페이스 + `NilBacking` 도입
18. RequestContext[U], TransitionContext[U] 작성
19. RoleBacking[U], OwnerBacking[U] 변경 + `policyBacking()` 구현
20. rule 함수를 제네릭 함수로 변경
21. Guard, Machine에 타입 파라미터 전파
22. 테스트 갱신

### Step 5: approve, price 프레임워크 (Approver/User 제네릭 + Backing 인터페이스)

23. `Backing` 마커 인터페이스 도입
24. ApprovalContext[A], ApproverBacking[A] 변경
25. PriceContext[U], MemberBacking[U] 변경
26. Flow[A], Pricer[U] 타입 파라미터 전파
27. 테스트 갱신

### Step 6: internal 패키지 갱신

28. internal/analyzer, internal/codegen, internal/graph — Graph 타입 파라미터 반영
29. codegen 템플릿이 제네릭 코드를 생성하도록 갱신

### Step 7: 검증

30. `go build ./...`
31. `go vet ./...`
32. `go test ./...` 전체 PASS
33. README, manual-for-ai.md 갱신

## 검증 기준

1. `Graph[C, G, B]`로 타입화된 claim/ground/backing을 사용한다
2. rule 함수 시그니처에서 claim/ground의 `any` 타입 단언이 제거되었다
3. 동질 backing 그래프에서 backing 타입 단언이 완전히 제거되었다
4. 마커 인터페이스 backing 그래프에서 잘못된 backing 타입이 컴파일 에러를 발생시킨다
5. evidence는 `any`를 유지한다
6. Engine(레거시) API는 변경 없이 동작한다
7. 잘못된 claim/ground 타입을 넘기면 컴파일 에러가 발생한다
8. 프레임워크 backing의 추출 함수가 `func(U) string` (타입 안전)이다
9. NilBacking 패턴이 nil backing을 대체한다
10. `go build ./...` 통과
11. `go vet ./...` 통과
12. `go test ./...` 전체 PASS
13. README, manual-for-ai.md 갱신 완료

## 의존성

- Phase 001-017: 전체 기존 구현
- Phase 018: Gin 분리 완료 (Guard 시그니처가 net/http 기반이어야 제네릭 전환 깔끔)
- Go 1.18+ (제네릭 지원, 현재 go.mod는 1.25.0이므로 충족)

## 예상 규모

- 수정 파일: ~119개
- 삭제 파일: 1개 (`to_rule_func.go`)
- 신규 파일: 5개 (프레임워크별 `backing.go`)
- 예상 난이도: **높음** (전면 API 파괴, 모든 테스트 갱신, 타입 파라미터 전파 복잡)
- **breaking change** — 메이저 버전 범프 필요 (v2)
- 핵심 난점: 코어 `Graph[C, G, B]` 변경이 6개 프레임워크 + internal 전체에 전파되므로, Step 1 완료 후 빌드가 깨진 상태에서 Step 2-6을 연속 진행해야 함
