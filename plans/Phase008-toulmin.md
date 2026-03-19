# Phase 008: 등록 시 순환 감지 + scanner 정리 + internal 테스트 보강

## 목표

순환 감지를 평가 시점(maxDepth)에서 그래프 구성 시점으로 이동하여
잘못된 그래프를 일찍 거부한다. calc에서 순환 방어 코드를 제거하여 평가 로직을 단순화한다.
Phase 003 도입 후 미사용된 scanner 패키지와 관련 dead code를 제거한다.
테스트가 없는 internal 패키지(graphdef, codegen, graph)에 단위 테스트를 추가한다.

## 배경

Phase 007까지 핵심 엔진의 알고리즘과 구조는 안정되었으나,
안정성과 품질 측면에서 미비점이 남아있다.

### 현재 문제

1. **순환 감지가 평가 시점에 있음**: `calc`/`calcTrace`가 `depth >= 100`이면 0.0을 반환한다. 순환은 **그래프 구조 오류**이지 평가 중에 처리할 런타임 상황이 아니다. 그래프 구성이 완료된 시점에 순환을 감지하면 (a) 사용자에게 명확한 에러를 반환하고 (b) calc에서 depth/visiting 방어 코드가 불필요해진다
2. **maxDepth=100이 깊은 비순환 체인을 자름**: 순환이 아닌 101단계 이상 체인도 0.0을 반환한다. 등록 시 순환을 거부하면 calc에서 깊이 제한이 불필요하므로 이 문제가 해소된다
3. **scanner 패키지 dead code**: Phase 002에서 도입한 `internal/scanner`(ScanDir, ExtractRules, extractRuleLines, RuleDecl)와 이를 사용하는 `internal/cli/scan_and_parse.go`, `internal/cli/parse_decls.go`가 Phase 003 이후 어디서도 호출되지 않는다. 미사용 코드가 남아있으면 유지보수 혼란과 컴파일 대상 증가를 유발한다
4. **internal 패키지 테스트 부재**: `internal/graphdef`(ParseYAML, Validate), `internal/codegen`(GenerateGraph), `internal/graph`(ValidateGraph)에 테스트가 없다. YAML 파싱 오류, 잘못된 defeat 참조, 코드 생성 결과 등 경계 케이스가 검증되지 않는다
5. **ParseYAML qualifier 0.0 구분 불가**: `qualifier == 0`이면 무조건 1.0으로 덮어쓴다. 사용자가 의도적으로 qualifier=0.0을 지정해도 1.0이 된다. YAML에서 미지정(zero value)과 명시적 0.0을 구분해야 한다
6. **funcID nil panic**: `funcID`에서 `runtime.FuncForPC`가 nil을 반환하면 `.Name()` 호출 시 panic이 발생한다. 방어 코드가 없다

## 핵심 설계

### 등록 시 순환 감지 — detectCycle

그래프의 edge map이 완성된 시점에 DFS로 순환을 감지한다.
순환이 발견되면 에러를 반환하여 평가를 차단한다.

```go
// detectCycle — DFS로 방향 그래프의 순환을 감지한다
func detectCycle(edges map[string][]string) error
```

DFS 상태: 미방문(0), 탐색 중(1), 완료(2). 탐색 중인 노드를 재방문하면 순환이다.

```go
func detectCycle(edges map[string][]string) error {
    state := make(map[string]int) // 0=unvisited, 1=visiting, 2=done
    var visit func(id string) error
    visit = func(id string) error {
        if state[id] == 1 {
            return fmt.Errorf("cycle detected at %q", id)
        }
        if state[id] == 2 {
            return nil
        }
        state[id] = 1
        for _, aid := range edges[id] {
            if err := visit(aid); err != nil {
                return err
            }
        }
        state[id] = 2
        return nil
    }
    for id := range edges {
        if err := visit(id); err != nil {
            return err
        }
    }
    return nil
}
```

