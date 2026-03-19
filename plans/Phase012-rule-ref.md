# Phase 012: Rule 참조 반환 + 체이닝 제거 — 코어 리팩토링

## 목표

Warrant/Rebuttal/Defeater가 `*Rule` 참조를 반환하도록 변경한다.
Defeat는 `*Rule` 두 개를 받아 관계만 선언한다.
체이닝(`*GraphBuilder` 반환)을 제거하고, 정의와 관계를 분리한다.

## 배경

### 현재 문제

1. **Defeat가 rule 정의를 재구성한다**: `Defeat(IsIPInList, blocklist, IsAuthenticated, nil)` — 이미 등록된 rule의 fn+backing을 다시 적어야 한다
2. **체이닝이 읽기를 해친다**: 한 줄에 정의+관계가 섞여 graph 구조가 한눈에 안 들어온다
3. **같은 함수 + 다른 backing 식별이 번거롭다**: ruleID를 내부에서 재구성해야 하므로 DefeatWith 같은 별도 메서드가 필요했다

### 해결: 정의할 때 정의하고, 가리킬 때 가리킨다

```go
// 정의 — 각 rule이 참조를 반환
auth    := g.Warrant(IsAuthenticated, nil, 1.0)
admin   := g.Warrant(IsInRole, "admin", 1.0)
blocked := g.Rebuttal(IsIPInList, blocklist, 1.0)
allowed := g.Defeater(IsIPInList, whitelist, 1.0)

// 관계 — 참조로 가리킴, 반복 없음
g.Defeat(blocked, auth)
g.Defeat(allowed, blocked)
```

## 핵심 설계

### Rule 구조체

```go
// pkg/toulmin/rule.go
type Rule struct {
    id string  // ruleID (funcID + "#" + backing)
}
```

외부에 노출되는 것은 타입뿐이다. id는 비공개. Rule은 Defeat에서 참조로만 사용된다.

### GraphBuilder API 변경

```go
// 현재 — 체이닝, *GraphBuilder 반환
func (b *GraphBuilder) Warrant(fn any, backing any, qualifier float64) *GraphBuilder

// 변경 — *Rule 반환, 체이닝 제거
func (g *Graph) Warrant(fn any, backing any, qualifier float64) *Rule
func (g *Graph) Rebuttal(fn any, backing any, qualifier float64) *Rule
func (g *Graph) Defeater(fn any, backing any, qualifier float64) *Rule
func (g *Graph) Defeat(from *Rule, to *Rule)
```

Defeat는 반환값 없음 — 관계 선언이므로.

### GraphBuilder → Graph 이름 변경

체이닝이 사라지면 "Builder" 패턴이 아니다. `GraphBuilder` → `Graph`로 이름을 변경한다.

```go
// 현재
g := toulmin.NewGraph("voting")

// 변경 없음 — NewGraph는 그대로, 반환 타입만 *Graph
g := toulmin.NewGraph("voting")
```

### NewGraph 반환 타입

```go
func NewGraph(name string) *Graph
```

### Evaluate/EvaluateTrace 위치

`*Graph`의 메서드로 유지:

```go
results, err := g.Evaluate(claim, ground)
results, err := g.EvaluateTrace(claim, ground)
```

### 사용 예시

```go
g := toulmin.NewGraph("route:admin")

// 정의
auth    := g.Warrant(route.IsAuthenticated, nil, 1.0)
admin   := g.Warrant(route.IsInRole, "admin", 1.0)
blocked := g.Rebuttal(route.IsIPInList, blocklist, 1.0)
limited := g.Rebuttal(route.IsRateLimited, limiter, 1.0)
allowed := g.Defeater(route.IsIPInList, whitelist, 1.0)
internal := g.Defeater(route.IsInternalService, nil, 1.0)

// 관계
g.Defeat(blocked, auth)
g.Defeat(limited, auth)
g.Defeat(allowed, blocked)
g.Defeat(internal, limited)

// 평가
results, err := g.Evaluate(nil, ctx)
```

### 기존 테스트 변환 예시

```go
// 현재 (체이닝)
g := NewGraph("test").
    Warrant(WarrantA, nil, 1.0).
    Rebuttal(RebuttalB, nil, 1.0).
    Defeat(RebuttalB, nil, WarrantA, nil)

// 변경 (참조)
g := NewGraph("test")
w := g.Warrant(WarrantA, nil, 1.0)
r := g.Rebuttal(RebuttalB, nil, 1.0)
g.Defeat(r, w)
```

### 하위 호환성

체이닝 API를 제거하므로 **breaking change**다. 하위 호환성은 유지하지 않는다.
- pkg/route: 새 API로 전환
- internal/cli: 새 API로 전환
- internal/codegen: 생성 코드 템플릿 변경
- 기존 사용자: API 변경 필요

## 범위

### 포함

