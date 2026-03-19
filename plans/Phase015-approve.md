# Phase 015: 승인 워크플로우 프레임워크 — pkg/approve

## 목표

toulmin 기반 승인 워크플로우 프레임워크를 `pkg/approve`에 구현한다.
다단계 결재의 각 단계를 defeats graph로 판정한다.
긴급 승인, 상위자 오버라이드, 금액 기준 분기 같은 예외를 defeat edge로 선언적 처리한다.

## 배경

### 현재 문제

1. **승인 로직이 if-else 지옥이다**: `if amount > 1000 && approver.Role == "manager" && !isFrozen && (isUrgent || approver.Level >= 3)` 같은 분기가 서비스 레이어에 박힌다
2. **단계별 조건이 코드에 흩어진다**: 1차 승인(팀장), 2차 승인(부서장), 최종 승인(CFO) 각각의 조건이 다른 함수에 산재한다. 전체 흐름 파악이 어렵다
3. **예외 처리가 하드코딩이다**: 긴급 결재, CEO 직권 승인, 소액 자동 승인 같은 예외가 if문으로 박혀서 규칙 추가/변경 시 사이드이펙트가 발생한다

### toulmin이 해결하는 것

- 승인 단계 하나 = graph 하나 (승인 조건의 선언적 정의)
- 승인 조건 = rule 함수 (1-2 depth)
- 예외 = defeat edge (긴급 승인, 직권 승인)
- 승인 판정 근거 = EvaluateTrace (감사 추적)

### claim/ground/backing 3-element 분리 원칙

toulmin의 `func(claim any, ground any, backing any) (bool, any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 승인 프레임워크에서 claim은 승인 요청 (금액, 항목, 요청자)
- **ground = 판정 재료**: ground는 승인 판정에 필요한 컨텍스트 (결재자, 예산, 조직 구조)
- **backing = 규칙의 판정 기준**: backing은 rule 함수가 판정할 때 사용하는 기준값 (role명, 레벨 임계값, 금액 임계값). graph 구성 시 고정되며, 클로저에 숨기지 않고 엔진이 관리하는 명시적 값이다

프레임워크는 Flow 구조와 판정 흐름을 제공하고, **도메인 데이터는 ground로, 판정 기준은 backing으로 사용자가 주입한다.** rule 함수는 claim/ground/backing에서 데이터를 꺼내 판단만 하므로, 프레임워크가 도메인을 몰라도 동작한다.

| 역할 | 승인 프레임워크에서 |
|---|---|
| claim | ApprovalRequest (Amount, Category, RequesterID) |
| ground | ApprovalContext (Approver, Budget, OrgTree, Metadata) |
| backing | 규칙별 판정 기준 (role명, 레벨 임계값, 금액 임계값). nil이면 backing 없음 |
| rule 함수 | claim/ground/backing에서 조건 하나만 판단 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 예외) |
| verdict | > 0 승인, <= 0 거부 (undecided = 거부) |

정책(Phase 010)에서 backing이 일급 시민이 되어 클로저가 사라졌다. 승인에서는 **claim이 ApprovalRequest로 활성화되고, backing이 role명/임계값 등 판정 기준을 명시적으로 전달한다.** HasApprovalRole이 `backing.(string)`으로 role명을 받아 `ground.Approver.Role == role`을 비교하는 것이 이를 보여준다.

## 핵심 설계

### ApprovalRequest

```go
// pkg/approve/approval_request.go
type ApprovalRequest struct {
    Amount      float64
    Category    string
    RequesterID string
    Description string
    Metadata    map[string]any
}
```

### ApprovalContext

```go
// pkg/approve/approval_context.go
type ApprovalContext struct {
    Approver    *Approver
    Budget      *Budget
    OrgTree     OrgTree
    Metadata    map[string]any
}

// pkg/approve/approver.go
type Approver struct {
    ID    string
    Role  string
    Level int
}

// pkg/approve/budget.go
type Budget struct {
    Remaining float64
    Frozen    bool
}