**감지 시점:** `newEvalContext`에서 edges 구성 직후, calc 실행 전.

**API 변경:**

```go
// 변경 전
func newEvalContext(rules []RuleMeta, defeatEdges []defeatEdge, roleMap map[string]string) *evalContext

// 변경 후
func newEvalContext(rules []RuleMeta, defeatEdges []defeatEdge, roleMap map[string]string) (*evalContext, error)
```

```go
// 변경 전
func (b *GraphBuilder) Evaluate(claim any, ground any) []EvalResult
func (e *Engine) Evaluate(claim any, ground any) []EvalResult

// 변경 후
func (b *GraphBuilder) Evaluate(claim any, ground any) ([]EvalResult, error)
func (e *Engine) Evaluate(claim any, ground any) ([]EvalResult, error)
```

EvaluateTrace도 동일하게 `([]EvalResult, error)` 반환으로 변경한다.

**calc 단순화:** 순환이 사전에 거부되므로 maxDepth 상수와 depth 파라미터를 제거한다.

```go
// 변경 전
func (ctx *evalContext) calc(id string, claim, ground any, depth int) float64 {
    if depth >= maxDepth {
        return 0.0
    }
    // ...
    raw := (ctx.calc(aid, claim, ground, depth+1) + 1.0) / 2.0
}

// 변경 후
func (ctx *evalContext) calc(id string, claim, ground any) float64 {
    // depth 체크 불필요 — 순환 없는 그래프가 보장됨
    // ...
    raw := (ctx.calc(aid, claim, ground) + 1.0) / 2.0
}
```

### ParseYAML qualifier — pointer로 미지정 구분

```go
// 변경 전
type RuleDef struct {
    Name      string  `yaml:"name"`
    Role      string  `yaml:"role"`
    Qualifier float64 `yaml:"qualifier"`
}

// 변경 후
type RuleDef struct {
    Name      string   `yaml:"name"`
    Role      string   `yaml:"role"`
    Qualifier *float64 `yaml:"qualifier"`
}
```

ParseYAML에서 nil이면 1.0 기본값, 비nil이면 사용자 지정값(0.0 포함)을 사용한다.

### funcID nil 방어

```go
func funcID(fn func(any, any) (bool, any)) string {
    ptr := reflect.ValueOf(fn).Pointer()
    f := runtime.FuncForPC(ptr)
    if f == nil {
        return fmt.Sprintf("unknown_%d", ptr)
    }
    return f.Name()
}
```

### scanner 패키지 + 관련 dead code 제거

Phase 003에서 GraphBuilder + YAML 방식이 도입되면서 `//rule:` 어노테이션 기반
scanner 방식은 완전히 대체되었다. 다음 파일이 프로젝트 어디서도 호출되지 않는다:

**삭제 대상:**
- `internal/scanner/scan_dir.go` — Go 소스 파일 목록 수집
- `internal/scanner/extract_rules.go` — AST 파싱으로 `//rule:` 함수 추출
- `internal/scanner/extract_rule_lines.go` — `//rule:` 주석 라인 필터
- `internal/scanner/rule_decl.go` — RuleDecl 구조체
- `internal/cli/scan_and_parse.go` — scanner 호출 + parseDecls 조합
- `internal/cli/parse_decls.go` — RuleDecl → RuleMeta 변환

삭제 후 `internal/scanner/` 디렉토리 자체를 제거한다.

### internal 패키지 테스트 추가

각 패키지의 핵심 함수에 대해 정상/오류/경계 케이스를 검증한다.

**graphdef 테스트:**
- ParseYAML: 정상 파싱, qualifier 미지정 시 기본값 1.0, qualifier 0.0 명시 시 0.0 유지, 잘못된 YAML 오류, 빈 파일
- Validate: 정상 통과, graph name 누락, 존재하지 않는 from/to 참조

