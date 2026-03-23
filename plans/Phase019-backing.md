# Phase 019: Backing 인터페이스 강제 — 함수 필드 금지, YAML 직렬화 완전 호환

## 목표

Backing을 `any`에서 `Backing` 인터페이스로 변경한다. Backing struct에 func 필드를 금지하고, YAML 직렬화 가능한 데이터 타입만 허용한다. 함수 로직은 ground 준비 단계로 이동한다.

## 배경

### 현재 문제

1. **Backing이 `any`**: string, int, struct, nil, func — 아무거나 들어간다. 예측 불가.
2. **Backing struct에 func 필드**: `RoleBacking.RoleFunc`, `OwnerBacking.UserIDFunc`, `IPListBacking.Check` 등. YAML 직렬화 불가.
3. **YAML 동적 로딩 차단**: 프로덕션에서 YAML로 정책을 즉시 변경하려면 backing이 직렬화 가능해야 한다. func 필드가 있으면 불가.
4. **Backing의 의미 왜곡**: Toulmin 모델에서 backing은 "warrant의 판정 기준"(데이터)이지, 데이터 추출 로직(함수)이 아니다.

### 영향 범위

| 패키지 | Backing 타입 | func 필드 |
|---|---|---|
| policy | `RoleBacking` | `RoleFunc func(any) string` |
| policy | `OwnerBacking` | `UserIDFunc func(any) string`, `ResourceIDFunc func(any) string` |
| policy | `IPListBacking` | `Check func(string) bool` |
| price | `DiscountBacking` | 확인 필요 |

## 설계

### 1. Backing 인터페이스

```go
// pkg/toulmin/backing.go
type Backing interface {
    BackingName() string  // 타입 식별 — YAML type 필드, 에러 메시지, 로깅
    Validate() error      // 자기 검증 — 필수 필드 누락, 값 범위 등
}
```

BackingBase 불필요. 두 메서드 구현이면 끝. 프레임워크가 Backing struct를 제공하고, 사용자는 가져다 쓴다.

```
toulmin (코어)        → Backing interface 정의 (BackingName, Validate)
policy (프레임워크)    → RoleBacking, IPListBacking 구현
price (프레임워크)     → DiscountBacking 구현
써드파티               → 자기 Backing 구현
```

### 2. Backing 규칙

| 규칙 | 설명 |
|---|---|
| nil 허용 | backing 불필요한 룰 (IsAuthenticated 등) |
| Backing 인터페이스 필수 | nil이 아니면 BackingName() + Validate() 구현 |
| func 필드 금지 | LoadGraph에서 reflect 검사, 에러 반환 |
| primitive 필드만 | string, int, float64, bool + slice/map |
| Validate() | LoadGraph에서 호출. 도메인별 필수 필드/값 범위 검증 |

### 3. API 변경

```go
// 변경 전
func (g *Graph) Warrant(fn any, backing any, qualifier float64) *Rule
func (g *Graph) Rebuttal(fn any, backing any, qualifier float64) *Rule
func (g *Graph) Defeater(fn any, backing any, qualifier float64) *Rule

// 변경 후
func (g *Graph) Warrant(fn any, backing Backing, qualifier float64) *Rule
func (g *Graph) Rebuttal(fn any, backing Backing, qualifier float64) *Rule
func (g *Graph) Defeater(fn any, backing Backing, qualifier float64) *Rule
```

### 4. 프레임워크 Backing 리팩토링

**policy — RoleBacking**

```go
// 변경 전
type RoleBacking struct {
    Role     string
    RoleFunc func(any) string
}

// 변경 후
type RoleBacking struct {
    Role string `yaml:"role"`
}

func (b *RoleBacking) BackingName() string { return "RoleBacking" }
func (b *RoleBacking) Validate() error {
    if b.Role == "" {
        return fmt.Errorf("RoleBacking: role is required")
    }
    return nil
}
```

RoleFunc 제거. role 추출은 ground 준비 단계(ContextFunc)로 이동:

```go
// 변경 전
func buildCtx(r *http.Request) *RequestContext {
    return &RequestContext{User: getUserFromJWT(r), ClientIP: r.RemoteAddr}
}
// 프레임워크가 내부에서 RoleFunc(ctx.User) 호출

// 변경 후
func buildCtx(r *http.Request) *RequestContext {
    user := getUserFromJWT(r)
    return &RequestContext{
        User:     user,
        Role:     user.Role,       // 여기서 추출 완료
        ClientIP: r.RemoteAddr,
    }
}
// 프레임워크는 ctx.Role만 읽음
```

