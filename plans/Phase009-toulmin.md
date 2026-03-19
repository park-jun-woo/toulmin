# Phase 009: CLI 순환 감지 — YAML/Go 파일 그래프 검증

## 목표

`toulmin graph` CLI에 순환 감지를 통합한다.
YAML 파일은 코드 생성 전에, Go 파일은 AST 분석으로 defeat 관계를 추출하여 순환을 검증한다.
Phase 008의 `detectCycle` 함수를 재사용한다.

## 배경

Phase 008에서 `detectCycle`을 `pkg/toulmin`에 도입하여 Evaluate 시점에 순환을 거부한다.
그러나 사용자가 순환을 인지하는 시점이 코드 실행 시점(Evaluate 호출)으로 늦다.

### 현재 문제

1. **YAML → 코드 생성 시 순환 미검증**: `toulmin graph example.yaml`이 순환이 있는 YAML도 코드를 생성한다. 생성된 코드가 Evaluate 시점에야 에러를 반환한다
2. **Go 파일의 그래프 구조 검증 수단 없음**: 사용자가 직접 작성한 GraphBuilder 코드에서 `Defeat(A, B)` + `Defeat(B, A)` 같은 순환을 코드 리뷰 외에 감지할 방법이 없다
3. **피드백 지연**: 순환 오류를 컴파일 → 실행 → Evaluate 호출 후에야 알 수 있다. YAML 편집 또는 Go 파일 작성 시점에 감지하면 피드백이 빨라진다

## 핵심 설계

### CLI 흐름 확장

`toulmin graph` 명령이 입력 파일 확장자에 따라 동작을 분기한다:

```
toulmin graph example.yaml          → 순환 검증 → 코드 생성 (기존 + 순환 감지 추가)
toulmin graph example.go            → AST 분석 → 순환 검증 → 결과 출력 (신규)
toulmin graph example.yaml --check  → 순환 검증만 (코드 생성 안 함)
```

### YAML 경로: graphdef.Validate에 순환 감지 통합

기존 `graphdef.Validate`가 rule name 참조만 검증한다. 여기에 순환 감지를 추가한다.

```go
// Validate에서 기존 검증 후 순환 감지 추가
func Validate(def *GraphDef) error {
    // ... 기존: graph name, rule name 참조 검증
    // 신규: defeat edges로 순환 감지
    edges := buildEdgesFromDef(def)
    if err := toulmin.DetectCycle(edges); err != nil {
        return fmt.Errorf("graph validation failed: %w", err)
    }
    return nil
}
```

이를 위해 Phase 008의 `detectCycle`을 **공개 함수** `DetectCycle`으로 export한다.

```go
// pkg/toulmin/detect_cycle.go
// 변경: detectCycle → DetectCycle (공개)
func DetectCycle(edges map[string][]string) error
```

### Go 경로: AST 분석으로 defeat 관계 추출

Go 소스 파일을 파싱하여 GraphBuilder 체인에서 defeat 관계를 추출한다.

**추출 대상 패턴:**

```go
// 패턴 1: 메서드 체인
var g = toulmin.NewGraph("name").
    Warrant(FnA).
    Rebuttal(FnB).
    Defeat(FnB, FnA)

// 패턴 2: 분리된 호출
g := toulmin.NewGraph("name")
g.Warrant(FnA)
g.Defeat(FnB, FnA)
```

**분석 방법:**

1. `go/parser`로 AST 파싱
2. `Defeat` 메서드 호출을 찾아 인자 2개의 식별자(함수 이름)를 추출
3. `Warrant`/`Rebuttal`/`Defeater` 호출에서 등록된 함수 이름 수집
4. defeat edges + 등록된 이름으로 순환 감지

```go
// internal/analyzer/extract_defeats.go
type DefeatGraph struct {
    Name    string            // graph name
    Rules   []string          // 등록된 함수 이름
    Defeats map[string][]string // to → []from
}

// ExtractDefeats — Go 소스에서 GraphBuilder의 defeat 관계를 추출한다
func ExtractDefeats(path string) ([]DefeatGraph, error)
```

**한계와 범위:**
- 함수 이름은 AST상의 식별자(Ident)로 추출한다. 변수에 담긴 함수는 추적하지 않는다
- 한 파일에 여러 GraphBuilder가 있으면 각각 독립적으로 분석한다
- cross-file 분석은 하지 않는다 (단일 파일 범위)

### 출력 형식

```
$ toulmin graph example.yaml --check
✓ filefunc-check: no cycles detected (3 rules, 2 defeats)

$ toulmin graph broken.yaml --check
✗ filefunc-check: cycle detected at "RuleB"

$ toulmin graph example.go
✓ graph "filefunc-check": no cycles detected (3 rules, 2 defeats)

$ toulmin graph broken.go
✗ graph "filefunc-check": cycle detected at "RuleB"
```