// pkg/approve/org_tree.go
type OrgTree interface {
    IsDirectManager(approverID, requesterID string) bool
    Level(userID string) int
}
```

### 범용 rule 함수

모든 rule 함수는 `func(claim any, ground any, backing any) (bool, any)` 시그니처를 따른다. 클로저 팩토리를 사용하지 않으며, 판정 기준은 backing으로 명시적 전달한다.

```go
// pkg/approve/rule_is_direct_manager.go
// backing: nil (backing 불필요)
func IsDirectManager(claim any, ground any, backing any) (bool, any) {
    req := claim.(*ApprovalRequest)
    ctx := ground.(*ApprovalContext)
    return ctx.OrgTree.IsDirectManager(ctx.Approver.ID, req.RequesterID), nil
}

// pkg/approve/rule_is_under_budget.go
// backing: nil (backing 불필요)
func IsUnderBudget(claim any, ground any, backing any) (bool, any) {
    req := claim.(*ApprovalRequest)
    ctx := ground.(*ApprovalContext)
    return req.Amount <= ctx.Budget.Remaining, nil
}

// pkg/approve/rule_is_budget_frozen.go
// backing: nil (backing 불필요)
func IsBudgetFrozen(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ApprovalContext)
    return ctx.Budget.Frozen, nil
}

// pkg/approve/rule_has_approval_role.go
// backing: string (role명, 예: "finance", "cfo")
func HasApprovalRole(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ApprovalContext)
    role := backing.(string)
    return ctx.Approver.Role == role, nil
}

// pkg/approve/rule_is_above_level.go
// backing: int (최소 레벨 임계값, 예: 3)
func IsAboveLevel(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ApprovalContext)
    minLevel := backing.(int)
    return ctx.Approver.Level >= minLevel, nil
}

// pkg/approve/rule_is_small_amount.go
// backing: float64 (금액 임계값, 예: 10000)
func IsSmallAmount(claim any, ground any, backing any) (bool, any) {
    req := claim.(*ApprovalRequest)
    threshold := backing.(float64)
    return req.Amount <= threshold, nil
}

// pkg/approve/rule_is_urgent.go
// backing: nil (backing 불필요)
func IsUrgent(claim any, ground any, backing any) (bool, any) {
    req := claim.(*ApprovalRequest)
    urgent, _ := req.Metadata["urgent"].(bool)
    return urgent, nil
}

// pkg/approve/rule_is_ceo_override.go
// backing: nil (backing 불필요)
func IsCEOOverride(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ApprovalContext)
    return ctx.Approver.Role == "ceo", nil
}
```

### Flow -- 다단계 승인 흐름

```go
// pkg/approve/flow.go
type Flow struct {
    name  string
    steps []*Step  // 순서대로 실행
}

type Step struct {
    Name  string
    Graph *toulmin.Graph
}

func NewFlow(name string) *Flow

// Step -- 승인 단계 추가
func (f *Flow) Step(name string, g *toulmin.Graph) *Flow

// Evaluate -- 전체 흐름 판정 (모든 단계 통과해야 승인)
func (f *Flow) Evaluate(req *ApprovalRequest, ctxBuilder StepContextFunc) (*FlowResult, error)

// StepContextFunc -- 단계별로 다른 ground 주입
type StepContextFunc func(stepName string) *ApprovalContext
```

### FlowResult

```go
// pkg/approve/flow_result.go
type FlowResult struct {
    Approved bool
    Steps    []StepResult
}

type StepResult struct {
    Name    string
    Verdict float64
    Trace   []toulmin.TraceEntry
}
```

### 사용 예시

클로저 팩토리를 사용하지 않는다. 판정 기준은 `Warrant(fn, backing, qualifier)` 형태로 backing 인자에 직접 전달한다. backing이 불필요한 rule은 `nil`을 명시적으로 전달한다. Warrant/Rebuttal/Defeater는 `*Rule` 참조를 반환하며 체이닝하지 않는다. Defeat는 `*Rule` 참조 두 개를 받아 관계를 선언한다. Rebuttal만으로는 공격이 일어나지 않으며 반드시 Defeat edge를 선언해야 한다. 예외를 처리하는 rule은 Defeater로 등록해야 한다.

```go
// 경비 승인 워크플로우
// 클로저 없음 — backing으로 판정 기준 명시
// 체이닝 없음 — 정의와 관계 분리 (Phase 012 참조 패턴)

