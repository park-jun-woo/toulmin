# Phase 007: 공격자 목록 캐시 + funcID full path 도입

## 목표

공격자 목록(attackerSet)을 evalContext 생성 시 한 번만 만들어 캐시하여 isWarrant과 calcTrace의 반복 순회를 제거한다.
내부 규칙 식별자를 package path 포함 full name으로 변경하여 함수 이름 충돌을 원천 방지한다.
Evaluate와 EvaluateTrace의 verdict 차이를 테스트로 명시한다.

## 배경

Phase 006에서 evalContext 추출로 4중 중복을 제거했으나,
성능과 안전성 측면에서 개선할 부분이 남아있다.

### 현재 문제

1. **공격자 목록 반복 생성**: `calcTrace`가 호출될 때마다 전체 edges를 순회하여 "누가 공격자인가" 목록을 만든다. edges는 평가 중 변하지 않으므로 한 번만 만들어 캐시하면 된다. 현재는 재귀 호출마다 반복하여 O(N×E) 복잡도
2. **isWarrant O(E) 순회**: Evaluate 루프 안에서 매 rule마다 전체 edges를 순회하여 attacker 여부를 확인한다. 공격자 목록을 캐시하면 O(1) 조회로 대체 가능
3. **FuncName short name 충돌**: 현재 `FuncName`이 `strings.LastIndex(full, ".")`로 short name만 추출한다. 다른 패키지의 동명 함수(`pkgA.IsAdult` vs `pkgB.IsAdult`)나 익명 함수(`func1`, `func2`)가 동일 이름으로 등록되어 규칙이 덮어써지고 잘못된 verdict를 반환할 수 있다
4. **Evaluate vs EvaluateTrace verdict 차이 미검증**: Evaluate는 warrant 간 캐시를 공유하고, EvaluateTrace는 warrant마다 reset한다. 복수 warrant가 동일 attacker를 공유하면 verdict가 달라질 수 있으나 이에 대한 테스트가 없다

## 핵심 설계

### 공격자 목록을 evalContext 필드로 캐시

edges의 value 쪽에 등장하는 노드 = 공격자. 이 목록은 graph가 변하지 않는 한 고정이므로 한 번만 만들어 캐시한다.

```go
type evalContext struct {
    // ... 기존 필드
    attackerSet map[string]bool // 신규: 캐시된 공격자 목록
}
```

`newEvalContext`에서 edges 구성 직후 한 번만 생성한다:

```go
ctx.attackerSet = make(map[string]bool)
for _, attackers := range ctx.edges {
    for _, aid := range attackers {
        ctx.attackerSet[aid] = true
    }
}
```

### isWarrant 시그니처 변경

```go
// 변경 전
func isWarrant(edges map[string][]string, strength Strength, name string) bool

// 변경 후
func isWarrant(attackerSet map[string]bool, strength Strength, name string) bool
```

edges 전체를 순회하지 않고 `attackerSet[name]` O(1) 조회로 대체.

### calcTrace에서 attackerSet 재사용

```go
// 변경 전: 매 호출마다 attackerSet 생성
attackerSet := make(map[string]bool)
for _, attackers := range ctx.edges {
    for _, aid := range attackers {
        attackerSet[aid] = true
    }
}
role = inferRole(ctx.strMap, attackerSet, id)

// 변경 후: ctx.attackerSet 직접 사용
role = inferRole(ctx.strMap, ctx.attackerSet, id)
```

### 내부 식별자를 full path로 변경 — funcID 도입

내부 규칙 식별에 full path를 사용하고, 표시용에는 기존 short name을 유지한다.

```go
// funcID — 내부 식별용 (full path, 고유 보장)
func funcID(fn func(any, any) (bool, any)) string {
    return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// FuncName — 표시용 (short name, trace/JSON 출력)
func FuncName(fn func(any, any) (bool, any)) string {
    // 기존 로직 유지: LastIndex(".") 이후 반환
}
```

**변경 범위:**
- `evalContext` 내부 맵 키: funcID (full path)
- `GraphBuilder.Defeat`: funcID로 edge 연결
- `GraphBuilder.Warrant/Rebuttal/Defeater`: funcID로 등록
- `TraceEntry.Name`, `EvalResult.Name`: FuncName (short name) — 표시용 유지
- `Engine.Register`: 사용자가 직접 Name을 지정하므로 변경 없음

이렇게 하면 다른 패키지의 동명 함수(`pkgA.IsAdult` vs `pkgB.IsAdult`)와
익명 함수(`pkg.TestFoo.func1` vs `pkg.TestFoo.func2`)가 자연스럽게 구분된다.
중복 검사 panic이 불필요해진다.

### Evaluate vs EvaluateTrace 동작 차이 문서화

verdict 차이는 의도된 설계이다:
- `Evaluate`: 캐시 공유로 실행 효율 우선
- `EvaluateTrace`: per-warrant 격리로 정확한 trace 우선

테스트를 추가하여 이 차이를 명시적으로 검증한다.

## 범위

### 포함