## 범위

### 포함

1. **DetectCycle 공개화**: Phase 008의 `detectCycle` → `DetectCycle` (export)
2. **graphdef.Validate 순환 감지 통합**: EdgeDef → edges 변환 + DetectCycle 호출
3. **internal/analyzer 패키지 신규**: Go AST에서 defeat 관계 추출
4. **CLI 확장**: Go 파일 입력 지원, `--check` 플래그
5. **테스트**: analyzer 패키지 단위 테스트, CLI 통합 테스트

### 제외

- cross-file 분석 (단일 파일만)
- Engine.Register 기반 코드 분석 (GraphBuilder 패턴만)
- LSP/에디터 통합 — 별도 Phase에서 검토

## 산출물

```
pkg/
  toulmin/
    detect_cycle.go               — detectCycle → DetectCycle (export)
internal/
  analyzer/
    extract_defeats.go            — Go AST에서 defeat 관계 추출
    defeat_graph.go               — DefeatGraph 구조체
    extract_defeats_test.go       — 신규
  graphdef/
    validate.go                   — 순환 감지 추가
    build_edges_from_def.go       — EdgeDef → edges 변환 (신규)
  cli/
    new_graph_cmd.go              — --check 플래그 추가
    run_graph.go                  — 확장자 분기 + 순환 감지 통합
    run_graph_check.go            — Go 파일 분석 로직 (신규)
```

## 단계

### Step 1: DetectCycle 공개화

- `pkg/toulmin/detect_cycle.go`: `detectCycle` → `DetectCycle`
- `pkg/toulmin/new_eval_context.go`: 호출부 수정 (`DetectCycle`)

### Step 2: graphdef.Validate에 순환 감지 통합

- `internal/graphdef/build_edges_from_def.go` 신규: `EdgeDef` 슬라이스를 `map[string][]string`으로 변환
- `internal/graphdef/validate.go`: 기존 검증 후 `DetectCycle` 호출
- 기존 `toulmin graph example.yaml` 흐름에서 순환이 감지되면 코드 생성 전에 에러

### Step 3: analyzer 패키지 작성

- `internal/analyzer/defeat_graph.go`: DefeatGraph 구조체
- `internal/analyzer/extract_defeats.go`: Go AST 파싱
  - `ast.CallExpr`에서 `Defeat` 셀렉터 탐색
  - 인자 2개의 `ast.Ident.Name` 추출 → edge
  - `Warrant`/`Rebuttal`/`Defeater` 셀렉터에서 함수 이름 수집
  - `NewGraph` 호출에서 graph 이름 추출

### Step 4: CLI 확장

- `internal/cli/new_graph_cmd.go`: `--check` 플래그 추가
- `internal/cli/run_graph.go`: 입력 파일 확장자 분기 (`.yaml`/`.yml` → 기존 흐름, `.go` → 분석 흐름)
- `internal/cli/run_graph_check.go` 신규: Go 파일 분석 + DetectCycle + 결과 출력
- YAML + `--check`: Validate만 실행, 코드 생성 안 함

### Step 5: 테스트

- `internal/analyzer/extract_defeats_test.go`:
  - 메서드 체인 패턴에서 defeat 관계 추출
  - 분리된 호출 패턴에서 추출
  - 순환 있는 Go 코드 → DefeatGraph → DetectCycle 에러
  - Defeat 없는 코드 → 빈 edges
  - 한 파일에 여러 GraphBuilder → 각각 독립 분석
- `internal/graphdef/validate_test.go` 보강 (Phase 008에서 작성한 테스트에 순환 케이스 추가):
  - 순환 있는 GraphDef → 에러
  - 순환 없는 GraphDef → nil

### Step 6: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. `toulmin graph example.yaml`이 순환 있는 YAML에서 코드 생성 전에 에러를 반환한다
2. `toulmin graph example.yaml --check`가 순환 없으면 성공 메시지, 있으면 에러를 출력한다
3. `toulmin graph example.go`가 Go 파일의 defeat 관계를 추출하고 순환을 감지한다
4. `DetectCycle`이 pkg/toulmin에서 export되어 internal에서 import 가능하다
5. 메서드 체인 패턴과 분리 호출 패턴 모두에서 defeat 관계가 정확히 추출된다
6. analyzer 패키지 테스트가 모두 PASS한다
7. filefunc validate 0 위반

## 의존성

- Phase 008: `DetectCycle` 함수, Evaluate API 변경
