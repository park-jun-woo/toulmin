# Claim, Not Fact: 툴민 논증 모델의 소프트웨어 규칙 엔진 적용

## 초록

소프트웨어 규칙 엔진 — Rego/OPA, Drools, Semgrep, JSON Schema — 은 공통된 설계 전제를 공유한다: 검증 대상 데이터는 "fact(사실)"이다. 본 논문은 검증 대상이 fact가 아니라 **claim(주장)**이며, 툴민의 논증 모델(1958)이 60년 이상 간과되어 온 소프트웨어 규칙 엔진의 올바른 설계 기초임을 논증한다.

본 논문은 `toulmin`을 제시한다. 툴민의 6요소 — Claim, Ground, Warrant, Backing, Qualifier, Rebuttal — 을 구현한 Go 규칙 엔진으로, 규칙 강도에 Nute의 strict/defeasible/defeater 분류를, verdict 연산에 Amgoud의 h-Categoriser를 [-1, 1] 스케일로 적용한다. 규칙은 Go 함수(`func(claim, ground, backing) → (bool, evidence)`)로 작성되어 defeats 그래프로 조직되며, 엔진은 verdict를 연산하고 완전한 trace 설명 가능성을 제공한다. 별도의 DSL은 불필요하다.

filefunc의 22개 코드 구조 규칙을 툴민 warrant로 변환하여 3개 프로젝트에 대한 정량적 효과를 측정함으로써 설계를 검증한다. 나아가 fullend의 컴파일타임/런타임 정책 통합을 통해 시점 선택적 아키텍처를 시연한다.

---

## 1. 서론

### 1.1 문제: fact 전제

현대 소프트웨어 규칙 엔진은 보편적으로 입력을 "fact"로 취급한다. Drools는 Java 객체를 "fact"로 working memory에 적재한다. Rego는 `input` 문서를 암묵적으로 참인 것으로 전제하고 정책을 평가한다. JSON Schema는 문서 구조가 주어진 것으로 간주하고 검증한다.

이 용어는 인식론적으로 부정확하다. **Fact(사실)**는 이미 참으로 확립된 것이다 — 검증이 필요 없다. 그런데 규칙 엔진의 존재 이유는 입력 데이터가 규칙을 충족하는지 **검증하는 것**이다. 검증 대상을 "이미 참인 것"이라 부르는 것은 모순이다.

검증 대상은 **claim(주장)** — 참일 수도 거짓일 수도 있는 단언이며, 규칙에 의해 타당성이 판정되어야 한다. 이것은 단순한 용어 교정이 아니다. "fact" 전제가 기존 규칙 엔진의 세 가지 구조적 한계에 기여해왔다(§3.1).

### 1.2 미싱링크

툴민은 1958년에 논증의 구조를 6요소로 분석했다[1]: Claim, Ground, Warrant, Backing, Qualifier, Rebuttal. 이후 수십 년간 규칙 엔진들이 등장했다 — CLIPS(1985), Jess(1995), Drools(2001), Rego(2016) — 모두 fact 기반 설계를 공유한다. 툴민의 모델이 이미 올바른 구조를 제공하고 있었음에도, 소프트웨어 규칙 엔진의 설계 기초로 적용된 적이 없다.

이 공백이 존재하는 이유는 툴민의 저작이 철학과 수사학 분야에 출판되었기 때문이다. 규칙 엔진 개발자가 논증 이론을 참조할 이유가 없었다. 한편 각 규칙 엔진은 독립적으로 툴민 모델의 단편들을 재발명해왔다: LegalRuleML은 defeasibility 메커니즘(Rebuttal)을, OPA는 `# METADATA` 어노테이션(부분적 Backing)을, Drools는 salience(부분적 Qualifier)를 추가했다. 어느 것도 툴민을 참조하지 않았다.

### 1.3 기여

1. 툴민의 논증 모델이 소프트웨어 규칙 엔진 설계에 자연스럽게 매핑되며, "fact" 용어가 기존 시스템의 구조적 한계에 기여해왔음을 논증한다(§3).
2. 툴민의 6요소, Nute의 규칙 강도 분류[2], Amgoud의 h-Categoriser[3]를 결합한 Go 규칙 엔진 `toulmin`을 구현한다(§4).
3. filefunc의 22개 규칙에 대한 정량적 검증(§5)과, fullend의 시점 통합 사례(§6)를 통해 설계를 실증한다.