// Step 1: manager 단계 graph 구성
mgrGraph := toulmin.NewGraph("expense:manager")
mgrDirect  := mgrGraph.Warrant(approve.IsDirectManager, nil, 1.0)
mgrBudget  := mgrGraph.Warrant(approve.IsUnderBudget, nil, 1.0)
mgrFrozen  := mgrGraph.Rebuttal(approve.IsBudgetFrozen, nil, 1.0)
mgrUrgent  := mgrGraph.Defeater(approve.IsUrgent, nil, 1.0)              // 예외 rule은 Defeater로 등록
mgrGraph.Defeat(mgrFrozen, mgrBudget)                                     // Rebuttal -> Warrant 공격 edge 필수
mgrGraph.Defeat(mgrUrgent, mgrFrozen)                                     // Defeater -> Rebuttal 예외 처리

// Step 2: finance 단계 graph 구성
finGraph := toulmin.NewGraph("expense:finance")
finRole    := finGraph.Warrant(approve.HasApprovalRole, "finance", 1.0)   // backing: "finance" (role명)
finBudget  := finGraph.Warrant(approve.IsUnderBudget, nil, 1.0)
finFrozen  := finGraph.Rebuttal(approve.IsBudgetFrozen, nil, 1.0)
finCEO     := finGraph.Defeater(approve.IsCEOOverride, nil, 1.0)
finGraph.Defeat(finFrozen, finBudget)
finGraph.Defeat(finCEO, finFrozen)

// Flow 구성
f := approve.NewFlow("expense").
    Step("manager", mgrGraph).
    Step("finance", finGraph)

// 소액 자동 승인 -- backing으로 금액 임계값 전달
autoGraph := toulmin.NewGraph("expense:auto")
_ = autoGraph.Warrant(approve.IsSmallAmount, float64(10000), 1.0)         // backing: 10000 (금액 임계값)

req := &approve.ApprovalRequest{
    Amount:      50000,
    Category:    "travel",
    RequesterID: "emp-123",
}

result, err := f.Evaluate(req, func(stepName string) *approve.ApprovalContext {
    switch stepName {
    case "manager":
        return &approve.ApprovalContext{
            Approver: manager,
            Budget:   deptBudget,
            OrgTree:  orgTree,
        }
    case "finance":
        return &approve.ApprovalContext{
            Approver: financeHead,
            Budget:   companyBudget,
            OrgTree:  orgTree,
        }
    }
    return nil
})
// result.Approved: 모든 단계 verdict > 0이면 true (undecided = 거부)
// result.Steps[0].Trace: 각 단계별 판정 근거 (backing 값 포함)
```

### backing이 있는 rule 간 Defeat 예시

같은 함수를 다른 backing으로 등록한 경우, 각각의 `*Rule` 참조로 Defeat edge를 선언한다. DefeatWith는 불필요하다 (Phase 012에서 제거됨):

```go
// 레벨 3 이상이면 승인, 레벨 5 이상이면 레벨 3 규칙을 무시
g := toulmin.NewGraph("level-override")
level3 := g.Warrant(approve.IsAboveLevel, 3, 1.0)    // backing: 3 (최소 레벨)
level5 := g.Rebuttal(approve.IsAboveLevel, 5, 1.0)   // backing: 5 (상위 레벨)
g.Defeat(level5, level3)                               // *Rule 참조로 직접 지정
```

## 범위

### 포함

1. **ApprovalRequest, ApprovalContext 구조체**: 승인 판정에 필요한 요청/컨텍스트
2. **Approver, Budget 구조체, OrgTree 인터페이스**: 조직 구조 추상화
3. **범용 rule 함수**: IsDirectManager, IsUnderBudget, IsBudgetFrozen, HasApprovalRole, IsAboveLevel, IsSmallAmount, IsUrgent, IsCEOOverride -- 모두 `func(claim any, ground any, backing any) (bool, any)` 순수 함수
4. **Flow**: 다단계 승인 흐름 (Step 등록, 순차 판정)
5. **FlowResult, StepResult**: 단계별 판정 결과 + trace (backing 포함)
6. **테스트**: rule 함수 단위 테스트, Flow 통합 테스트

### 제외

- 승인 상태 퍼시스턴스 (DB 저장은 사용자 책임)
- 알림/이메일 발송
- UI/대시보드
- 병렬 승인 (동시 결재) -- 별도 Phase에서 검토

## 산출물

```
pkg/
  approve/
    approval_request.go           -- ApprovalRequest 구조체
    approval_context.go           -- ApprovalContext 구조체
    approver.go                   -- Approver 구조체
    budget.go                     -- Budget 구조체
    org_tree.go                   -- OrgTree 인터페이스
    rule_is_direct_manager.go     -- IsDirectManager (backing: nil)
    rule_is_under_budget.go       -- IsUnderBudget (backing: nil)
    rule_is_budget_frozen.go      -- IsBudgetFrozen (backing: nil)
    rule_has_approval_role.go     -- HasApprovalRole (backing: string, role명)
    rule_is_above_level.go        -- IsAboveLevel (backing: int, 레벨 임계값)
    rule_is_small_amount.go       -- IsSmallAmount (backing: float64, 금액 임계값)
    rule_is_urgent.go             -- IsUrgent (backing: nil)
    rule_is_ceo_override.go       -- IsCEOOverride (backing: nil)
    flow.go                       -- Flow (NewFlow, Step, Evaluate)
    flow_result.go                -- FlowResult, StepResult
    rule_test.go                  -- rule 함수 단위 테스트
    flow_test.go                  -- Flow 통합 테스트
