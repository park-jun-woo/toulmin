# Phase 013: 상태 전이 프레임워크 — pkg/state

## 목표

toulmin 기반 상태 전이 프레임워크를 `pkg/state`에 구현한다.
상태 전이의 허용/거부를 defeats graph로 판정한다.
Mermaid stateDiagram + 런타임 가드를 하나의 모델로 대체한다.

## 배경

### 현재 문제

1. **상태 전이 로직이 if-else로 흩어진다**: `if current == "pending" && isOwner && !isExpired` 같은 분기가 핸들러마다 반복된다
2. **다이어그램과 코드가 분리된다**: Mermaid stateDiagram은 문서이고, 실제 가드 로직은 코드에 따로 있다. 불일치가 발생한다
3. **전이 가드에 예외가 생기면 복잡해진다**: 관리자 오버라이드, 만료 예외 등이 if-else 중첩을 만든다

### toulmin이 해결하는 것

- 전이 하나 = graph 하나 (전이 허용 조건의 선언적 정의)
- 전이 가드 = rule 함수 (1-2 depth)
- 예외 = defeat edge
- 전이 판정 근거 = EvaluateTrace

### claim/ground/backing 분리 원칙

toulmin의 `(claim any, ground any, backing any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 상태 전이 프레임워크에서 claim은 전이 요청 (from, to, event)
- **ground = 판정 재료**: ground는 전이 판정에 필요한 컨텍스트 (현재 상태, 사용자, 리소스)
- **backing = 규칙의 판정 기준**: backing은 rule 함수의 판정 기준 (ID 추출 함수, 만료 시간 추출 함수 등). 선언 시점에 고정되며, 클로저에 숨기지 않고 명시적으로 전달한다.

프레임워크는 Machine 구조와 판정 흐름을 제공하고, **도메인 데이터는 ground로 사용자가 주입한다.** rule 함수는 claim/ground/backing에서 데이터를 꺼내 판단만 하므로, 프레임워크가 도메인을 몰라도 동작한다.

| 역할 | 상태 전이 프레임워크에서 |
|---|---|
| claim | TransitionRequest (from, to, event) |
| ground | TransitionContext (CurrentState, User, Resource, Metadata) |
| backing | 규칙의 판정 기준 (ID 추출 함수, 만료 시간 추출 함수 등) |
| rule 함수 | claim/ground/backing에서 조건 하나만 판단 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 예외) |
| verdict | 전이 허용/거부 판정 |

정책 프레임워크(Phase 010)와의 차이: 정책에서는 claim이 nil(라우트 매칭으로 확정)이지만, 상태 전이에서는 **claim이 TransitionRequest로 활성화된다.** IsCurrentState가 `claim.From == ground.CurrentState`를 비교하는 것이 이를 보여준다. 같은 graph, 같은 rule 함수 세트가 claim/ground의 의미만 바꿔서 다른 도메인에 적용된다.

## 핵심 설계

### TransitionContext

```go
// pkg/state/transition_context.go
type TransitionContext struct {
    CurrentState    string
    User            any
    Resource        any
    Metadata        map[string]any
}
```

### TransitionRequest

```go
// pkg/state/transition_request.go
type TransitionRequest struct {
    From   string
    To     string
    Event  string
}
```

### 범용 rule 함수

클로저 팩토리를 사용하지 않는다. backing으로 판정 기준을 명시적으로 전달한다.

```go
// pkg/state/rule_is_current_state.go
// backing 불필요 — claim과 ground만으로 판정
func IsCurrentState(claim any, ground any, backing any) (bool, any) {
    req := claim.(*TransitionRequest)
    ctx := ground.(*TransitionContext)
    return ctx.CurrentState == req.From, nil
}

// pkg/state/rule_is_owner.go
// backing: OwnerBacking{OwnerIDFunc, UserIDFunc} — ID 추출 함수 쌍
type OwnerBacking struct {
    OwnerIDFunc func(any) string
    UserIDFunc  func(any) string
}

func IsOwner(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*TransitionContext)
    b := backing.(*OwnerBacking)
    return b.UserIDFunc(ctx.User) == b.OwnerIDFunc(ctx.Resource), nil
}

// pkg/state/rule_is_expired.go
// backing: func(any) time.Time — 만료 시간 추출 함수
func IsExpired(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*TransitionContext)
    expiryFunc := backing.(func(any) time.Time)
    return time.Now().After(expiryFunc(ctx.Resource)), nil
}
```

### Machine — 전이 그래프 등록 및 실행

```go
// pkg/state/machine.go
type Machine struct {
    transitions map[string]*toulmin.Graph  // key: "from:event"
}