**codegen 테스트:**
- GenerateGraph: 정상 코드 생성, 생성 결과가 gofmt 유효, warrant/rebuttal/defeater 각 role 포함, defeat edge 포함

**graph 테스트:**
- ValidateGraph: 정상 통과, 존재하지 않는 defeat target 오류

## 범위

### 포함

1. **detect_cycle.go**: DFS 기반 순환 감지 함수 신규
2. **new_eval_context.go**: detectCycle 호출, 에러 반환 시그니처 변경
3. **eval_context.go**: maxDepth 상수 제거
4. **eval_context_calc.go**: depth 파라미터 제거
5. **eval_context_calc_trace.go**: depth 파라미터 제거
6. **Evaluate/EvaluateTrace 4곳**: 반환 타입 `([]EvalResult, error)` 변경, calc 호출에서 depth 제거
7. **func_id.go**: nil 방어 코드 추가
8. **graphdef/rule_def.go**: Qualifier를 *float64로 변경
9. **graphdef/parse_yaml.go**: nil → 1.0 기본값 로직 변경
10. **scanner 패키지 제거**: internal/scanner/ 디렉토리 전체
11. **scanner 관련 CLI dead code 제거**: scan_and_parse.go, parse_decls.go
12. **internal 테스트 3개 패키지**: graphdef, codegen, graph

### 제외

- h-Categoriser 수식 변경 없음

## 산출물

```
pkg/
  toulmin/
    detect_cycle.go               — 신규: DFS 순환 감지
    eval_context.go               — maxDepth 제거
    new_eval_context.go           — detectCycle 호출, error 반환
    eval_context_calc.go          — depth 파라미터 제거
    eval_context_calc_trace.go    — depth 파라미터 제거
    engine_evaluate.go            — ([]EvalResult, error) 반환
    engine_evaluate_trace.go      — ([]EvalResult, error) 반환
    graph_builder_evaluate.go     — ([]EvalResult, error) 반환
    graph_builder_evaluate_trace.go — ([]EvalResult, error) 반환
    func_id.go                    — nil 방어 추가
    engine_test.go                — 순환 시 에러 반환 테스트, API 변경 반영
    graph_builder_test.go         — 순환 시 에러 반환 테스트, 깊은 체인 테스트, API 변경 반영
internal/
  scanner/                        — 디렉토리 전체 삭제
    scan_dir.go                   — 삭제
    extract_rules.go              — 삭제
    extract_rule_lines.go         — 삭제
    rule_decl.go                  — 삭제
  cli/
    scan_and_parse.go             — 삭제 (scanner 의존)
    parse_decls.go                — 삭제 (scanner 의존)
    run_evaluate.go               — Evaluate API 변경 반영
  graphdef/
    rule_def.go                   — Qualifier *float64로 변경
    parse_yaml.go                 — nil qualifier 처리 로직 변경
    parse_yaml_test.go            — 신규
    validate_test.go              — 신규
  codegen/
    generate_graph_test.go        — 신규
  graph/
    validate_graph_test.go        — 신규
```

## 단계

### Step 1: detectCycle 함수 작성

- `detect_cycle.go` 신규: `detectCycle(edges map[string][]string) error`
- DFS 상태 기반 (unvisited/visiting/done) 순환 감지
- 순환 발견 시 관련 노드 이름을 포함한 에러 반환

### Step 2: newEvalContext에서 순환 감지

- `new_eval_context.go`: 반환 타입을 `(*evalContext, error)`로 변경
- edges 구성 직후 `detectCycle(ctx.edges)` 호출
- 에러 발생 시 nil, error 반환

### Step 3: calc/calcTrace에서 depth 제거

- `eval_context.go`: `maxDepth` 상수 제거
- `eval_context_calc.go`: `depth int` 파라미터 제거, `depth >= maxDepth` 체크 제거
- `eval_context_calc_trace.go`: 동일 변경

### Step 4: Evaluate/EvaluateTrace 반환 타입 변경