```

## 단계

### Step 1: 구조체 및 인터페이스 정의

- ApprovalRequest, ApprovalContext, Approver, Budget, OrgTree
- FlowResult, StepResult

### Step 2: rule 함수 구현

- 각 rule 함수를 파일 하나에 하나씩 구현 (filefunc 규칙 준수)
- 각 함수는 1-2 depth 유지
- 모든 rule 함수는 `func(claim any, ground any, backing any) (bool, any)` 시그니처
- 클로저 팩토리 사용하지 않음 -- HasApprovalRole, IsAboveLevel, IsSmallAmount는 backing 인자로 판정 기준 수신

### Step 3: Flow 구현

- NewFlow: 빈 Flow 생성
- Step: graph를 단계로 등록
- Evaluate: 단계를 순차 실행, 각 단계마다 StepContextFunc으로 ground 주입, verdict <= 0이면 거부 (undecided = 거부)

### Step 4: 테스트

- rule 함수 단위 테스트: 각 조건별 true/false, backing 값 전달 검증
- Flow 통합 테스트:
  - 모든 단계 통과 -> Approved: true
  - 중간 단계 거부 -> Approved: false, 거부 단계 식별
  - defeat edge 동작 (예산 동결 + 긴급 -> 승인)
  - CEO 오버라이드 동작
  - 소액 자동 승인 graph 동작 (backing으로 금액 임계값 전달)
  - backing이 있는 rule 간 *Rule 참조 Defeat 동작

### Step 5: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. 모든 rule 함수가 `func(claim any, ground any, backing any) (bool, any)` 시그니처를 따른다
2. HasApprovalRole, IsAboveLevel, IsSmallAmount가 backing 인자에서 판정 기준을 받아 동작한다
3. 클로저 팩토리가 존재하지 않는다
4. `Warrant(fn, backing, qualifier)` 형태로 backing을 명시적으로 전달한다
5. backing이 nil인 rule은 `Warrant(fn, nil, qualifier)` 로 명시적 선언한다
6. 같은 함수 + 다른 backing을 가진 rule이 `*Rule` 참조로 구분되어 `Defeat`로 동작한다
7. Flow.Evaluate가 모든 단계를 순차 실행하고 결과를 집계한다
8. 중간 단계에서 거부되면 이후 단계를 실행하지 않는다
9. StepContextFunc으로 단계별 다른 ground를 주입할 수 있다
10. defeat edge가 예외를 정확히 처리한다 (긴급 -> 예산 동결 무시)
11. FlowResult.Steps에 각 단계의 trace가 포함된다 (backing 값 포함)
12. verdict > 0 승인, verdict <= 0 거부 (undecided = 거부)
13. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
- Phase 010: backing 일급 시민 (`func(claim, ground, backing)` 시그니처, `Warrant(fn, backing, qualifier)` API)
- Phase 012: Rule 참조 반환 + 체이닝 제거 (Warrant/Rebuttal/Defeater → `*Rule`, `Defeat(*Rule, *Rule)`, GraphBuilder → Graph, DefeatWith 제거)