func NewMachine() *Machine

// Add — 전이 하나에 대한 판정 graph 등록
func (m *Machine) Add(from, event, to string, g *toulmin.Graph)

// Can — 전이 가능 여부 판정 (verdict 반환)
func (m *Machine) Can(req *TransitionRequest, ctx *TransitionContext) (float64, error)

// CanTrace — 전이 가능 여부 + 판정 근거
func (m *Machine) CanTrace(req *TransitionRequest, ctx *TransitionContext) (*TraceResult, error)
```

### TraceResult

```go
// pkg/state/trace_result.go
type TraceResult struct {
    Verdict float64
    From    string
    To      string
    Event   string
    Trace   []toulmin.TraceEntry
}
```

### 사용 예시

**주의**: 클로저를 사용하지 않는다. backing으로 판정 기준을 명시적으로 전달한다. Warrant/Rebuttal/Defeater는 `*Rule`을 반환하며 체이닝하지 않는다. Rebuttal만으로는 공격이 일어나지 않으며 반드시 Defeat edge를 선언해야 한다. 예외를 처리하는 rule은 Defeater로 등록해야 한다. verdict <= 0이면 전이를 거부한다.

```go
m := state.NewMachine()

// backing 준비
ownerBacking := &state.OwnerBacking{
    OwnerIDFunc: func(r any) string { return r.(*Proposal).OwnerID },
    UserIDFunc:  func(u any) string { return u.(*User).ID },
}
expiryFunc := func(r any) time.Time { return r.(*Proposal).ExpiresAt }

// pending → accepted 전이
g := toulmin.NewGraph("proposal:accept")
current := g.Warrant(state.IsCurrentState, nil, 1.0)
owner := g.Warrant(state.IsOwner, ownerBacking, 1.0)
auth := g.Warrant(isAuthenticated, nil, 1.0)
expired := g.Rebuttal(state.IsExpired, expiryFunc, 1.0)
override := g.Defeater(isAdminOverride, nil, 1.0)
g.Defeat(expired, current)
g.Defeat(override, expired)
m.Add("pending", "accept", "accepted", g)

// pending → rejected 전이
g2 := toulmin.NewGraph("proposal:reject")
g2.Warrant(state.IsCurrentState, nil, 1.0)
g2.Warrant(state.IsOwner, ownerBacking, 1.0)
g2.Warrant(isAuthenticated, nil, 1.0)
m.Add("pending", "reject", "rejected", g2)

// 전이 판정
req := &state.TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
ctx := &state.TransitionContext{
    CurrentState: proposal.Status,
    User:         currentUser,
    Resource:     proposal,
}

verdict, err := m.Can(req, ctx)
// verdict > 0: 전이 허용
// verdict <= 0: 전이 거부 (undecided는 거부)
```

### 같은 함수 + 다른 backing이 필요한 경우

Warrant/Rebuttal/Defeater가 `*Rule`을 반환하므로, 같은 함수를 다른 backing으로 등록하면 다른 `*Rule` 참조를 얻는다. Defeat는 `*Rule` 참조로 관계를 선언한다:

```go
// 두 가지 소유권 검증이 필요한 경우
creatorBacking := &state.OwnerBacking{
    OwnerIDFunc: func(r any) string { return r.(*Doc).CreatorID },
    UserIDFunc:  func(u any) string { return u.(*User).ID },
}
approverBacking := &state.OwnerBacking{
    OwnerIDFunc: func(r any) string { return r.(*Doc).ApproverID },
    UserIDFunc:  func(u any) string { return u.(*User).ID },
}

g := toulmin.NewGraph("doc:approve")
approver := g.Warrant(state.IsOwner, approverBacking, 1.0)
creator := g.Rebuttal(state.IsOwner, creatorBacking, 1.0)
g.Defeat(creator, approver)
```

### 상태 다이어그램 추출

```go
// pkg/state/diagram.go
// Mermaid — Machine에 등록된 전이 목록에서 Mermaid stateDiagram 생성
func (m *Machine) Mermaid() string
```

```
stateDiagram-v2
    pending --> accepted : accept
    pending --> rejected : reject
    accepted --> completed : complete