- `engine_evaluate.go`: `[]EvalResult` → `([]EvalResult, error)`
- `engine_evaluate_trace.go`: 동일 변경
- `graph_builder_evaluate.go`: 동일 변경
- `graph_builder_evaluate_trace.go`: 동일 변경
- `internal/cli/run_evaluate.go`: Evaluate 호출부 에러 처리 추가

### Step 5: funcID nil 방어

- `func_id.go`: `runtime.FuncForPC` nil 체크 추가
- fallback으로 pointer 기반 문자열 반환

### Step 6: scanner 패키지 + 관련 dead code 제거

- `internal/scanner/` 디렉토리 전체 삭제 (scan_dir.go, extract_rules.go, extract_rule_lines.go, rule_decl.go)
- `internal/cli/scan_and_parse.go` 삭제
- `internal/cli/parse_decls.go` 삭제
- 삭제 후 `go build ./...` 성공 확인 (다른 코드에서 import하지 않는지 검증)

### Step 7: ParseYAML qualifier 수정

- `graphdef/rule_def.go`: `Qualifier float64` → `Qualifier *float64`
- `graphdef/parse_yaml.go`: `qualifier == 0` 체크 → `qualifier == nil` 체크로 변경
- `internal/codegen/generate_graph.go`: `r.Qualifier` → `*r.Qualifier` 역참조 (nil이면 이미 ParseYAML에서 기본값 설정되므로 안전)

### Step 8: graphdef 테스트 작성

- `parse_yaml_test.go`:
  - 정상 YAML 파싱 (graph name, rules, defeats 확인)
  - qualifier 미지정 시 기본값 1.0
  - qualifier 0.0 명시 시 0.0 유지
  - 잘못된 YAML 형식 → 에러
  - 존재하지 않는 파일 → 에러
- `validate_test.go`:
  - 정상 GraphDef → nil
  - graph name 빈 문자열 → 에러
  - defeat에서 존재하지 않는 rule 참조 → 에러

### Step 9: codegen 테스트 작성

- `generate_graph_test.go`:
  - warrant만 있는 GraphDef → 유효한 Go 코드 생성
  - warrant + rebuttal + defeat → Defeat 메서드 포함 확인
  - 생성된 코드가 gofmt 유효성 통과
  - 빈 rules → 유효한 빈 graph

### Step 10: graph 테스트 작성

- `validate_graph_test.go`:
  - 정상 RuleMeta 목록 → nil
  - Defeats에 존재하지 않는 target → 에러
  - Defeats 빈 목록 → nil

### Step 11: 기존 테스트 수정 + 보강

- 기존 모든 테스트: Evaluate/EvaluateTrace 반환 타입 변경 반영 (`results, err := ...`)
- `engine_test.go`: 순환 그래프(A↔B) 테스트를 에러 반환 검증으로 변경
- `graph_builder_test.go`: 순환 시 에러 테스트 추가, depth > 100 비순환 체인 정상 평가 테스트 추가

### Step 12: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. 순환 그래프 구성 시 Evaluate가 에러를 반환한다
2. 비순환 그래프에서 기존 verdict와 동일한 결과를 반환한다
3. depth 101 이상의 비순환 체인에서 정상 verdict를 반환한다
4. calc/calcTrace에 depth 파라미터가 없다
5. `funcID`가 nil function pointer에 대해 panic하지 않는다
6. `internal/scanner/` 디렉토리가 존재하지 않는다
7. `internal/cli/scan_and_parse.go`, `internal/cli/parse_decls.go`가 존재하지 않는다
8. `go build ./...` 성공 (scanner 제거 후 컴파일 오류 없음)
9. ParseYAML에서 qualifier 미지정 시 1.0, 명시적 0.0 시 0.0이 된다
10. internal 패키지(graphdef, codegen, graph) 테스트가 모두 PASS한다
11. filefunc validate 0 위반

## 의존성

- 없음 (내부 개선)
