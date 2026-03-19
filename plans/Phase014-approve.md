# Phase 012: 승인 워크플로우 프레임워크 — pkg/approve

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

### claim/ground 분리 원칙

toulmin의 `(claim any, ground any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 승인 프레임워크에서 claim은 승인 요청 (금액, 항목, 요청자)
- **ground = 판정 재료**: ground는 승인 판정에 필요한 컨텍스트 (결재자, 예산, 조직 구조)

프레임워크는 Flow 구조와 판정 흐름을 제공하고, **도메인 데이터는 ground로 사용자가 주입한다.** rule 함수는 claim/ground에서 데이터를 꺼내 판단만 하므로, 프레임워크가 도메인을 몰라도 동작한다.

| 역할 | 승인 프레임워크에서 |
|---|---|
| claim | ApprovalRequest (Amount, Category, RequesterID) |
| ground | ApprovalContext (Approver, Budget, OrgTree, Metadata) |
| rule 함수 | claim/ground에서 조건 하나만 판단 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 예외) |
| verdict | 승인/거부 판정 |

정책(Phase 010)에서는 claim이 nil, 상태 전이(Phase 011)에서는 claim이 TransitionRequest였다. 승인에서는 **claim이 ApprovalRequest로 활성화된다.** IsUnderBudget이 `claim.Amount <= ground.Budget.Remaining`을 비교하는 것이 이를 보여준다.

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

```go
// pkg/approve/rule_is_direct_manager.go
func IsDirectManager(claim any, ground any) (bool, any) {
    req := claim.(*ApprovalRequest)
    ctx := ground.(*ApprovalContext)
    return ctx.OrgTree.IsDirectManager(ctx.Approver.ID, req.RequesterID), nil
}

// pkg/approve/rule_is_under_budget.go
func IsUnderBudget(claim any, ground any) (bool, any) {
    req := claim.(*ApprovalRequest)
    ctx := ground.(*ApprovalContext)
    return req.Amount <= ctx.Budget.Remaining, nil
}

// pkg/approve/rule_is_budget_frozen.go
func IsBudgetFrozen(claim any, ground any) (bool, any) {
    ctx := ground.(*ApprovalContext)
    return ctx.Budget.Frozen, nil
}

// pkg/approve/rule_has_approval_role.go
func HasApprovalRole(role string) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*ApprovalContext)
        return ctx.Approver.Role == role, nil
    }
}

// pkg/approve/rule_is_above_level.go
func IsAboveLevel(minLevel int) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*ApprovalContext)
        return ctx.Approver.Level >= minLevel, nil
    }
}

// pkg/approve/rule_is_small_amount.go
func IsSmallAmount(threshold float64) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        req := claim.(*ApprovalRequest)
        return req.Amount <= threshold, nil
    }
}

// pkg/approve/rule_is_urgent.go
func IsUrgent(claim any, ground any) (bool, any) {
    req := claim.(*ApprovalRequest)
    urgent, _ := req.Metadata["urgent"].(bool)
    return urgent, nil
}

// pkg/approve/rule_is_ceo_override.go
func IsCEOOverride(claim any, ground any) (bool, any) {
    ctx := ground.(*ApprovalContext)
    return ctx.Approver.Role == "ceo", nil
}
```

### Flow — 다단계 승인 흐름

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

// Step — 승인 단계 추가
func (f *Flow) Step(name string, g *toulmin.Graph) *Flow

// Evaluate — 전체 흐름 판정 (모든 단계 통과해야 승인)
func (f *Flow) Evaluate(req *ApprovalRequest, ctxBuilder StepContextFunc) (*FlowResult, error)

// StepContextFunc — 단계별로 다른 ground 주입
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

**주의**: 클로저 rule은 변수에 저장 후 재사용해야 한다. Rebuttal만으로는 공격이 일어나지 않으며 반드시 Defeat edge를 선언해야 한다. 예외를 처리하는 rule은 Defeater로 등록해야 한다.

```go
// 클로저는 변수에 저장 후 재사용
hasFinance := approve.HasApprovalRole("finance")
smallAmount := approve.IsSmallAmount(10000)

