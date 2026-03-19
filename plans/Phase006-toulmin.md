# Phase 006: Evaluate 공통 로직 추출 — 4중 중복 제거 (구현 완료)

## 목표

Engine.Evaluate, Engine.EvaluateTrace, GraphBuilder.Evaluate, GraphBuilder.EvaluateTrace에
분산된 동일 calc 재귀 로직을 공통 함수로 추출한다.
warrant 식별 기준을 통일하고, 존재하지 않는 rule 참조 시 panic을 방지한다.

## 배경

Phase 005에서 lazy graph evaluation을 도입하면서 calc 클로저가 4개 파일에 복제되었다.
h-Categoriser 수식(`raw = qualifier / (1 + Σattackers)`, `verdict = 2*raw - 1`)이
4곳에 동일하게 존재하여, 수식 변경 시 4곳을 모두 수정해야 한다.

또한 warrant 식별 기준이 Engine(`len(r.Defeats) > 0`)과
GraphBuilder(`attackerSet[r.Name]`)로 달라, 동일 규칙 구성에서 미묘한 동작 차이가 가능하다.

### 현재 문제

1. **calc 재귀 4중 복제**: 동일 수식이 engine_evaluate.go, engine_evaluate_trace.go, graph_builder_evaluate.go, graph_builder_evaluate_trace.go에 존재
2. **warrant 식별 기준 불일치**: Engine은 `Defeats` 필드 기반, GraphBuilder는 defeat edge의 `from` 등장 여부 기반
3. **nil function panic**: defeat edge가 등록되지 않은 rule을 참조하면 `fnMap[id]`가 nil이어서 panic 발생
4. **dead code**: `CalcAcceptability` + `BuildSubgraph`가 lazy evaluation 도입 후 엔진 내부에서 사용되지 않음

## 핵심 설계

### evalContext 내부 구조체

4개 Evaluate 함수가 공유하는 상태와 로직을 하나의 구조체로 캡슐화한다.

```go
type evalContext struct {
    fnMap    map[string]func(any, any) (bool, any)
    qualMap  map[string]float64
    strMap   map[string]Strength
    edges    map[string][]string
    ran      map[string]bool
    active   map[string]bool
    evidence map[string]any
}
```

### calc 메서드 통합

```go
func (ctx *evalContext) calc(id string, claim, ground any, depth int) float64
```

- depth >= maxDepth → 0.0
- fnMap에 id가 없으면 → -1.0 (panic 방지)
- 미실행이면 func 실행 + 캐시
- active=false → -1.0
- Strict이면 attacker 무시
- h-Categoriser 수식 적용

### warrant 식별 기준 통일

Engine과 GraphBuilder 모두 **edges 기반**으로 통일한다:
- warrant = edges의 from(공격자)으로 등장하지 않고, Strength != Defeater인 rule
- `isWarrant(edges, strength, name)` 헬퍼 함수 도입

### 상태 초기화 분리

- `Evaluate`: 모든 warrant 공유 캐시 (ran/active/evidence 1회 생성)
- `EvaluateTrace`: warrant마다 상태 리셋 (per-warrant trace)
- `evalContext.reset()` 메서드로 상태 초기화

## 범위

### 포함

1. **evalContext 구조체**: calc 재귀에 필요한 상태 캡슐화
2. **calc 메서드**: h-Categoriser + lazy func 실행 통합
3. **calcTrace 메서드**: calc + TraceEntry 수집
4. **isWarrant 헬퍼**: warrant 식별 기준 통일
5. **nil function guard**: fnMap에 없는 id 참조 시 -1.0 반환
6. **Engine.Evaluate/EvaluateTrace 리팩터링**: evalContext 사용
7. **GraphBuilder.Evaluate/EvaluateTrace 리팩터링**: evalContext 사용
8. **dead code 제거**: CalcAcceptability, BuildSubgraph, addDefeatEdges, RuleGraph.Attackers — 엔진 내부에서 미사용 시 제거 검토

### 제외

- h-Categoriser 수식 변경 없음 — 구조만 변경
- 외부 API 시그니처 변경 없음 — Engine.Evaluate, GraphBuilder.Evaluate의 반환 타입 동일

## 산출물