1. **Rule 구조체 신규**: id만 가진 참조 타입
2. **GraphBuilder → Graph 이름 변경**: 구조체, 파일명, 메서드 수신자 전부
3. **Warrant/Rebuttal/Defeater 반환 타입**: `*GraphBuilder` → `*Rule`
4. **Defeat 시그니처**: `(fromFn, fromBacking, toFn, toBacking)` → `(from *Rule, to *Rule)`
5. **DefeatWith 제거**: Defeat가 *Rule로 통합
6. **기존 테스트 전면 수정**: 체이닝 → 참조 방식
7. **pkg/route 수정**: 새 API 적용
8. **internal/cli, internal/codegen 수정**: 새 API 적용

### 제외

- Engine API (RuleMeta 직접 등록 방식) — 별도 검토
- YAML 코드젠 출력 형식 — 별도 Phase

## 산출물

```
pkg/
  toulmin/
    rule.go                       — Rule 구조체 (신규)
    graph.go                      — Graph 구조체 (GraphBuilder에서 이름 변경)
    new_graph.go                  — NewGraph → *Graph
    graph_warrant.go              — Warrant → *Rule (이름/반환 변경)
    graph_rebuttal.go             — Rebuttal → *Rule
    graph_defeater.go             — Defeater → *Rule
    graph_defeat.go               — Defeat(from *Rule, to *Rule)
    graph_evaluate.go             — Evaluate (수신자 변경)
    graph_evaluate_trace.go       — EvaluateTrace (수신자 변경)
    *_test.go                     — 전면 수정
  route/
    guard.go                      — *Graph 타입 참조
    guard_debug.go                — *Graph 타입 참조
    guard_test.go                 — 새 API 적용
internal/
  codegen/
    generate_graph.go             — 참조 방식 코드 생성
    generate_graph_test.go        — 기대값 변경
    format_defeats.go             — 변수명 참조
  analyzer/
    extract_defeats.go            — AST 분석 패턴 변경
    extract_defeats_test.go       — 테스트 수정
  cli/
    run_evaluate.go               — 새 API 적용
```

## 단계

### Step 1: Rule 구조체 생성

- `pkg/toulmin/rule.go`: `type Rule struct { id string }`

### Step 2: GraphBuilder → Graph 이름 변경

- 구조체명, 파일명, 모든 수신자 변경
- NewGraph 반환 타입 변경

### Step 3: Warrant/Rebuttal/Defeater 반환 변경

- `*GraphBuilder` → `*Rule` 반환
- 내부에서 rule 등록 후 `&Rule{id: ruleID}` 반환

### Step 4: Defeat 시그니처 변경

- `(fromFn, fromBacking, toFn, toBacking)` → `(from *Rule, to *Rule)`
- DefeatWith 파일 삭제 (이미 삭제됨)
- from.id, to.id로 defeatEdge 생성

### Step 5: 코어 테스트 전면 수정

- 체이닝 → 참조 방식으로 전환
- 모든 기존 테스트 PASS 확인

### Step 6: pkg/route 수정

- Guard/GuardDebug: `*GraphBuilder` → `*Graph`
- guard_test.go: 새 API 적용

### Step 7: internal/cli 수정

- internal/cli/run_evaluate.go: 새 API 적용 (Engine API는 변경 없음, GraphBuilder 참조 제거)

### Step 8: internal/codegen 수정

`toulmin graph voting.yaml`로 생성되는 코드가 체이닝에서 참조 방식으로 변경된다:

```go
// 현재 생성 코드 (체이닝)
var VotingGraph = toulmin.NewGraph("voting").
    Warrant(IsAdult, nil, 1.0).
    Rebuttal(HasCriminalRecord, nil, 1.0).
    Defeat(HasCriminalRecord, nil, IsAdult, nil)

// 변경 생성 코드 (참조)
var VotingGraph = func() *toulmin.Graph {
    g := toulmin.NewGraph("voting")
    isAdult := g.Warrant(IsAdult, nil, 1.0)
    hasCriminalRecord := g.Rebuttal(HasCriminalRecord, nil, 1.0)
    g.Defeat(hasCriminalRecord, isAdult)
    return g
}()
```

수정 대상:
- `internal/codegen/generate_graph.go`: 참조 방식 코드 생성
- `internal/codegen/generate_graph_test.go`: 기대값 변경
- `internal/codegen/format_defeats.go`: Defeat가 변수명 참조
- `internal/analyzer/extract_defeats.go`: AST 분석 패턴 변경 (체이닝 → 분리 호출)

### Step 9: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS

## 검증 기준

1. Warrant/Rebuttal/Defeater가 `*Rule`을 반환한다
2. Defeat가 `*Rule` 두 개를 받아 관계를 선언한다
3. 같은 함수 + 다른 backing이 다른 `*Rule`로 구분된다
4. 체이닝 없이 정의와 관계가 분리된다
5. `*Graph`에서 Evaluate/EvaluateTrace가 정상 동작한다
6. 기존 모든 verdict 계산 결과가 동일하다
7. 전체 테스트 PASS

## 의존성

- Phase 010: backing 일급 시민
- Phase 011: pkg/route (함께 수정)
