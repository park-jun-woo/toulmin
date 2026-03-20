# toulmin

**if-else를 그만 중첩하라. 규칙을 선언하고, 관계를 선언하라.**

Go 룰 엔진. 규칙은 Go 함수다. 예외는 그래프 엣지다. DSL 없음. 사이드카 없음. 새 언어 없음.

규칙은 Go 함수다. 각 함수는 1-2 depth:

```go
func isAuthenticated(claim, ground, backing any) (bool, any) {
    return ground.(*Req).User != nil, nil
}
func isIPBlocked(claim, ground, backing any) (bool, any) {
    return blockedIPs[ground.(*Req).IP], nil
}
func isInternalIP(claim, ground, backing any) (bool, any) {
    return strings.HasPrefix(ground.(*Req).IP, "10."), nil
}
func isRateLimited(claim, ground, backing any) (bool, any) { /* ... */ }
func isPremiumUser(claim, ground, backing any) (bool, any) { /* ... */ }
func isIncidentMode(claim, ground, backing any) (bool, any) { /* ... */ }
```

요구사항은 진화한다. 양쪽이 어떻게 대응하는지 보라:

```go
// 월요일: "인증된 사용자만 접근, IP 차단 적용, 내부망은 차단 면제"
g := toulmin.NewGraph("api:access")
auth    := g.Warrant(isAuthenticated, nil, 1.0)
blocked := g.Rebuttal(isIPBlocked, nil, 1.0)
exempt  := g.Defeater(isInternalIP, nil, 1.0)
g.Defeat(blocked, auth)
g.Defeat(exempt, blocked)

// 화요일: "Rate limiting 추가"
limited := g.Rebuttal(isRateLimited, nil, 1.0)
g.Defeat(limited, auth)

// 수요일: "프리미엄 사용자는 Rate limit 면제"
premium := g.Defeater(isPremiumUser, nil, 1.0)
g.Defeat(premium, limited)

// 목요일: "장애 대응 중에는 프리미엄도 제한"
incident := g.Rebuttal(isIncidentMode, nil, 1.0)
g.Defeat(incident, premium)

results, _ := g.Evaluate(nil, req)
// results[0].Verdict > 0: 허용
```

매일 2줄 추가, 기존 코드 변경 없음. 같은 진화를 if-else로:

```go
// 월요일
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") {
            allow = true
        }
    } else {
        allow = true
    }
}

// 화요일: "Rate limiting 추가" — 어디에 끼워넣지?
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") {
            allow = true
        }
    } else if isRateLimited(ip) {
        allow = false
    } else {
        allow = true
    }
}

// 수요일: "프리미엄 사용자는 Rate limit 면제"
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") {
            allow = true
        }
    } else if isRateLimited(ip) {
        if isPremium(user) {       // 3중 중첩
            allow = true
        }
    } else {
        allow = true
    }
}

// 목요일: "장애 대응 중에는 프리미엄도 제한"
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") {
            allow = true
        }
    } else if isRateLimited(ip) {
        if isPremium(user) {
            if !incidentMode {     // 4중 중첩, 구조 파악 불가
                allow = true
            }
        }
    } else {
        allow = true
    }
}
```

toulmin: **요구사항당 2줄, 구조 불변.** if-else: **매번 전체 구조를 뜯어고친다.**

## 설치

```bash
go get github.com/park-jun-woo/toulmin/pkg/toulmin
```

## 핵심 개념

### 규칙은 Go 함수다

```go
func(claim any, ground any, backing any) (bool, any)
```

- `ground` = 요청마다 달라지는 판정 재료 (사용자, IP, 컨텍스트)
- `backing` = 그래프 선언 시 고정되는 판정 기준 (임계값, 역할명, 설정)
- 반환 = `(판정 결과, 증거)`. 증거는 도메인별 자유 타입.

```go
func isInRole(claim, ground, backing any) (bool, any) {
    user := ground.(*User)
    role := backing.(string)
    return user.Role == role, user.Role
}
```

### Defeats Graph

세 종류의 노드, 한 종류의 엣지:

| 노드 | 역할 |
|---|---|
| **Warrant** | 주장 — 공격받을 수 있다 |
| **Rebuttal** | 반박 — 주장을 공격한다 (자기 결론 있음) |
| **Defeater** | 예외 — 공격만 한다 (자기 결론 없음) |

```go
g := toulmin.NewGraph("voting")
auth     := g.Warrant(isAdult, nil, 1.0)
criminal := g.Rebuttal(hasCriminalRecord, nil, 1.0)
expunged := g.Defeater(isExpunged, nil, 1.0)
g.Defeat(criminal, auth)      // 전과가 성인을 공격
g.Defeat(expunged, criminal)  // 말소가 전과를 무력화
```

**Defeater가 차별점이다.** "예외의 예외"를 if-else로는 구조적으로 표현할 수 없다. defeats graph에서는 엣지 한 줄이다.

### Verdict

h-Categoriser가 연속값 `[-1, +1]`을 계산한다:

```
raw(a) = w(a) / (1 + Σ raw(attackers))
verdict = 2 × raw - 1
```

| Verdict | 의미 |
|---|---|
| +1.0 | 확정 |
| 0.0 | 미결 |
| -1.0 | 완전 반박 |