```

등록된 graph에서 from/to/event를 추출하므로, 다이어그램과 런타임 가드가 항상 일치한다.

## 범위

### 포함

1. **TransitionContext, TransitionRequest 구조체**: 전이 판정에 필요한 컨텍스트
2. **OwnerBacking 구조체**: IsOwner rule의 backing 타입
3. **범용 rule 함수**: IsCurrentState(backing nil), IsOwner(backing OwnerBacking), IsExpired(backing expiryFunc) — 모두 순수 함수
4. **Machine**: 전이 graph 등록, Can/CanTrace 판정
5. **TraceResult**: 판정 근거 구조체
6. **Mermaid 다이어그램 생성**: Machine → stateDiagram 추출
7. **테스트**: rule 함수 단위 테스트, Machine 통합 테스트

### 제외

- 상태 퍼시스턴스 (DB 저장은 사용자 책임)
- 이벤트 버스/pub-sub 연동
- 상태 전이 YAML 정의 및 코드젠 — fullend에서 처리
- Gin 어댑터 — pkg/policy의 Guard와 조합하여 사용

## 산출물

```
pkg/
  state/
    transition_context.go         — TransitionContext 구조체
    transition_request.go         — TransitionRequest 구조체
    trace_result.go               — TraceResult 구조체
    owner_backing.go              — OwnerBacking 구조체 (backing 타입)
    rule_is_current_state.go      — IsCurrentState (backing nil, 순수 함수)
    rule_is_owner.go              — IsOwner (backing OwnerBacking, 순수 함수)
    rule_is_expired.go            — IsExpired (backing expiryFunc, 순수 함수)
    machine.go                    — Machine (NewMachine, Add, Can, CanTrace)
    diagram.go                    — Mermaid 다이어그램 생성
    rule_test.go                  — rule 함수 단위 테스트
    machine_test.go               — Machine 통합 테스트
    diagram_test.go               — Mermaid 출력 테스트
```

## 단계

### Step 1: 구조체 정의

- `pkg/state/transition_context.go`: TransitionContext
- `pkg/state/transition_request.go`: TransitionRequest
- `pkg/state/trace_result.go`: TraceResult
- `pkg/state/owner_backing.go`: OwnerBacking

### Step 2: rule 함수 구현

- IsCurrentState: `func(claim any, ground any, backing any) (bool, any)` — backing nil, claim.From == ground.CurrentState
- IsOwner: `func(claim any, ground any, backing any) (bool, any)` — backing으로 OwnerBacking을 받아 ID 추출 후 비교
- IsExpired: `func(claim any, ground any, backing any) (bool, any)` — backing으로 만료 시간 추출 함수를 받아 만료 판정

### Step 3: Machine 구현

- NewMachine: 빈 Machine 생성
- Add: from + event 키로 graph 등록, to 저장
- Can: 키로 graph 조회 → Evaluate → verdict 반환
- CanTrace: 키로 graph 조회 → EvaluateTrace → TraceResult 반환

### Step 4: Mermaid 다이어그램 생성

- Machine에 등록된 전이 목록을 순회하여 `stateDiagram-v2` 포맷 출력
- 등록 순서대로 출력

### Step 5: 테스트

- rule 함수 단위 테스트: 상태 일치/불일치, 소유자 판정 (backing으로 ID 추출), 만료 판정 (backing으로 만료 시간 추출)
- Machine 통합 테스트:
  - 유효한 전이 → verdict > 0
  - 상태 불일치 → verdict <= 0
  - defeat edge 동작 (만료됐지만 관리자 → 허용)
  - 미등록 전이 → 에러
- Mermaid 출력 테스트: 등록된 전이가 다이어그램에 포함

### Step 6: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. IsCurrentState가 claim.From과 ground.CurrentState를 비교한다 (backing nil)
2. IsOwner가 backing(OwnerBacking)에서 ID 추출 함수를 받아 소유자를 판정한다
3. IsExpired가 backing(expiryFunc)에서 만료 시간 추출 함수를 받아 만료를 판정한다
4. Machine.Can이 등록된 graph로 전이 가능 여부를 판정한다
5. Machine.CanTrace가 판정 근거를 TraceResult로 반환한다
6. defeat edge가 예외를 정확히 처리한다 (만료 + 관리자 오버라이드)
7. 미등록 전이에 대해 에러를 반환한다
8. Mermaid()가 등록된 전이와 일치하는 다이어그램을 생성한다
9. 다이어그램과 런타임 가드가 동일한 소스(Machine)에서 나온다
10. verdict <= 0이면 전이를 거부한다 (undecided 포함)
11. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
- Phase 010: backing 일급 시민 — rule 시그니처 `(claim, ground, backing)`, API `(fn, backing, qualifier)`
- Phase 012: Rule 참조 반환 + 체이닝 제거 — Warrant/Rebuttal/Defeater가 `*Rule` 반환, Defeat는 `*Rule` 참조, GraphBuilder → Graph, DefeatWith 제거
