# Request 001: Rule 시그니처 확장 — Evidence 반환

## 요청자

filefunc (github.com/park-jun-woo/filefunc)

## 문제

현재 rule 시그니처 `func(claim any, ground any) bool`은 판정만 반환한다. 위반의 **구체적 사유**(어디의 뭐가 왜 틀렸는지)를 전달할 수 없다.

filefunc에서 발생하는 실제 문제:

1. toulmin이 verdict > 0 (위반)을 판정
2. 구체적 사유(예: "depth 4 exceeds maximum of 2")를 알려면 원본 Check 함수를 **다시 호출**해야 함
3. 동일 로직이 2번 실행됨 — Rule 함수(판정) + Check 함수(메시지)

## 요청

### rule 시그니처 확장

```go
// 현재
func(claim any, ground any) bool

// 요청
func(claim any, ground any) (bool, any)
```

두 번째 반환값 `any`는 rule이 남기는 증거(evidence). 도메인별로 다른 타입. toulmin은 내용을 모르고 전달만 한다.

### TraceEntry에 Evidence 필드 추가

```go
// 현재
type TraceEntry struct {
    Name      string
    Role      string
    Activated bool
    Qualifier float64
}

// 요청
type TraceEntry struct {
    Name      string
    Role      string
    Activated bool
    Qualifier float64
    Evidence  any       // rule 함수가 반환한 증거. nil이면 없음.
}
```

### EvalResult에 Evidence 필드 추가

```go
// 현재
type EvalResult struct {
    Name    string
    Verdict float64
}

// 요청
type EvalResult struct {
    Name     string
    Verdict  float64
    Evidence any     // warrant의 증거 (verdict > 0일 때 유의미)
}
```

## 사용 예시 (filefunc)

```go
// Rule 함수 — 판정 + 증거를 동시에 반환
func RuleF1(claim any, ground any) (bool, any) {
    gf := ground.(*ValidateGround).File
    if len(gf.Funcs) > 1 {
        return true, &Evidence{
            File: gf.Path, Rule: "F1",
            Got: len(gf.Funcs), Expected: 1,
        }
    }
    return false, nil
}

// 소비 — EvalResult에서 증거를 꺼냄
results := ValidateGraph.EvaluateTrace(claim, ground)
for _, r := range results {
    if r.Verdict > 0 {
        ev := r.Evidence.(*Evidence)
        fmt.Printf("[%s] %s: got %v, expected %v\n", ev.Rule, ev.File, ev.Got, ev.Expected)
    }
}
```

## 영향 범위

### 변경 필요

| 파일 | 변경 내용 |
|---|---|
| `pkg/toulmin/rule_meta.go` | `Fn` 필드 타입 `func(any, any) bool` → `func(any, any) (bool, any)` |
| `pkg/toulmin/eval_result.go` | `Evidence any` 필드 추가 |
| `pkg/toulmin/trace_entry.go` | `Evidence any` 필드 추가 |
| `pkg/toulmin/engine_evaluate.go` | `r.Fn(claim, ground)` 반환값 2개 처리 |
| `pkg/toulmin/engine_evaluate_trace.go` | 동일 + trace에 evidence 저장 |
| `pkg/toulmin/graph_builder_evaluate.go` | 동일 |
| `pkg/toulmin/graph_builder_evaluate_trace.go` | 동일 |
| `pkg/toulmin/build_subgraph.go` | activated 수집 시 evidence 보존 |
| `pkg/toulmin/engine_test.go` | 테스트 rule 함수 시그니처 변경 |
| `pkg/toulmin/graph_builder_test.go` | 동일 |

### 하위 호환성

**깨짐.** `func(any, any) bool` → `func(any, any) (bool, any)`은 시그니처 변경. 기존 사용자의 rule 함수가 컴파일 실패.

마이그레이션: 기존 `return true` → `return true, nil`. 증거가 필요 없으면 `nil`.

### 대안: 하위 호환 유지

두 시그니처를 모두 지원하는 방법:

```go
type RuleMeta struct {
    Fn       func(any, any) bool         // 기존 (evidence 없음)
    FnTrace  func(any, any) (bool, any)  // 확장 (evidence 있음)
}
```

`FnTrace`가 있으면 우선 사용, 없으면 `Fn` 사용. 하위 호환 유지.

하지만 이러면 RuleMeta가 복잡해지고, Graph Builder의 Warrant/Defeater/Rebuttal도 두 시그니처를 받아야 함. 복잡도 증가.

## 권장

하위 호환을 깨고 시그니처를 일괄 변경하는 것을 권장. 이유:

1. toulmin은 아직 초기 단계 (사용자: filefunc 1개)
2. 마이그레이션 비용이 낮음 (`return true` → `return true, nil`)
3. 깨끗한 단일 시그니처 유지
4. evidence가 `any`이므로 도메인 무관성 유지