1. **evalContext에 공격자 목록 캐시 필드 추가**: `eval_context.go`
2. **newEvalContext에서 공격자 목록 한 번 생성**: `new_eval_context.go`
3. **isWarrant 시그니처 변경**: edges → attackerSet 조회
4. **calcTrace에서 ctx.attackerSet 사용**: 반복 생성 제거
5. **Evaluate/EvaluateTrace 호출부 수정**: isWarrant 호출 시 ctx.attackerSet 전달
6. **funcID 도입**: 내부 식별자를 full path로 변경, FuncName은 표시용으로 유지
7. **GraphBuilder 등록/연결을 funcID 기반으로 변경**: Warrant/Rebuttal/Defeater/Defeat
8. **TraceEntry.Name에 short name 유지**: funcID→FuncName 변환
9. **테스트 추가**: 다른 패키지 동명 함수 구분, Evaluate vs EvaluateTrace 차이, 경계값

### 제외

- h-Categoriser 수식 변경 없음
- 외부 API 시그니처 변경 없음
- maxDepth 초과 시 경고/에러 반환 — 별도 Phase에서 검토

## 산출물

```
pkg/
  toulmin/
    eval_context.go               — attackerSet 필드 추가
    new_eval_context.go           — 공격자 목록 캐시 추가
    eval_context_calc_trace.go    — ctx.attackerSet 사용으로 변경
    is_warrant.go                 — 시그니처 변경 (attackerSet 조회)
    engine_evaluate.go            — isWarrant 호출부 수정
    engine_evaluate_trace.go      — isWarrant 호출부 수정
    graph_builder_evaluate.go     — isWarrant 호출부 수정
    graph_builder_evaluate_trace.go — isWarrant 호출부 수정
    func_name.go                  — FuncName은 표시용으로 유지
    func_id.go                    — funcID 신규 (내부 식별용 full path)
    graph_builder_warrant.go      — funcID 기반 등록으로 변경
    graph_builder_rebuttal.go     — funcID 기반 등록으로 변경
    graph_builder_defeater.go     — funcID 기반 등록으로 변경
    graph_builder_defeat.go       — funcID 기반 edge 연결로 변경
    engine_test.go                — Evaluate vs EvaluateTrace 차이 테스트 추가
    graph_builder_test.go         — 동명 함수 구분 테스트, 경계값 테스트 추가
```

## 단계

### Step 1: evalContext에 공격자 목록 캐시

- `eval_context.go`에 `attackerSet map[string]bool` 필드 추가
- `new_eval_context.go`에서 edges 구성 직후 공격자 목록을 한 번 생성하여 캐시

### Step 2: isWarrant 시그니처 변경

- `is_warrant.go`: `edges map[string][]string` → `attackerSet map[string]bool`
- 내부 로직: `attackerSet[name]` O(1) 조회로 변경
- 호출부 4곳 수정: Engine.Evaluate, Engine.EvaluateTrace, GraphBuilder.Evaluate, GraphBuilder.EvaluateTrace

### Step 3: calcTrace에서 attackerSet 재사용

- `eval_context_calc_trace.go`에서 attackerSet 생성 코드 제거
- `ctx.attackerSet`을 `inferRole`에 직접 전달

### Step 4: funcID 도입 + GraphBuilder 변경

- `func_id.go` 신규: `funcID(fn) string` — `runtime.FuncForPC`의 full name 반환
- `func_name.go` 유지: `FuncName(fn) string` — 기존 short name (표시용)
- `graph_builder_warrant.go`, `graph_builder_rebuttal.go`, `graph_builder_defeater.go`: `FuncName` → `funcID`로 변경 (RuleMeta.Name에 full path 저장)
- `graph_builder_defeat.go`: `FuncName` → `funcID`로 변경
- `graph_builder_evaluate.go`, `graph_builder_evaluate_trace.go`: EvalResult.Name과 TraceEntry.Name 출력 시 short name 변환 (full path에서 LastIndex(".") 이후 추출)

### Step 5: 테스트 추가

- **동명 함수 구분 테스트**: 같은 short name이지만 다른 full path인 함수가 올바르게 구분되는지 검증
- **Evaluate vs EvaluateTrace 차이 테스트**: 복수 warrant가 동일 attacker를 공유하는 경우
- **qualifier 0.0 테스트**: 경계값 검증
- **깊은 체인 테스트**: 10단계 이상 defeat 체인에서 verdict 수렴 검증

### Step 6: 기존 테스트 PASS 확인

- 모든 기존 테스트가 변경 없이 PASS하는지 확인

## 검증 기준

1. 기존 모든 테스트 PASS (verdict 동일)
2. `isWarrant`가 O(1)로 동작한다 (edges 순회 없음)
3. `calcTrace`에서 attackerSet을 매번 생성하지 않는다
4. 다른 패키지의 동명 함수가 funcID로 구분되어 서로 다른 규칙으로 등록된다
5. 다른 graph에서 동일 함수를 재사용하면 정상 동작한다 (TestGraphBuilderFuncReuse PASS)
6. TraceEntry.Name, EvalResult.Name은 short name으로 표시된다
7. Evaluate와 EvaluateTrace의 동작 차이가 테스트로 문서화된다
7. filefunc validate 0 위반

## 의존성

- 없음 (내부 최적화 + 방어 로직)
