# Phase 004: Evaluation Trace — 판정 사유 추적

## 목표

Evaluate 결과에 trace를 추가하여 판정 사유를 설명 가능하게 한다.
각 규칙의 이름, 역할, 활성화 여부, qualifier를 verdict와 함께 반환한다.

## 배경

현재 EvalResult는 verdict만 반환한다.
filefunc validate 등 소비자가 "왜 이 verdict인가"를 설명하려면
어떤 규칙이 true/false였고, 어떤 qualifier로 적용됐는지 알아야 한다.

## 범위

### 포함

1. **TraceEntry 구조체**: name, role, activated, qualifier
2. **EvalResult에 Trace 필드 추가**
3. **GraphBuilder.Evaluate에서 trace 수집**: 모든 규칙(true/false 포함) 기록
4. **Engine.Evaluate에서도 trace 수집**: 기존 API 호환
5. **문서 갱신**: manual-for-ai.md, README.md

### 제외

- trace 포맷팅/메시지 조립 — 소비자 책임

## 설계

### TraceEntry

```go
type TraceEntry struct {
    Name      string  `json:"name"`
    Role      string  `json:"role"`      // "warrant", "rebuttal", "defeater"
    Activated bool    `json:"activated"` // func(claim, ground) 결과
    Qualifier float64 `json:"qualifier"` // 적용된 확신도
}
```

### EvalResult 변경

```go
type EvalResult struct {
    Name    string       `json:"name"`
    Verdict float64      `json:"verdict"`
    Trace   []TraceEntry `json:"trace"`
}
```

### 출력 예시

```json
{
  "name": "IsAdult",
  "verdict": 1.0,
  "trace": [
    {"name": "IsAdult", "role": "warrant", "activated": true, "qualifier": 1.0},
    {"name": "HasCriminalRecord", "role": "rebuttal", "activated": false, "qualifier": 1.0}
  ]
}
```

## 산출물

```
pkg/
  toulmin/
    trace_entry.go            — TraceEntry 구조체
    eval_result.go            — EvalResult (Trace 필드 추가)
    graph_builder_evaluate.go — trace 수집 로직 추가
    engine_evaluate.go        — trace 수집 로직 추가
    engine_test.go            — trace 검증 테스트 추가
    graph_builder_test.go     — trace 검증 테스트 추가
```

## 단계

### Step 1: TraceEntry 구조체

- `trace_entry.go` 생성

### Step 2: EvalResult에 Trace 필드 추가

- `eval_result.go` 수정

### Step 3: GraphBuilder.Evaluate에서 trace 수집

- 모든 규칙(activated + non-activated) 순회
- 각 규칙의 name, role, activated, qualifier 기록
- 역할(role)은 builder 등록 시점의 메서드(Warrant/Rebuttal/Defeater)로 결정

### Step 4: Engine.Evaluate에서 trace 수집

- 기존 API에도 동일하게 trace 추가
- 역할은 RuleMeta의 Defeats/Strength로 추론

### Step 5: 테스트

- trace에 모든 규칙(true/false)이 포함되는지 검증
- trace의 role이 정확한지 검증
- trace의 qualifier가 정확한지 검증

### Step 6: 문서 갱신

- `manual-for-ai.md` — TraceEntry 설명, 출력 예시 추가
- `README.md` — trace 사용 예시 추가

## 검증 기준

1. EvalResult.Trace에 해당 그래프의 모든 규칙이 포함된다
2. activated=true인 규칙과 false인 규칙이 모두 포함된다
3. role이 warrant/rebuttal/defeater로 정확히 표시된다
4. qualifier 값이 builder/engine 등록 시 지정한 값과 일치한다
5. 기존 verdict 값에 영향 없음

## 의존성

- 없음 (기존 구조체 확장)