---

## 2. 배경 및 관련 연구

### 2.1 툴민의 논증 모델

툴민[1]은 논증을 6요소로 분석했다:

- **Claim(주장)**: 제기되는 단언.
- **Ground(근거/데이터)**: 주장을 뒷받침하는 증거.
- **Warrant(보증)**: 근거가 주장을 뒷받침한다고 판단하는 규칙이나 원리.
- **Backing(뒷받침)**: 보증이 왜 유효한지에 대한 정당성 근거.
- **Qualifier(한정)**: 확신의 정도 ("확실히", "아마도", "가능하게").
- **Rebuttal(반박)**: 주장이 성립하지 않는 예외 조건.

Gabriel et al.[4]은 이 중 5요소(Backing 제외)를 AgentSpeak/Jason 기반 BDI 에이전트 시스템에 구현했으며, qualify 함수로 SWW - SWR(warrant 가중치 합 - rebuttal 가중치 합)을 계산했다. 이것이 우리가 아는 한 툴민 모델의 유일한 소프트웨어 구현 선행 연구다.

### 2.2 무력화 가능 추론

Nute[2]는 규칙을 세 강도로 분류했다: **strict**(무력화 불가), **defeasible**(rebuttal에 의해 무력화 가능), **defeater**(자신의 결론 없이 다른 규칙의 결론만 차단). 이 분류는 ASPIC+[5]와 Dung의 추상 논증 프레임워크[6] 내에서 형식화되었다.

### 2.3 점진적 의미론

Amgoud와 Ben-Naim[3]은 논증 프레임워크를 위한 가중 h-Categoriser를 정의했다:

```
Acc(a) = w(a) / (1 + Σ Acc(attackers))
```

여기서 `w(a)`는 초기 가중치이고, `attackers`는 `a`를 공격하는 노드들이다. h-Categoriser는 **보상 원리(Compensation)**를 만족하는 유일한 점진적 의미론이다: 공격자가 자신도 공격받으면(방어), 원래 노드의 수용도가 상승한다. 수렴이 보장된다[3].

### 2.4 기존 규칙 시스템

| 시스템 | 입력 모델 | Defeasibility | 메타데이터 | 시점 |
|--------|----------|---------------|-----------|------|
| Rego/OPA | fact (`input`) | `default`/`else` 패턴 | `# METADATA` (비시맨틱) | Conftest: 다중 시점 |
| Drools | fact (working memory) | salience/agenda | 없음 | 런타임 전용 |
| Semgrep | AST 노드 | 없음 | 없음 | 컴파일타임 전용 |
| JSON Schema | JSON 문서 | 없음 | `description` | 컴파일타임 전용 |
| LegalRuleML | fact | XML 요소 (볼트온) | 없음 | N/A |

입력을 claim으로 취급하면서, 1급 defeasibility, 시맨틱 backing/qualifier, 시점 선택적 평가를 동시에 제공하는 기존 시스템은 존재하지 않는다.

### 2.5 시점 통합

StaRVOOrS[7]와 ppDATE[8]는 하나의 명세로 정적 검증(KeY)과 런타임 모니터링(LARVA) 양쪽에서 검증을 수행할 수 있음을 보였다. Cockburn의 Hexagonal Architecture[9]는 도메인 로직과 외부 시스템을 분리하는 어댑터 패턴을 형식화했다.

---

## 3. Fact에서 Claim으로

### 3.1 "Fact" 명명이 설계 한계에 기여하는 방식

규칙 엔진의 구조적 한계 원인은 복합적이다 — 기술 스택의 제약, 성능 최적화, 유스케이스 범위 등이 모두 기여한다. 그러나 "fact"라는 명명은 특정 설계 방향을 촉진하는 인지적 틀(cognitive frame)로 작용한다. 용어가 설계 철학을 형성하고, 설계 철학이 구현을 제약한다.

**기여 1: fact → 검증 인식 약화 → 정합성 불일치.**
Drools가 working memory 항목을 "fact"로 명명하면, 개발자는 인지적으로 이미 검증된 것으로 취급하게 된다. fact 간 교차 검증은 개발자가 명시적으로 규칙을 추가해야만 수행된다. Java의 객체 모델도 기여하지만, "fact" 명명은 검증 누락을 자연스럽게 느끼게 하는 인지적 틀을 제공한다.