```
pkg/
  toulmin/
    eval_context.go               — evalContext 구조체 (신규)
    eval_context_calc.go          — evalContext.calc 메서드 (신규)
    eval_context_calc_trace.go    — evalContext.calcTrace 메서드 (신규)
    eval_context_reset.go         — evalContext.reset 메서드 (신규)
    new_eval_context.go           — newEvalContext 팩토리 (신규)
    is_warrant.go                 — isWarrant 헬퍼 (신규)
    engine_evaluate.go            — evalContext 사용으로 변경
    engine_evaluate_trace.go      — evalContext 사용으로 변경
    graph_builder_evaluate.go     — evalContext 사용으로 변경
    graph_builder_evaluate_trace.go — evalContext 사용으로 변경
    engine_test.go                — 기존 테스트 PASS 검증
    graph_builder_test.go         — 기존 테스트 PASS 검증 + nil guard 테스트 추가
```

### 제거 후보

```
pkg/
  toulmin/
    calc_acceptability.go         — inline calc로 대체됨
    build_subgraph.go             — lazy evaluation으로 불필요
    add_defeat_edges.go           — BuildSubgraph 전용이었음
    rule_graph_attackers.go       — BuildSubgraph + CalcAcceptability 전용이었음
    rule_graph.go                 — 위 함수들의 의존 타입
    node.go                       — RuleGraph 전용 타입
    new_graph.go (이름 충돌 주의) — GraphBuilder 팩토리와 별개
```

※ CLI의 `internal/graph/validate_graph.go`에서 RuleGraph를 사용하는지 확인 후 제거 결정

## 단계

### Step 1: evalContext 구조체 + 팩토리

- `eval_context.go`: 상태 필드 정의
- `new_eval_context.go`: `[]RuleMeta` + `[]defeatEdge` (또는 RuleMeta.Defeats) → evalContext 생성
- Engine용: RuleMeta.Defeats에서 edges 구성
- GraphBuilder용: defeatEdge 슬라이스에서 edges 구성

### Step 2: calc 메서드

- `eval_context_calc.go`: 재귀 calc 로직 1곳으로 통합
- nil function guard: `fnMap[id]`가 없으면 -1.0 반환
- 기존 4곳의 calc 클로저와 동일한 결과 보장

### Step 3: calcTrace 메서드

- `eval_context_calc_trace.go`: calc + TraceEntry 수집
- trace 슬라이스를 evalContext에 포함하거나 반환값으로 처리
- inferRole 또는 roles map 활용

### Step 4: reset 메서드

- `eval_context_reset.go`: ran, active, evidence, trace 초기화
- EvaluateTrace에서 warrant마다 호출

### Step 5: isWarrant 헬퍼

- `is_warrant.go`: edges와 strength 기반 warrant 판별
- Engine과 GraphBuilder 양쪽에서 동일 로직 사용

### Step 6: Engine.Evaluate/EvaluateTrace 리팩터링

- evalContext 생성 → warrant 순회 → calc/calcTrace 호출
- 기존 테스트 PASS 확인

### Step 7: GraphBuilder.Evaluate/EvaluateTrace 리팩터링

- evalContext 생성 → warrant 순회 → calc/calcTrace 호출
- 기존 테스트 PASS 확인

### Step 8: dead code 정리

- CalcAcceptability, BuildSubgraph, addDefeatEdges 사용처 확인
- internal/graph, internal/cli에서 미사용이면 제거
- 사용처가 있으면 유지하되, 엔진 평가 경로와의 관계를 문서에 명시

### Step 9: 테스트 추가

- nil function guard 테스트: defeat edge가 미등록 rule을 가리키는 경우
- warrant 식별 일관성 테스트: Engine과 GraphBuilder에서 동일 규칙 구성 시 동일 결과

### Step 10: 문서 갱신

- manual-for-ai.md — 내부 구조 변경 반영
- README.md — 외부 API 변경 없으므로 최소 수정

## 검증 기준

1. 기존 모든 테스트 PASS (verdict 동일)
2. h-Categoriser 수식이 1곳(evalContext.calc)에만 존재한다
3. Engine과 GraphBuilder에서 동일 규칙 구성 시 동일 verdict를 반환한다
4. defeat edge가 미등록 rule을 참조해도 panic이 발생하지 않는다
5. EvaluateTrace의 per-warrant 상태 리셋이 정상 동작한다
6. filefunc validate 0 위반

## 의존성

- 없음 (내부 리팩터링)
