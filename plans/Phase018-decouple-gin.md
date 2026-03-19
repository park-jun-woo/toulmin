# Phase 018: Gin 결합 제거 — 프레임워크 독립 net/http 전환

## 목표

pkg/policy, pkg/feature, pkg/moderate에서 Gin 의존을 완전히 제거한다. 프레임워크 패키지는 `net/http` 기반 미들웨어만 제공한다. Gin 래퍼는 제공하지 않는다 — 프레임워크 독립이 원칙이다. 메인 모듈의 go.mod에서 gin-gonic/gin을 삭제한다.

## 배경

### 현재 문제

1. **gin이 indirect 의존인데 pkg에 침투**: go.mod에 `gin v1.12.0 // indirect`로 걸려 있으나 실제로 8개 파일이 gin을 직접 import한다. Gin 안 쓰는 프로젝트가 toulmin을 쓰면 gin + 간접 의존(quic-go, mongodb driver 등) 전부가 따라온다
2. **웹 프레임워크 종속**: Echo, Chi, net/http, Fiber 사용자는 Gin 어댑터를 쓸 수 없다
3. **테스트도 gin 의존**: policy 테스트가 `gin.TestMode`, `gin.New()`를 사용한다

### 영향 파일 (8개 + 테스트 1개)

| 패키지 | 파일 | gin 사용 |
|---|---|---|
| policy | `context_builder_func.go` | `func(*gin.Context) *RequestContext` |
| policy | `gin_guard.go` | `gin.HandlerFunc`, `c.AbortWithStatusJSON`, `c.Next` |
| policy | `gin_guard_debug.go` | `gin.HandlerFunc`, `c.Header`, `c.AbortWithStatusJSON` |
| policy | `gin_guard_test.go` | `gin.SetMode`, `gin.New()`, `gin.H` |
| feature | `feature_context_func.go` | `func(*gin.Context) *UserContext` |
| feature | `gin_inject.go` | `gin.HandlerFunc`, `c.Set`, `c.Next` |
| feature | `gin_require.go` | `gin.HandlerFunc`, `c.AbortWithStatus` |
| moderate | `content_extractor.go` | `func(*gin.Context) (*Content, *ContentContext)` |
| moderate | `gin_guard.go` | `gin.HandlerFunc`, `c.AbortWithStatusJSON` |

### 설계 원칙

- **프레임워크 독립** — net/http 표준 라이브러리만 사용한다. 특정 웹 프레임워크 래퍼를 제공하지 않는다
- **Gin/Echo/Chi 사용자는 README 예제로 안내** — 10줄 래퍼를 위해 별도 모듈을 유지하지 않는다
- **net/http 미들웨어 시그니처** — `func(http.Handler) http.Handler` 표준 패턴

## 핵심 설계

### 프레임워크-무관 미들웨어

#### policy

```go
// pkg/policy/context_func.go
type ContextFunc func(r *http.Request) *RequestContext
```

```go
// pkg/policy/guard.go
// Guard — verdict <= 0이면 403 거부하는 net/http 미들웨어
func Guard(g *toulmin.Graph, ctxFn ContextFunc) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            rc := ctxFn(r)
            results, err := g.Evaluate(rc, rc)
            if err != nil || results[0].Verdict <= 0 {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusForbidden)
                w.Write([]byte(`{"error":"forbidden"}`))
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// GuardDebug — trace 헤더 포함 버전
func GuardDebug(g *toulmin.Graph, ctxFn ContextFunc) func(http.Handler) http.Handler
```

#### feature

```go
// pkg/feature/context_func.go
type ContextFunc func(r *http.Request) *UserContext
```

```go
// pkg/feature/middleware.go
// Require — 피처 비활성이면 404 반환
func Require(flags *Flags, name string, ctxFn ContextFunc) func(http.Handler) http.Handler

// Inject — 활성 피처 목록을 context.Value에 주입
func Inject(flags *Flags, ctxFn ContextFunc) func(http.Handler) http.Handler
```

#### moderate

```go
// pkg/moderate/extract_func.go
type ExtractFunc func(r *http.Request) (*Content, *ContentContext)
```

```go
// pkg/moderate/guard.go
// Guard — Block→403, Flag→202, Allow→next
func Guard(m *Moderator, extractFn ExtractFunc) func(http.Handler) http.Handler
```

### README 사용 예제

각 프레임워크 패키지 README에 Gin/Echo/Chi 통합 예제를 포함한다. 별도 모듈이나 래퍼는 제공하지 않는다.

```go
// Gin — gin.WrapH로 net/http 미들웨어 사용
r := gin.Default()
guard := policy.Guard(g, func(r *http.Request) *policy.RequestContext {
    // http.Request에서 추출
})
r.Use(gin.WrapH(guard(http.DefaultServeMux)))

// 또는 직접 Evaluate 호출 (권장)
r.GET("/admin", func(c *gin.Context) {
    rc := buildCtx(c)
    results, _ := g.Evaluate(rc, rc)
    if results[0].Verdict <= 0 {
        c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
        return
    }
    c.Next()
})
```