**기여 2: fact → 데이터 구조 바인딩 → 규칙 중복.**
Rego가 `input`을 고정된 JSON 구조로 취급하면, 규칙은 구체적 경로(`input.user.role`)를 참조한다. Rego의 JSON 질의 문법도 기여 요인이지만, fact 기반 설계 철학이 규칙을 추상 속성이 아닌 구체적 데이터 형태에 바인딩하는 방향을 촉진한다. 동일한 논리적 규칙 — "관리자만 삭제할 수 있다" — 을 HTTP 요청, gRPC 메타데이터, CLI 인자에 각각 재작성해야 한다.

**기여 3: fact → 진리값 고정 → defeasibility의 사후 추가.**
Fact는 참이다. 참인 것은 뒤집히지 않는다. 이 전제 위에 설계된 시스템은 defeasibility를 위해 외부 예외 메커니즘을 필요로 한다. LegalRuleML이 `<Override>`, `<Defeater>`, `<Superiority>` XML 요소를 필요로 했던 것은 — 핵심 모델 바깥에 추가된 하위 시스템 전체다. 입력을 claim으로 재정의하면 defeasibility가 내재된 설계가 촉진되지만, Rebuttal(§4)이라는 구조적 장치는 여전히 필요하다.

**선례: JWT의 claims.**
JWT(RFC 7519)는 토큰 필드를 "facts"가 아니라 "claims"라 명명한다. `sub`, `exp`, `iss`는 토큰 발급자의 주장이다. 서명 검증, 만료 확인, issuer 대조가 기대되는 것은 정확히 데이터가 claim으로 프레이밍되었기 때문이다. 용어가 설계의 방향을 형성한 사례다.

### 3.2 툴민 매핑

| 툴민 요소 | 규칙 엔진에서의 역할 | 수량 |
|---|---|---|
| Claim | 판정 대상 명제. rule의 입력 | 1 |
| Ground | 판정 근거 데이터. Ground Adapter가 공급 | 0..N |
| Warrant | rule (bool 함수). defeats 그래프에서 공격받는 노드 | 1..N |
| Backing | warrant의 정당성 근거. 어노테이션(문서) + 런타임 값(판정 기준) | 1..N |
| Qualifier | 각 rule의 초기 가중치 w ∈ [0.0, 1.0] **(§3.3 참조)** | 1 (rule당) |
| Rebuttal | rule (bool 함수). defeats 그래프에서 공격하는 노드 | 0..N |

### 3.3 Qualifier의 재배치: Claim에서 Rule로

툴민의 원래 모델에서 Qualifier는 **Claim에 부착**된다. "**아마도(presumably)** 이 환자에게 페니실린을 투여해야 한다" — Qualifier는 주장의 확신도를 표현하는 양상 한정사(modal qualifier)다. Warrant, Ground, Rebuttal을 종합하여 화자가 Claim에 부여하는 확신의 정도다.

본 모델은 Qualifier를 Claim에서 **각 Rule(Warrant, Rebuttal, Defeater)**로 재배치한다. 이것은 논증 이론에서 규칙 엔진으로의 적용에서 불가피한 공학적 교정이다.

규칙 엔진에서 claim은 검증 대상일 뿐이다. "이 파일에 함수가 3개다" — 사실 확인이며, 확신도가 붙을 대상이 아니다. 판정의 질을 결정하는 것은 **규칙의 확신도**다:

- "파일당 함수 하나" — qualifier 1.0 (확실한 규칙)
- "권장 100줄 이하" — qualifier 0.7 (유연한 규칙)
- "알레르기 가능성" — qualifier 0.95 (강한 반박)

이 재배치는 Amgoud의 h-Categoriser[3]와 자연스럽게 결합한다. h-Categoriser의 `w(a)`가 각 노드(rule)의 초기 가중치이며, 최종 산출물인 verdict가 툴민 원래 모델에서 Qualifier가 담당하던 역할 — 판정의 확신도 — 을 대신한다.

| | 툴민 원래 모델 | 본 모델 |
|--|--------------|--------|
| Qualifier 위치 | Claim (출력 측) | 각 Rule (입력 측) |
| 역할 | 최종 판정의 확신도 | 각 rule의 초기 가중치 |
| 확신도 산출 | 화자가 종합적으로 판단 | h-Categoriser가 rule 가중치들로 계산 |
| 최종 확신도 | Qualifier 자체 | **verdict** (h-Categoriser 출력) |

