# Phase 005: Lazy Graph Evaluation — 그래프 순회 기반 규칙 실행 (구현 완료)

## 목표

규칙 실행을 일괄에서 그래프 순회 방식으로 전환한다.
warrant부터 시작하여 그래프를 따라 필요한 규칙만 실행하고,
trace에는 관련 규칙만 포함한다. verdict 결과는 동일하다.

## 배경

현재 Evaluate는 모든 규칙을 먼저 실행한 뒤 서브그래프를 구성한다.
warrant가 false여도 rebuttal을 실행하고, trace에 무관한 규칙도 포함된다.
h-Categoriser가 이미 재귀 구조이므로, 재귀 탐색 시점에 func을 실행하면
필요한 규칙만 평가하고 trace도 자연스럽게 관련된 것만 남는다.

## 핵심 설계

### 실행 흐름 변경

**현재 (일괄 실행)**
```
1. 모든 규칙 func 실행 → true/false 수집
2. true인 것만 모아서 서브그래프 구성
3. h-Categoriser로 verdict 계산
```

**변경 (그래프 순회)**
```
1. warrant 노드에서 시작
2. func 실행 → false면 스킵
3. true면 attackers 탐색, 각 attacker func 실행
4. 재귀적으로 h-Categoriser + func 실행 동시 수행
5. 실행된 규칙만 trace에 기록
```

### verdict 동일성

rule 함수는 `func(claim, ground) bool` — claim/ground에만 의존한다.
다른 규칙의 결과에 의존하지 않으므로 실행 순서가 verdict에 영향을 주지 않는다.

## 범위

### 포함

1. **GraphBuilder.Evaluate**: 그래프 순회 기반으로 변경
2. **GraphBuilder.EvaluateTrace**: 그래프 순회 + trace 수집
3. **기존 Engine API**: 그래프 순회 방식 적용 (Engine은 defeats가 RuleMeta에 있으므로 별도 처리)
4. **테스트**: 기존 verdict 동일성 검증 + 불필요한 규칙 미실행 검증
5. **문서 갱신**: manual-for-ai.md, README.md

### 제외

- h-Categoriser 알고리즘 변경 없음 — 실행 시점만 변경

## 산출물

```
pkg/
  toulmin/
    graph_builder_evaluate.go       — 그래프 순회 기반 Evaluate
    graph_builder_evaluate_trace.go — 그래프 순회 기반 EvaluateTrace
    engine_evaluate.go              — 그래프 순회 기반 Evaluate
    engine_evaluate_trace.go        — 그래프 순회 기반 EvaluateTrace
    graph_builder_test.go           — 미실행 검증 테스트 추가
```

## 단계

### Step 1: GraphBuilder.Evaluate 그래프 순회

- warrant 노드(attackers에 포함되지 않는 노드)를 시작점으로
- 각 노드에서 func 실행, false면 verdict 기여 없음
- true면 attackers 재귀 탐색 + h-Categoriser 연산
- 이미 평가한 노드는 캐시 (중복 실행 방지)

### Step 2: GraphBuilder.EvaluateTrace 그래프 순회 + trace

- Step 1과 동일하되 실행된 규칙마다 TraceEntry 수집

### Step 3: Engine.Evaluate/EvaluateTrace 적용

- Engine도 동일한 패턴으로 변경

### Step 4: 테스트

- 기존 verdict 동일성 검증 (모든 기존 테스트 PASS)
- 불필요한 규칙 미실행 검증: warrant=false일 때 rebuttal func 호출 안 됨
- trace에 관련 규칙만 포함되는지 검증

### Step 5: 문서 갱신

- manual-for-ai.md — Evaluation Flow 섹션 업데이트
- README.md

## 검증 기준

1. 기존 모든 테스트 PASS (verdict 동일)
2. warrant=false일 때 rebuttal func이 실행되지 않는다
3. trace에 실행된 규칙만 포함된다
4. 중복 노드는 한 번만 실행된다 (캐시)
5. filefunc validate 0 위반

## 의존성

- 없음