```go
// Chi
r := chi.NewRouter()
r.Use(policy.Guard(g, buildCtx))

// Echo
e := echo.New()
e.Use(echo.WrapMiddleware(policy.Guard(g, buildCtx)))
```

## 범위

### 포함

1. pkg/policy — net/http Guard, GuardDebug, ContextFunc
2. pkg/feature — net/http Require, Inject, ContextFunc
3. pkg/moderate — net/http Guard, ExtractFunc
4. 기존 gin 파일 9개 삭제
5. go.mod에서 gin 의존성 제거
6. net/http 미들웨어 테스트
7. README 갱신 (Gin/Echo/Chi 통합 예제)

### 제외

- Gin 래퍼 모듈 (제공하지 않음 — 프레임워크 독립 원칙)
- 코어 엔진 변경 없음
- rule 함수 변경 없음
- 제네릭 전환 (Phase 019에서 별도 진행)

## 산출물

### 삭제 (9개)

| 파일 | 이유 |
|---|---|
| `pkg/policy/context_builder_func.go` | → `pkg/policy/context_func.go` |
| `pkg/policy/gin_guard.go` | → `pkg/policy/guard.go` (net/http) |
| `pkg/policy/gin_guard_debug.go` | → `pkg/policy/guard.go` (net/http) |
| `pkg/policy/gin_guard_test.go` | → `pkg/policy/guard_test.go` (net/http) |
| `pkg/feature/feature_context_func.go` | → `pkg/feature/context_func.go` |
| `pkg/feature/gin_inject.go` | → `pkg/feature/middleware.go` (net/http) |
| `pkg/feature/gin_require.go` | → `pkg/feature/middleware.go` (net/http) |
| `pkg/moderate/content_extractor.go` | → `pkg/moderate/extract_func.go` |
| `pkg/moderate/gin_guard.go` | → `pkg/moderate/guard.go` (net/http) |

### 신규 (8개)

| 파일 | 내용 |
|---|---|
| `pkg/policy/context_func.go` | `ContextFunc` 타입 (net/http) |
| `pkg/policy/guard.go` | `Guard`, `GuardDebug` (net/http 미들웨어) |
| `pkg/policy/guard_test.go` | net/http 미들웨어 테스트 |
| `pkg/feature/context_func.go` | `ContextFunc` 타입 (net/http) |
| `pkg/feature/middleware.go` | `Require`, `Inject` (net/http 미들웨어) |
| `pkg/feature/middleware_test.go` | net/http 미들웨어 테스트 |
| `pkg/moderate/extract_func.go` | `ExtractFunc` 타입 (net/http) |
| `pkg/moderate/guard.go` | `Guard` (net/http 미들웨어) |

## 단계

### Step 1: net/http 미들웨어 작성

1. `pkg/policy/context_func.go` — `ContextFunc` 타입
2. `pkg/policy/guard.go` — `Guard`, `GuardDebug` (net/http)
3. `pkg/feature/context_func.go` — `ContextFunc` 타입
4. `pkg/feature/middleware.go` — `Require`, `Inject` (net/http)
5. `pkg/moderate/extract_func.go` — `ExtractFunc` 타입
6. `pkg/moderate/guard.go` — `Guard` (net/http)

### Step 2: 테스트 작성

7. `pkg/policy/guard_test.go` — httptest 기반 테스트
8. `pkg/feature/middleware_test.go` — httptest 기반 테스트

### Step 3: 기존 gin 파일 삭제

9. 9개 파일 삭제
10. go.mod에서 gin 관련 의존성 정리 (`go mod tidy`)

### Step 4: 검증

11. `go build ./...` — 메인 모듈에 gin import 없음
12. `go test ./...` — 전체 PASS
13. `go mod graph | grep gin` — gin 의존 없음 확인

### Step 5: 문서 갱신

14. pkg/policy/README.md — Guard 사용법을 net/http 기반으로 변경 + Gin/Echo/Chi 통합 예제
15. pkg/feature/README.md — Require/Inject 사용법 변경 + 통합 예제
16. pkg/moderate/README.md — Guard 사용법 변경 + 통합 예제

## 검증 기준

1. 메인 모듈 `go.mod`에 gin-gonic/gin이 없다
2. `go mod graph | grep gin` 결과 비어 있다
3. `pkg/**/*.go`에 `gin-gonic/gin` import가 없다
4. `go build ./...` 통과
5. `go test ./...` 전체 PASS
6. net/http 미들웨어가 기존 gin 미들웨어와 동일한 판정 동작 (verdict 기준, HTTP 상태 코드)
7. README에 Gin/Echo/Chi 통합 예제 포함

## 의존성

- Phase 001-017: 전체 기존 구현

## 예상 규모

- 삭제 파일: 9개
- 신규 파일: 8개
- 수정 파일: 3개 (README)
- 예상 난이도: 낮음 (기계적 이전, 로직 변경 없음)