### 3.4 Warrant의 인식론적 지위

엄밀히 말하면 warrant 자체도 claim이다 — "관리자만 삭제할 수 있다"는 조직의 합의이지 물리법칙이 아니다. 물리법칙조차 재현성이 극히 높은 claim이다. 그러나 이 인식을 무한히 적용하면 공리의 무한 후퇴에 빠진다. 괴델의 불완전성 정리[10]가 보여주듯, 충분히 강력한 형식 체계는 자기 자신의 무모순성을 증명할 수 없다. 켈젠의 근본규범(Grundnorm)[11]도 법체계에서 동일한 구조를 지적한다.

따라서 본 모델은 실용적 결단을 취한다: **Warrant는 시스템 내에서 fact로 취급한다.** 이는 모든 공리 체계가 공유하는 불가피한 전제이며, 모델의 적용 범위를 "warrant가 주어졌을 때 claim의 타당성을 판정하는 것"으로 한정한다.

---

## 4. 엔진 설계

### 4.1 증거를 반환하는 함수로서의 규칙

모든 규칙 — warrant, rebuttal, defeater — 은 동일한 시그니처를 공유한다:

```go
type Rule func(claim any, ground any, backing any) (bool, any)
```

규칙은 `(판정, 증거)`를 반환한다. bool은 판정이고, 두 번째 값은 도메인 특화 증거(예: 에러 상세, 위반 맥락)다. 세 번째 매개변수 `backing`은 규칙의 판정 기준을 런타임 값으로 전달한다.

**ground와 backing의 구분:**

| | Ground | Backing |
|--|--------|---------|
| 정의 | 판정 대상의 사실 | 규칙의 판정 기준 |
| 시점 | 요청마다 다름 | 선언 시점에 고정 |
| 예시 | 파일의 함수 개수, 사용자의 역할 | 허용 함수 개수 임계값, 인가 정책 |

이를 통해 엔진은 verdict를 연산할 뿐 아니라 **왜 그런 판정이 나왔는지를 설명**할 수 있다. backing이 1급 런타임 값이므로, 같은 함수를 다른 backing으로 등록하여 서로 다른 판정 기준의 규칙을 생성할 수 있다(§4.2 참조).

warrant, rebuttal, defeater의 구분은 함수 시그니처가 아니라 **defeats 그래프**에서의 위치에 의해 결정된다:

| 역할 | 그래프에서의 위치 |
|------|----------------|
| Warrant | 공격받을 수 있는 노드 |
| Rebuttal | warrant를 공격하는 노드 (자신의 결론을 가짐) |
| Defeater | 자신의 결론 없이 공격만 하는 노드 |

### 4.2 Graph API: 로직과 구조의 분리

엔진은 두 가지 API를 제공한다. Graph API(권장)는 규칙 로직(Go 함수)과 그래프 구조(defeats, qualifier, strength)를 분리한다:

```go
g := toulmin.NewGraph("file-structure")
w := g.Warrant(CheckOneFileOneFunc, nil, 1.0)
d := g.Defeater(TestFileException, nil, 1.0)
g.Defeat(d, w)
```

`Warrant(fn, backing Backing, qualifier)` — backing은 두 번째 인자이며 `Backing` 인터페이스(`BackingName() string`, `Validate() error`)를 구현해야 한다. 규칙에 고정된 판정 기준이 없을 때는 `nil`을 전달한다. `Warrant`, `Rebuttal`, `Defeater`는 `*Rule` 참조를 반환하며, 이 참조를 `Defeat(from, to)`에 전달하여 defeats 관계를 선언한다. backing이 있을 때:

```go
g := toulmin.NewGraph("line-limit")
strict := g.Warrant(CheckLineCount, &LineLimit{Max: 100}, 0.7)
relaxed := g.Warrant(CheckLineCount, &LineLimit{Max: 200}, 0.5)
g.Defeat(relaxed, strict)
```

같은 `CheckLineCount` 함수를 서로 다른 backing(`Max: 100`, `Max: 200`)으로 등록하면 별개의 노드가 된다. `*Rule` 참조가 각 노드를 고유하게 식별하므로, 같은 함수의 서로 다른 backing 노드 간 defeats 관계를 자연스럽게 선언할 수 있다.

함수가 식별자다 — 문자열 이름이 필요 없다. 같은 함수를 다른 그래프에서 다른 역할과 defeats 관계로 재사용할 수 있다:

```go
// 같은 IsAdult 함수, 다른 그래프, 다른 defeats
votingGraph := toulmin.NewGraph("voting")
vw := votingGraph.Warrant(IsAdult, nil, 1.0)
vr := votingGraph.Rebuttal(HasCriminalRecord, nil, 1.0)
votingGraph.Defeat(vr, vw)

contractGraph := toulmin.NewGraph("contract")
cw := contractGraph.Warrant(IsAdult, nil, 1.0)
cr := contractGraph.Rebuttal(IsBankrupt, nil, 1.0)
contractGraph.Defeat(cr, cw)
```

이 분리는 모든 구조적 메타데이터(역할, qualifier, strength, defeats)가 함수가 아니라 그래프에 속함을 의미한다. 함수는 순수한 판정 로직만 담는다.

### 4.3 그래프 제약으로서의 Strength

Nute의 분류[2]가 공격 간선을 제어한다:

| Strength | 의미 |
|----------|------|
| Strict | 들어오는 공격 간선 불허 |
| Defeasible | 들어오는 공격 간선 허용 |
| Defeater | 나가는 공격 간선만 존재, 자체 verdict 없음 |

### 4.4 Verdict 연산: h-Categoriser

Amgoud의 가중 h-Categoriser[3]를 [-1, 1]로 선형 변환하여 적용한다:

```
raw(a) = w(a) / (1 + Σ raw(attackers))     [0, 1]
verdict(a) = 2 × raw(a) - 1                [-1, 1]
```

| Verdict | 의미 |
|---------|------|
| +1.0 | 위반 확정 |
| 0.0 | 판정불가 |
| -1.0 | 반박 확정 |

판정: `verdict > 0` → 위반, `verdict == 0` → 판정불가, `verdict < 0` → 반박됨.

### 4.5 Evaluate와 EvaluateTrace

엔진은 두 가지 평가 모드를 제공하며, 각각 연산 방식을 선택할 수 있다:

```go
// Evaluate — verdict + 증거 (기본: 행렬곱)
results, err := g.Evaluate(claim, ground)
results, err  = g.Evaluate(claim, ground, toulmin.Recursive) // 재귀 h-Categoriser

// EvaluateTrace — verdict + 증거 + 완전한 trace (설명 가능성)
results, err := g.EvaluateTrace(claim, ground)
results, err  = g.EvaluateTrace(claim, ground, toulmin.Recursive)
```

`Matrix`(기본값)는 행렬곱으로 verdict를 연산한다. `Recursive`는 수학적으로 증명된 재귀 h-Categoriser 탐색을 사용한다.

EvalResult:

```go
type EvalResult struct {
    Name     string       `json:"name"`
    Verdict  float64      `json:"verdict"`
    Evidence any          `json:"evidence,omitempty"`
    Trace    []TraceEntry `json:"trace"`
}

type TraceEntry struct {
    Name      string  `json:"name"`
    Role      string  `json:"role"`       // "warrant", "rebuttal", "defeater"
    Activated bool    `json:"activated"`
    Qualifier float64 `json:"qualifier"`
    Backing   any     `json:"backing,omitempty"`
    Evidence  any     `json:"evidence,omitempty"`
}
```

EvaluateTrace는 완전한 설명 가능성을 제공한다: 어떤 규칙이 활성화되었고, 어떤 역할로, 어떤 qualifier로, 어떤 backing(판정 기준)으로, 어떤 증거를 생산했는지. backing이 trace에 포함됨으로써, 동일 함수가 서로 다른 backing으로 등록된 경우에도 각 노드의 판정 기준을 구분할 수 있다. 이것이 심볼릭 추론의 핵심 장점이다 — verdict의 도출 과정이 투명하고 감사 가능하다.

### 4.6 평가 흐름

```
0. 순환 감지: 그래프 구성 시점에 defeat edges에 대해 DFS 수행
   → 순환 발견 시 error 반환 (func 실행 전에 차단)
1. 각 warrant 노드에서 시작
2. warrant func(claim, ground, backing) 실행 → false? 건너뛰기
   backing은 그래프 선언 시 고정된 값이 전달됨
3. true이면 attackers (rebuttal/defeater) 재귀적 순회
4. 각 attacker: func(claim, ground, backing) 실행 → false? 기여 0 → true? 더 깊이 재귀
5. 각 노드에서 h-Categoriser: raw(a) = w(a) / (1 + Σ raw(attackers))
   verdict(a) = 2 * raw(a) - 1
6. func 결과 캐싱 — 각 (func, backing) 쌍은 평가당 최대 한 번 실행
7. warrant의 공격 체인에서 도달 가능한 규칙만 실행
```