// 경비 승인 워크플로우
f := approve.NewFlow("expense").
    Step("manager",
        toulmin.NewGraph("expense:manager").
            Warrant(approve.IsDirectManager, 1.0).
            Warrant(approve.IsUnderBudget, 1.0).
            Rebuttal(approve.IsBudgetFrozen, 1.0).
            Defeater(approve.IsUrgent, 1.0).                      // 예외 rule은 Defeater로 등록
            Defeat(approve.IsBudgetFrozen, approve.IsUnderBudget). // Rebuttal → Warrant 공격 edge 필수
            Defeat(approve.IsUrgent, approve.IsBudgetFrozen),      // Defeater → Rebuttal 예외 처리
    ).
    Step("finance",
        toulmin.NewGraph("expense:finance").
            Warrant(hasFinance, 1.0).
            Warrant(approve.IsUnderBudget, 1.0).
            Rebuttal(approve.IsBudgetFrozen, 1.0).
            Defeater(approve.IsCEOOverride, 1.0).
            Defeat(approve.IsBudgetFrozen, approve.IsUnderBudget).
            Defeat(approve.IsCEOOverride, approve.IsBudgetFrozen),
    )

// 소액 자동 승인 — 별도 graph로 Flow 자체를 건너뛸 수도 있음
autoApprove := toulmin.NewGraph("expense:auto").
    Warrant(smallAmount, 1.0)

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
// result.Approved: 모든 단계 verdict > 0이면 true (undecided는 거부)
// result.Steps[0].Trace: 각 단계별 판정 근거
```

## 범위

### 포함

1. **ApprovalRequest, ApprovalContext 구조체**: 승인 판정에 필요한 요청/컨텍스트
2. **Approver, Budget 구조체, OrgTree 인터페이스**: 조직 구조 추상화
3. **범용 rule 함수**: IsDirectManager, IsUnderBudget, IsBudgetFrozen, HasApprovalRole, IsAboveLevel, IsSmallAmount, IsUrgent, IsCEOOverride
4. **Flow**: 다단계 승인 흐름 (Step 등록, 순차 판정)
5. **FlowResult, StepResult**: 단계별 판정 결과 + trace
6. **테스트**: rule 함수 단위 테스트, Flow 통합 테스트

### 제외

- 승인 상태 퍼시스턴스 (DB 저장은 사용자 책임)
- 알림/이메일 발송
- UI/대시보드
- 병렬 승인 (동시 결재) — 별도 Phase에서 검토

## 산출물

```
pkg/
  approve/
    approval_request.go           — ApprovalRequest 구조체
    approval_context.go           — ApprovalContext 구조체
    approver.go                   — Approver 구조체
    budget.go                     — Budget 구조체
    org_tree.go                   — OrgTree 인터페이스
    rule_is_direct_manager.go     — IsDirectManager
    rule_is_under_budget.go       — IsUnderBudget
    rule_is_budget_frozen.go      — IsBudgetFrozen
    rule_has_approval_role.go     — HasApprovalRole (클로저)
    rule_is_above_level.go        — IsAboveLevel (클로저)
    rule_is_small_amount.go       — IsSmallAmount (클로저)
    rule_is_urgent.go             — IsUrgent
    rule_is_ceo_override.go       — IsCEOOverride
    flow.go                       — Flow (NewFlow, Step, Evaluate)
    flow_result.go                — FlowResult, StepResult
    rule_test.go                  — rule 함수 단위 테스트
    flow_test.go                  — Flow 통합 테스트
```

## 단계

### Step 1: 구조체 및 인터페이스 정의

- ApprovalRequest, ApprovalContext, Approver, Budget, OrgTree
- FlowResult, StepResult

### Step 2: rule 함수 구현

- 각 rule 함수를 파일 하나에 하나씩 구현 (filefunc 규칙 준수)
- 각 함수는 1-2 depth 유지
- 클로저 rule: HasApprovalRole, IsAboveLevel, IsSmallAmount

### Step 3: Flow 구현

- NewFlow: 빈 Flow 생성
- Step: graph를 단계로 등록
- Evaluate: 단계를 순차 실행, 각 단계마다 StepContextFunc으로 ground 주입, 하나라도 verdict < 0이면 거부

### Step 4: 테스트

- rule 함수 단위 테스트: 각 조건별 true/false
- Flow 통합 테스트:
  - 모든 단계 통과 → Approved: true
  - 중간 단계 거부 → Approved: false, 거부 단계 식별
  - defeat edge 동작 (예산 동결 + 긴급 → 승인)
  - CEO 오버라이드 동작
  - 소액 자동 승인 graph 동작

### Step 5: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. IsDirectManager, IsUnderBudget 등 rule 함수가 claim/ground에서 올바르게 판정한다
2. Flow.Evaluate가 모든 단계를 순차 실행하고 결과를 집계한다
3. 중간 단계에서 거부되면 이후 단계를 실행하지 않는다
4. StepContextFunc으로 단계별 다른 ground를 주입할 수 있다
5. defeat edge가 예외를 정확히 처리한다 (긴급 → 예산 동결 무시)
6. FlowResult.Steps에 각 단계의 trace가 포함된다
7. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