이진 판정이 아니다. **프레임워크가 해석**한다:

- 보안: `verdict ≤ 0` → 거부
- 모더레이션: `verdict ≤ 0` → 차단, `0 < v ≤ 0.3` → 검토, `> 0.3` → 허용
- 피처 플래그: `verdict > 0` → 활성

### Backing

같은 함수 + 다른 backing = 다른 규칙. 클로저 팩토리 없이 규칙을 재사용한다:

```go
g := toulmin.NewGraph("access")
admin  := g.Warrant(isInRole, "admin", 1.0)
editor := g.Warrant(isInRole, "editor", 0.8)
```

backing이 `nil`이면 규칙에 판정 기준이 필요 없다는 뜻이다.

## Trace

`EvaluateTrace`는 각 규칙의 판정 근거를 추적한다:

```go
results, _ := g.EvaluateTrace(claim, ground)
for _, t := range results[0].Trace {
    fmt.Printf("%s role=%s activated=%v evidence=%v\n",
        t.Name, t.Role, t.Activated, t.Evidence)
}
```

모더레이션 로그, 감사 추적, 디버깅 — 별도 로깅 없이 엔진이 제공한다.

## 동적 로딩

`LoadGraph`는 정의 + 함수 레지스트리로 라이브 그래프를 생성한다. 그래프 구조와 backing은 재배포 없이 변경 가능 — 함수는 컴파일된 채로 유지.

```go
// 기동 시 컴파일된 함수 등록
funcs := map[string]any{
    "isAuthenticated": isAuthenticated,
    "isIPBlocked":     isIPBlocked,
    "isRateLimited":   isRateLimited,
}

// YAML, DB, API에서 그래프 구조 로드 — 재컴파일 불필요
backings := map[string]any{"isIPBlocked": fetchBlocklistFromRedis()}
g, err := toulmin.LoadGraph(def, funcs, backings)
results, _ := g.Evaluate(nil, req)
```

컴파일 실행 속도 + 동적 규칙 업데이트. DSL 파서 없음, 인터프리터 없음, VM 없음 — 그래프 재배선만.

## 프레임워크 패키지

코어 위에 도메인별 프레임워크를 제공한다. 규칙 함수와 래퍼가 미리 구현되어 있다.

| 패키지 | 도메인 | 핵심 API |
|---|---|---|
| `pkg/policy` | 접근 제어 (인증, 인가, IP, Rate limit) | `Guard` (net/http 미들웨어) |
| `pkg/state` | 상태 전이 (FSM) | `Machine.Can`, `Mermaid()` |
| `pkg/approve` | 다단계 결재 | `Flow.Evaluate` |
| `pkg/price` | 할인 판정 (쿠폰, 멤버십) | `Pricer.Evaluate` |
| `pkg/feature` | 피처 플래그 (롤아웃, 토글) | `Flags.IsEnabled` |
| `pkg/moderate` | 콘텐츠 모더레이션 (혐오, 스팸) | `Moderator.Review` |

프레임워크 없이 코어만 써도 된다. 위 킬러 예제처럼 직접 규칙 함수를 작성하는 게 가장 유연하다.

## 왜 toulmin인가

### vs if-else

- 규칙 추가: `g.Defeat(new, existing)` 한 줄 vs 중첩 구조 전체 리팩터링
- 예외 처리: 엣지 선언 vs 조건문 안에 조건문
- 판정 근거: `EvaluateTrace` 내장 vs 별도 로깅 구축
- 테스트: 규칙 함수 단위 테스트 vs 조합 폭발

### vs OPA/Casbin/Cedar

| | toulmin | OPA | Casbin | Cedar |
|---|---|---|---|---|
| 규칙 언어 | **Go 함수** | Rego (DSL) | PERM 모델 (설정) | Cedar (DSL) |
| 예외 처리 | **defeats graph** | 규칙 우선순위 | 정책 우선순위 | forbid/permit |
| 예외의 예외 | **Defeater** | 없음 | 없음 | 없음 |
| 판정 | **연속값 [-1,1]** | allow/deny | allow/deny | allow/deny |
| 판정 근거 | **Trace 내장** | Decision log | 없음 | 없음 |
| 의존성 | Go 표준 라이브러리 | Go + Rego 런타임 | Go | Rust + FFI |
| 학습 곡선 | Go만 알면 됨 | Rego 학습 필요 | PERM 모델 학습 | Cedar 문법 학습 |

### 학술적 기반

| 구성 요소 | 출처 |
|---|---|
| 6요소 논증 구조 | Toulmin (1958) |
| strict/defeasible/defeater | Nute (1994) |
| h-Categoriser | Amgoud & Ben-Naim (2013, 2017) |

## CLI

```bash
toulmin graph voting.yaml                  # YAML → Go 코드 생성
toulmin graph voting.yaml --check          # 순환 검증만
toulmin graph voting.yaml --dry-run        # stdout 출력
toulmin graph voting.go                    # Go 파일 순환 분석
```

## 사용 사례

- **[filefunc](https://github.com/park-jun-woo/filefunc)** — LLM 네이티브 Go 코드 구조 도구. `validate` 명령이 toulmin defeats graph로 규칙 예외(F5, F6 등)를 처리한다.

## 라이선스

MIT