### 4.7 YAML 그래프 정의와 코드 생성

그래프를 YAML로 정의하고 Go 코드로 컴파일할 수 있다:

```yaml
graph: file-structure
rules:
  - name: CheckOneFileOneFunc
    role: warrant
    qualifier: 1.0
  - name: TestFileException
    role: defeater
    qualifier: 1.0
defeats:
  - from: TestFileException
    to: CheckOneFileOneFunc
```

```bash
toulmin graph file-structure.yaml   # graph_gen.go 생성
```

### 4.8 Backing: 1급 런타임 값

backing은 **1급 런타임 값**이다. `Warrant(fn, backing Backing, qualifier)`의 두 번째 인자로 전달되며, `Backing` 인터페이스(`BackingName() string`, `Validate() error`)를 구현해야 한다. 엔진이 규칙 함수를 호출할 때 세 번째 매개변수로 주입된다. 이로써 같은 함수를 서로 다른 판정 기준으로 등록할 수 있다.

backing이 필요 없는 규칙은 `nil`을 전달하고, 함수 안에서 세 번째 매개변수를 무시한다. backing이 필요한 규칙은 이를 판정 기준으로 사용한다:

```go
func CheckOneFileOneFunc(claim any, ground any, backing any) (bool, any) {
    gf := ground.(*FileGround)
    if len(gf.Funcs) > 1 {
        return true, &Evidence{Got: len(gf.Funcs), Expected: 1}
    }
    return false, nil
}

func CheckLineCount(claim any, ground any, backing any) (bool, any) {
    gf := ground.(*FileGround)
    limit := backing.(*LineLimit)
    if gf.Lines > limit.Max {
        return true, &Evidence{Got: gf.Lines, Limit: limit.Max}
    }
    return false, nil
}
```

backing이 trace에 포함되므로, 동일 함수가 서로 다른 backing으로 등록된 경우에도 각 노드의 판정 기준을 구분할 수 있다.

규칙의 존재 이유를 문서화하고 싶을 때는 일반 Go 주석(`//`)을 사용한다. 별도의 어노테이션 문법은 불필요하다.

---

## 5. 사례 1: filefunc — 22개 규칙의 툴민 변환

### 5.1 배경

filefunc은 LLM 네이티브 Go 개발을 위한 코드 구조 컨벤션 및 CLI 도구다. 22개 규칙(F1-F6, Q1-Q3, A1-A16, C1-C4)을 Go 함수로 구현하여 컴파일타임에 AST를 대상으로 평가한다. 모든 규칙을 툴민 warrant로 변환했다.

### 5.2 Strength 분류

| Strength | 수 | 비율 | 예시 |
|----------|---|------|------|
| Strict | 15 | 68% | F1, F2, F3, F4, A1-A3, A6-A16 |
| Defeasible | 4 | 18% | Q1, Q2, Q3, C4 |
| Defeater | 3 | 14% | F5, F6, 테스트 파일 예외 |

대부분의 규칙이 strict다 — 코드 구조 컨벤션은 본질적으로 예외를 최소화한다. defeasible 규칙은 맥락에 따라 임계값이 달라지는 것들(깊이 제한 Q1-Q2, 크기 제한 Q3)이다.

### 5.3 정량적 결과

| 프로젝트 | 파일 수 (전→후) | 평균 LOC/파일 (전→후) | SRP 위반 해소 | depth 위반 해소 |
|---------|---------------|---------------------|-------------|--------------|
| filefunc | — (처음부터 준수) | 25.1 | 0 | 0 |
| fullend | 87→1,260 | 244→25.4 | 66→0 | 148→0 |
| whyso | 12→99 | 147.8→24.4 | 12→0 | 23→0 |

### 5.4 툴민 변환의 이점