**policy — OwnerBacking**

```go
// 변경 전
type OwnerBacking struct {
    UserIDFunc     func(any) string
    ResourceIDFunc func(any) string
}

// 변경 후 — 제거. ground에서 처리.
// RequestContext에 UserID, ResourceOwnerID 필드 추가.
```

**policy — IPListBacking**

```go
// 변경 전
type IPListBacking struct {
    Purpose string
    Check   func(string) bool
}

// 변경 후
type IPListBacking struct {
    Purpose string   `yaml:"purpose"`
    List    []string `yaml:"list"`
}

func (b *IPListBacking) BackingName() string { return "IPListBacking" }
func (b *IPListBacking) Validate() error {
    if b.Purpose == "" {
        return fmt.Errorf("IPListBacking: purpose is required")
    }
    return nil
}
```

### 5. RequestContext 확장

```go
type RequestContext struct {
    User          any
    Role          string   // RoleFunc 대체
    UserID        string   // UserIDFunc 대체
    ResourceOwner string   // ResourceIDFunc 대체
    ClientIP      string
    IPBlocked     bool     // IPCheck 대체 (또는 List에서 판정)
    Headers       map[string]string
}
```

추출 로직이 ground 준비(ContextFunc)로 집중된다. 프레임워크는 데이터만 읽는다.

### 6. LoadGraph 검사

```go
// LoadGraph 내부에서 backing 검증 순서
if backing != nil {
    // 1. func 필드 금지 (reflect)
    if err := validateBackingFields(backing); err != nil {
        return nil, err
    }
    // 2. 도메인 검증 (Backing.Validate)
    if err := backing.Validate(); err != nil {
        return nil, err
    }
}
```

```go
func validateBackingFields(b Backing) error {
    t := reflect.TypeOf(b)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    for i := 0; i < t.NumField(); i++ {
        if t.Field(i).Type.Kind() == reflect.Func {
            return fmt.Errorf("backing %s has func field %q — func fields are not allowed",
                b.BackingName(), t.Field(i).Name)
        }
    }
    return nil
}
```

### 7. YAML backing 직렬화

Backing이 순수 데이터 struct이면 YAML 직렬화가 자연스럽다:

```yaml
graph: api:access
rules:
  - name: isInRole
    role: warrant
    backing:
      role: "admin"
  - name: isIPInList
    role: rebuttal
    backing:
      purpose: "blocklist"
      list: ["10.0.0.1", "192.168.1.100"]
  - name: isAuthenticated
    role: warrant
    # backing 없음 (nil)
```

## 구현 순서

| Step | 내용 |
|---|---|
| 1 | `pkg/toulmin/backing.go` — Backing 인터페이스 (BackingName, Validate) |
| 2 | `pkg/toulmin/validate_backing.go` — func 필드 reflect 검사 |
| 3 | Warrant/Rebuttal/Defeater 시그니처 변경: `any` → `Backing` |
| 4 | LoadGraph에 backing 검사 추가 (validateBackingFields + Validate) |
| 5 | policy — RoleBacking, OwnerBacking, IPListBacking 리팩토링 (func 제거, BackingName/Validate 구현) |
| 6 | policy — RequestContext 확장, Rule 함수 수정 (추출 로직을 ground으로 이동) |
| 7 | price, state, approve, feature, moderate — Backing 리팩토링 |
| 8 | 전체 테스트 수정 |
| 9 | 문서 업데이트 |

## 검증

```bash
go test ./...
filefunc validate
```

- 기존 테스트 전부 통과
- LoadGraph에 func 필드 backing 주입 시 에러 확인
- YAML 파싱 → LoadGraph → Evaluate 전체 흐름 확인

## Breaking Change

**Major breaking change.** Warrant/Rebuttal/Defeater의 backing 파라미터 타입이 `any` → `Backing`으로 변경. 기존 사용자 코드에서 string/int를 직접 넘기던 코드가 전부 컴파일 에러. `BackingName()` + `Validate()`를 구현한 struct로 감싸야 한다. 프레임워크(policy, price 등)가 제공하는 Backing struct를 사용하면 직접 구현할 필요 없다.
