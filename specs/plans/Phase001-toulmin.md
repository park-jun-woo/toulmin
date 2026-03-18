# Phase 001: toulmin 엔진 코어 + CLI

## 목표

Toulmin 논증 모델 기반 규칙 엔진의 최소 동작 구현.
Rule 등록 → defeats 그래프 구성 → h-Categoriser verdict 연산 → CLI 출력.

## 범위

### 포함

1. **도메인 모델**: Claim, Ground, Rule, Node, RuleGraph
2. **엔진 코어**: rule 평가, 서브그래프 구성, h-Categoriser verdict [-1, 1]
3. **어노테이션 파서**: `//rule:` 접두사 메타데이터 파싱
4. **CLI**: `toulmin evaluate` 커맨드 (cobra)

### 불필요 (구현 대상 아님)

- Ground Adapter — ground는 호출부가 넣는 값일 뿐, 엔진이 추상화할 대상 아님
- 부분 평가 (Partial Evaluation) — 엔진은 bool 결과와 그래프만 필요
- phase(시점) 선택 — 같은 함수를 언제 호출하느냐의 문제, 엔진 관심사 밖

## 산출물

```
go.mod
go.sum
codebook.yaml
cmd/
  toulmin/
    main.go                  — 엔트리포인트
internal/
  cli/
    root.go                  — cobra root command
    evaluate.go              — cobra evaluate command
pkg/
  toulmin/
    claim.go                 — PatternClaim, RelationClaim
    ground.go                — Ground interface
    rule.go                  — Rule type, RuleMeta
    node.go                  — Node (qualifier, strength)
    strength.go              — Strict, Defeasible, Defeater 상수
    graph.go                 — RuleGraph (nodes, edges, Attackers)
    graph_build.go           — BuildSubgraph (활성 rule → 서브그래프)
    verdict.go               — CalcAcceptability (h-Categoriser)
    engine.go                — Engine (rule 등록, 평가 오케스트레이션)
    annotation.go            — //rule: 파서
```

## 단계

### Step 1: 프로젝트 초기화

- go mod init
- cobra 의존성 추가
- main.go + root command

### Step 2: 도메인 모델

- Claim 타입 (PatternClaim, RelationClaim)
- Ground interface
- Rule 타입: `func(claim any, ground any) bool`
- RuleMeta: name, qualifier, strength, defeats, backing, what
- Node: RuleMeta + 평가 결과
- Strength 상수: Strict, Defeasible, Defeater

### Step 3: 엔진 코어

- RuleGraph: nodes map, defeats edges, Attackers(nodeID) 메서드
- BuildSubgraph: 활성 rule + defeats 간선 → 서브그래프 구성 (strict 노드는 들어오는 간선 거부)
- CalcAcceptability: h-Categoriser 재귀 연산, maxDepth=100, [-1, 1] 변환
- Engine: rule 등록, claim에 대해 전체 평가 흐름 오케스트레이션

### Step 4: 어노테이션 파서

- `//rule:warrant qualifier=1.0 strength=strict` 파싱
- `//rule:defeater defeats=CheckOneFileOneFunc` 파싱
- `//rule:backing "..."` 파싱
- `//rule:what ...` 파싱

### Step 5: CLI

- `toulmin evaluate` — 등록된 rule로 claim 평가, verdict 출력
- JSON 출력 포맷 (verdict, activated rules, graph)

## 검증 기준

1. warrant만 있는 경우: verdict = +1.0 (위반 확정)
2. warrant + defeater: verdict = 0.0 (판정불가)
3. warrant + defeater + defeater의 defeater: verdict > 0 (부분 복원, 보상 원리)
4. strict warrant에 공격 시도: 간선 거부, verdict 변화 없음
5. 순환 공격: maxDepth에서 0.0 반환

## 의존성

- github.com/spf13/cobra