1. **Backing의 1급 지위**: 판정 기준(임계값, 정책 등)이 런타임 backing 값으로 함수에 전달된다. 같은 함수를 다른 backing으로 등록하면 별개의 규칙이 된다.
2. **Rebuttal/Defeater의 구조화**: 예외 조건(F5 테스트 파일, F6 그룹 상수)이 명시적 `defeats` 관계로 선언된다. 코드 속 if 분기에 매몰되지 않는다.
3. **Strength 분류**: "이 규칙은 예외 없음"(strict) vs "이 규칙은 무력화 가능"(defeasible)이 선언적으로 표현된다.

---

## 6. 사례 2: fullend — 시점 선택적 아키텍처

### 6.1 배경

fullend는 9개 SSOT(OpenAPI, SSaC, Rego 등)의 정합성을 검증하고 코드를 생성하는 시스템이다. 18개 교차 검증 규칙이 Go 함수로 구현되어 있고, 런타임 정책은 별도의 Rego 파일이다. filefunc(§5)이 단일 시점(컴파일타임)에서 작동하는 반면, fullend는 시점 선택적 아키텍처의 필요성을 보여주는 사례다.

### 6.2 문제: 규칙 중복

동일한 논리적 규칙이 두 번 작성된다:

**컴파일타임** (Go crosscheck):
```go
func CheckRegoClaimsPresence(policies, openapi) []Error { ... }
```

**런타임** (Rego 정책):
```rego
allow if {
  input.action == "ExecuteWorkflow"
  input.claims.user_id > 0
}
```

### 6.3 툴민 해법

```go
// OWASP API Security Top 10 A2:2023
func AuthEndpointRequiresClaims(claim any, ground any, backing any) (bool, any) {
    g := ground.(*EndpointGround)
    if g.Security.Contains("bearerAuth") && !g.PolicyRule.References("claims") {
        return true, &Evidence{Endpoint: g.Path, Missing: "claims reference"}
    }
    return false, nil
}

// 공개 API endpoint는 설계상 인증이 불필요하다
func PublicEndpointException(claim any, ground any, backing any) (bool, any) {
    g := ground.(*EndpointGround)
    return g.Annotation.Contains("x-public"), nil
}

// 그래프 정의
g := toulmin.NewGraph("endpoint-auth")
w := g.Warrant(AuthEndpointRequiresClaims, nil, 1.0)
d := g.Defeater(PublicEndpointException, nil, 1.0)
g.Defeat(d, w)
```

**컴파일타임 Ground Adapter**: `EndpointGround.Security` ← OpenAPI 스펙 파싱, `EndpointGround.PolicyRule` ← Rego AST 파싱.

**런타임 Ground Adapter**: `EndpointGround.Security` ← 요청 라우트 미들웨어, `EndpointGround.PolicyRule` ← OPA 평가.

하나의 warrant 함수로 컴파일타임 교차 검증과 런타임 정책 적용을 모두 수행한다. 그래프 구조는 동일하며 Ground Adapter만 바뀐다. 함수 자체는 그래프 간에 재사용 가능하다 — 같은 `AuthEndpointRequiresClaims`가 다른 맥락에서 다른 defeats 관계로 참여할 수 있다.

---

## 7. 논의

### 7.1 본 논문이 주장하지 않는 것

본 논문은 이론적 독창성을 주장하지 않는다. 구조적 기초는 툴민(1958), 강도 분류는 Nute(1994), verdict 연산은 Amgoud(2017), 어댑터 패턴은 Cockburn(2005)과 StaRVOOrS(2015)다. 우리의 기여는 **독립적으로 발전한 이 조각들이 연결된다는 것을 발견한 것** — 툴민의 논증 모델이 소프트웨어 규칙 엔진의 누락된 설계 기초라는 것 — 이며, **구현을 통해 이를 실증한 것**이다.

### 7.2 한계

- Rule body가 Go 함수(엔진에 불투명)이므로, 부분 평가, 직렬화, 교차 언어 이식이 본 구현에서 지원되지 않는다.
- h-Categoriser의 [-1, 1] 변환은 선형 매핑이다. 복잡한 규칙 상호작용에는 더 정교한 verdict 집계가 필요할 수 있다.
- Claim 모델링에 추가 작업이 필요하다 — 두 주어 간의 관계 단언은 단순한 (subject, property, value) 삼중쌍의 표현력을 초과한다.
- 정량적 검증이 저자 자신의 프로젝트에 한정되어 있다. 유명 오픈소스 라이브러리에 대한 외부 검증을 계획 중이다.

### 7.3 기존 시스템과의 관계

본 엔진은 Rego, Drools, Semgrep을 대체하지 않는다. **상위 추상화**를 제공한다. Rego는 런타임 Ground Adapter 백엔드로 사용할 수 있다. Semgrep 규칙은 컴파일타임 Ground Adapter에서 재사용할 수 있다. 기존 시스템은 특정 시점의 Ground Adapter 구현체가 된다.

---

## 8. 결론

60년간 소프트웨어 규칙 엔진은 검증 대상을 "fact"로 취급하고 이 전제 위에 설계를 구축해왔다. 툴민의 논증 모델(1958)이 이미 올바른 구조 — Claim(Fact가 아닌), Ground, Warrant, Backing, Qualifier, Rebuttal — 을 제공했지만, 철학과 소프트웨어 공학 사이의 학문적 장벽이 그 적용을 가로막았다.

우리는 이 간극을 메우는 Go 규칙 엔진 `toulmin`을 구현했다. 규칙은 `func(claim, ground, backing) → (bool, evidence)` 시그니처의 Go 함수 — 판정과 그 근거다. backing은 규칙의 판정 기준으로서 선언 시점에 고정되어 런타임에 함수로 전달되며, ground는 요청마다 달라지는 판정 대상의 사실이다. 엔진은 Graph API를 통해 규칙을 defeats 그래프로 조직하고, 평가하고, Amgoud의 h-Categoriser를 통해 [-1, 1] 스케일의 verdict를 연산한다. Nute의 strict/defeasible/defeater 분류가 그래프의 공격 간선을 제어한다. EvaluateTrace는 완전한 설명 가능성을 제공한다 — 어떤 규칙이 활성화되었고, 어떤 역할로, 어떤 backing과 증거와 함께.

설계에 별도의 DSL이 필요 없다 — Go의 타입 시스템이 claim-rule 매칭을 처리하고, 엔진 코어는 수백 줄 규모다. 같은 규칙 함수를 다른 그래프에서 다른 defeats 관계로, 또는 다른 backing으로 재사용할 수 있으며, 규칙 로직과 그래프 구조가 분리된다. filefunc의 22개 규칙에 대한 3개 프로젝트 검증은 툴민 모델이 규칙 강도 분류, 구조화된 defeasibility, 명시적 backing을 자연스럽게 수용함을 보여준다. fullend의 시점 통합 사례는 Ground Adapter가 시점 선택적 규칙 재사용을 가능하게 함을 시연한다.

`toulmin` 라이브러리는 MIT 라이센스로 공개된다.

---

## 참고문헌

[1] Toulmin, S. *The Uses of Argument*. Cambridge University Press, 1958.

[2] Nute, D. "Defeasible Reasoning." In *Handbook of Logic in Artificial Intelligence and Logic Programming*, Vol. 3, Oxford University Press, 1994.

[3] Amgoud, L., and Ben-Naim, J. "Ranking-based Semantics for Argumentation Frameworks." *SUM 2013*, LNCS 8078, 2013.

[4] Gabriel, V.O., Panisson, A.R., Bordini, R.H., Adamatti, D.F., Billa, C.Z. "Reasoning in BDI agents using Toulmin's argumentation model." *Theoretical Computer Science*, 805:76–91, 2020.

[5] Modgil, S., and Prakken, H. "The ASPIC+ Framework for Structured Argumentation: A Tutorial." *Argument and Computation*, 5(1):31–62, 2014.

[6] Dung, P.M. "On the Acceptability of Arguments and its Fundamental Role in Nonmonotonic Reasoning, Logic Programming and n-person Games." *Artificial Intelligence*, 77(2):321–357, 1995.

[7] Ahrendt, W., et al. "StaRVOOrS: A Tool for Combined Static and Runtime Verification of Java." *Runtime Verification*, LNCS 9333, 2015.

[8] Ahrendt, W., et al. "A Specification Language for Static and Runtime Verification of Data and Control Properties." *FM 2015*, LNCS 9109, 2015.

[9] Cockburn, A. "Hexagonal Architecture (Ports and Adapters)." 2005.

[10] Gödel, K. "Über formal unentscheidbare Sätze der Principia Mathematica und verwandter Systeme I." *Monatshefte für Mathematik und Physik*, 38(1):173–198, 1931.

[11] Kelsen, H. *Reine Rechtslehre (Pure Theory of Law)*. Franz Deuticke, Wien, 1934.

