# Toulmin 기반 규칙 엔진 설계

## 핵심 구조

### Claim (struct — 입력)

판정 대상 명제. rule 평가의 입력.

```go
// PatternClaim — 단일 subject의 속성 검사
type PatternClaim struct {
    Subject  string   // 검증 대상 식별자 (e.g., 파일 경로)
    Type     string   // claim 유형 식별
}

// RelationClaim — 두 subject 간 교차 검사
type RelationClaim struct {
    Left     string   // 첫 번째 subject (e.g., OpenAPI endpoint)
    Right    string   // 두 번째 subject (e.g., Rego rule)
    Type     string
}
```

Claim 유형에 따라 적용 가능한 rule이 제한된다. Go 타입 시스템이 claim-rule 매칭을 컴파일타임에 검증한다.

### Ground (struct — 입력)

판정 근거 데이터. Ground Adapter가 공급.

```go
// FileGround — 코드 파일 검증용 근거
type FileGround struct {
    FuncCount    int
    TypeCount    int
    MaxDepth     int
    Control      string   // sequence | selection | iteration
    Dimension    int
    FileName     string
    Annotations  []string
}

// CrossGround — 교차 검증용 근거
type CrossGround struct {
    EndpointSecurity   []string
    PolicyRule         PolicyRule
    EndpointAnnotation []string
}
```

Ground는 0..N개. 완전 보류(§부분 평가) 시 0개.

### Rule (func — 연산)

모든 rule은 동일한 시그니처를 가진 boolean 함수다.

```go
type Rule func(claim any, ground any) bool
```

warrant, rebuttal, defeater의 구분은 타입이 아니라 **그래프에서의 관계(defeats 간선)**에 의해 결정된다.

| 구분 | 그래프에서의 역할 | 설명 |
|------|----------------|------|
| warrant | 공격받는 노드 | claim에 대한 판정을 주장 |
| rebuttal | 공격하는 노드 | warrant의 판정을 약화 |
| defeater | 공격하는 노드 | warrant의 판정을 차단 (자신은 판정 없음) |

### Strength

rule 자체의 속성이 아니라 **그래프에서 공격 가능 여부를 제어하는 제약**.

| strength | 의미 |
|----------|------|
| strict | 이 노드를 향하는 공격 간선 불허 |
| defeasible | 공격 간선 허용 |
| defeater | 이 노드에서 나가는 공격 간선만 존재, 들어오는 판정 없음 |

### Qualifier

`0.0–1.0` 실수. rule의 초기 가중치(w). h-Categoriser의 입력.

- 코드 규칙: 기본값 1.0 (결정적)
- 리얼월드 규칙: 1.0 미만 가능

### Verdict 값 범위

h-Categoriser 원래 출력은 `[0, 1]`이지만, `[-1, 1]`로 변환하여 중립점을 0으로 둔다.

```
변환: verdict = 2 * raw - 1   (raw = h-Categoriser 출력 [0,1])

-1.0  반박 확정 (fully rebutted)
 0.0  판정불가 (undecided)
+1.0  위반 확정 (fully confirmed)
```

### 판정 기준

```
verdict > 0   → 위반 (warrant가 우세)
verdict == 0  → 판정불가
verdict < 0   → 반박됨 (rebuttal이 우세)
```

Level(심각도)은 독립 속성이 아니라 verdict에서 파생된다.

## Toulmin 6요소 매핑

| Toulmin 요소 | 엔진에서의 역할 | 수량 |
|---|---|---|
| Claim | 판정 대상 명제. rule 평가의 입력 | 1 |
| Ground | 판정 근거 데이터. Ground Adapter가 공급 | 0..N |
| Warrant | rule (bool 함수). 그래프에서 공격받는 노드 | 1..N |
| Backing | warrant의 정당성 근거. 메타데이터 | 1..N |
| Qualifier | rule의 초기 가중치 w. h-Categoriser 입력 | 1 (per rule) |
| Rebuttal | rule (bool 함수). 그래프에서 공격하는 노드 | 0..N |

## 평가 흐름

```
1. Claim 추출
   "이 파일은 SRP를 준수한다"

2. 관련 rule 전부 평가 (boolean)
   모든 rule에 대해 func(claim, ground) → bool 실행
   true인 rule만 활성 노드로 수집

3. 활성 노드 + defeats 간선으로 서브그래프 구성
   (strict 노드는 들어오는 간선 불허)

4. h-Categoriser 연산
   raw(a) = w(a) / (1 + Σ raw(attackers))       [0, 1]
   verdict(a) = 2 * raw(a) - 1                   [-1, 1]
   순환 시 maxDepth에서 0.0(판정불가) 반환

5. 최종 판정
   verdict > 0  → 위반 (warrant 우세)
   verdict == 0 → 판정불가
   verdict < 0  → 반박됨 (rebuttal 우세)
```

## h-Categoriser

Amgoud & Ben-Naim (2017)의 가중 점진적 의미론.

### 수식

```
Acc(a) = w(a) / (1 + Σ Acc(attackers))
```

- `w(a)` = 초기 가중치 (qualifier)
- `attackers` = 이 노드를 공격하는 활성 노드들
- 공격이 없으면 `Acc = w` (그대로 통과)
- 공격이 있으면 수용도 감소
- 공격자도 공격받을 수 있으므로 반복 계산하여 수렴

### 특성

- **보상 원리(Compensation)**: 공격자가 방어(공격의 공격)되면 원래 노드의 수용도 복원. h-Categoriser가 유일하게 만족
- **수렴 보장**: 반복 계산이 항상 수렴
- **Backward compatible**: w=1.0이고 공격이 이진적이면 고전적 Dung과 동일 결과
- **구현**: 10줄 이내

### 구현

```go
const maxDepth = 100

// CalcAcceptability는 rule 그래프에서 노드의 최종 verdict를 계산한다.
// 반환값: [-1.0, 1.0]. 양수=위반, 0=판정불가, 음수=반박됨.
func CalcAcceptability(nodeID string, graph RuleGraph, depth int) float64 {
    if depth >= maxDepth {
        return 0.0 // 순환 — 판정불가
    }
    node := graph.Nodes[nodeID]
    attackerSum := 0.0
    for _, attackerID := range graph.Attackers(nodeID) {
        // attacker의 verdict는 [0,1] raw값으로 합산 (변환 전)
        raw := (CalcAcceptability(attackerID, graph, depth+1) + 1.0) / 2.0
        attackerSum += raw
    }
    raw := node.Qualifier / (1.0 + attackerSum) // [0, 1]
    return 2*raw - 1                             // [-1, 1]
}
```

## 예시

### 코드 규칙 (확신도 1.0, 공격 없음)

```
warrant "one-func-per-file" w=1.0
  attackers: 없음

raw = 1.0 / (1 + 0) = 1.0
verdict = 2*1.0 - 1 = +1.0 → 위반 확정
```

### 코드 규칙 (확신도 1.0, defeater 존재)

```
warrant A "one-func-per-file" w=1.0
  attacker: defeater D "test-file-exception" w=1.0

raw(D) = 1.0 / (1 + 0)   = 1.0  → verdict(D) = +1.0
raw(A) = 1.0 / (1 + 1.0) = 0.5  → verdict(A) =  0.0 → 판정불가 (무력화)
```

### 코드 규칙 (방어: 공격의 공격)

```
warrant A "one-func-per-file" w=1.0
  attacker: defeater B "test-file-exception" w=1.0
    attacker: defeater C "test-helper-not-exception" w=1.0

raw(C) = 1.0 / (1 + 0)   = 1.0   → verdict(C) = +1.0
raw(B) = 1.0 / (1 + 1.0) = 0.5   → verdict(B) =  0.0  ← C에 의해 무력화
raw(A) = 1.0 / (1 + 0.5) = 0.667 → verdict(A) = +0.33 ← B가 약화되어 A 부분 복원
```

### 리얼월드 규칙 (확신도 < 1.0)

```
warrant "페니실린 투여" w=0.85
  attacker: "알레르기 검사 양성" w=0.95

raw(attacker) = 0.95             → verdict(attacker) = +0.9
raw(warrant)  = 0.85 / (1+0.95)  → verdict(warrant)  = -0.13 → 반박됨
= 0.436
```

## 시점 선택적 아키텍처

### phase 속성

```
phase 미지정 → 모든 시점에 적용 (기본값)
phase: compiletime → 컴파일타임에만 적용
phase: runtime → 런타임에만 적용
```

### Ground Adapter

StaRVOOrS의 시점 통합 검증 + Hexagonal Architecture의 어댑터 패턴.
rule이 참조하는 추상 속성을 구체적 데이터 소스에서 추출하는 변환 계층.

```
adapter "compiletime-openapi" {
  endpoint.security → OpenAPIDoc.Paths[path].Operations[method].Security
}

adapter "runtime-request" {
  endpoint.security → request.header("Authorization").scheme
}
```

### 부분 평가

OPA의 Partial Evaluation 원리 적용.

| 유형 | 설명 |
|------|------|
| 완전 해소 | 컴파일타임 Ground만으로 판정 가능. 런타임 재평가 불필요 |
| 완전 보류 | 컴파일타임 Ground가 없어 런타임까지 보류 |
| 부분 해소 | 조건 일부 해소, 나머지 잔여 warrant로 런타임 이관 |

## 선행 연구 기반

| 요소 | 출처 |
|------|------|
| 6요소 구조 | Toulmin (1958) |
| strict/defeasible/defeater | Nute (1994) |
| h-Categoriser | Amgoud & Ben-Naim (2017) |
| Ground Adapter | StaRVOOrS (2015) + Hexagonal Architecture (2005) |
| 부분 평가 | OPA Partial Evaluation |
| Warrant의 인식론적 지위 | Gödel (1931), Kelsen (1934) |
| Toulmin SW 구현 선행 | Gabriel et al. (2020) |

## 참고문헌

### 논증 구조

- Toulmin, S. *The Uses of Argument*. Cambridge University Press, 1958.
  — 논증의 6요소 모델 (Claim, Ground, Warrant, Backing, Qualifier, Rebuttal). 본 엔진 설계의 전체 구조적 기초.

- Dung, P.M. "On the Acceptability of Arguments and its Fundamental Role in Nonmonotonic Reasoning, Logic Programming and n-person Games." *Artificial Intelligence*, 77(2):321–357, 1995.
  — 추상 논증 프레임워크. 논증을 노드, 공격을 간선으로 모델링. h-Categoriser의 이론적 토대.

- Modgil, S., and Prakken, H. "The ASPIC+ Framework for Structured Argumentation: A Tutorial." *Argument and Computation*, 5(1):31–62, 2014.
  — 구조화된 논증 프레임워크. strict/defeasible rule 구분, 전제·추론·결론에 대한 공격 형식화. Nute 분류의 Toulmin 위 통합에 참고.

### Defeasibility

- Nute, D. "Defeasible Reasoning." In *Handbook of Logic in Artificial Intelligence and Logic Programming*, Vol. 3, Oxford University Press, 1994.
  — 규칙의 세 강도 분류: strict (무력화 불가), defeasible (무력화 가능), defeater (차단만). 본 엔진의 strength 모델에 직접 적용.

- Prakken, H., and Sartor, G. "Argument-based Extended Logic Programming with Defeasible Priorities." *Journal of Applied Non-classical Logics*, 7(1):25–75, 1997.
  — 규칙 간 우선순위 자체도 무력화 가능하게 도출하는 변증법적 논증 체계.

### Verdict 연산 (h-Categoriser)

- Besnard, P., and Hunter, A. "A Logic-Based Theory of Deductive Arguments." *Artificial Intelligence*, 128(1-2):203–235, 2001.
  — h-Categoriser 원형. 비가중 버전: `Cat(a) = 1 / (1 + Σ Cat(attackers))`.

- Amgoud, L., and Ben-Naim, J. "Ranking-based Semantics for Argumentation Frameworks." *SUM 2013*, LNCS 8078, 2013.
  — 가중 h-Categoriser: `Hbs(a) = w(a) / (1 + Σ Hbs(attackers))`. 보상 원리(Compensation) 유일 만족 증명.

- Amgoud, L., and Ben-Naim, J. "Weighted Bipolar Argumentation Graphs: Axioms and Semantics." *IJCAI 2017*, pp. 5194–5198, 2017.
  — 가중 양극성 논증 그래프. 공격(attack)과 지지(support) 모두 [0,1] 가중치. h-Categoriser의 양극성 확장.

- Amgoud, L., and Ben-Naim, J. "Abstract Weighted Based Gradual Semantics in Argumentation Theory." arXiv:2401.11472, 2024.
  — 점진적 의미론의 통합 분석. h-Categoriser, Max-Based, Card-Based 의미론의 역문제(inverse problem) 분석.

### 퍼지 논증

- Tamani, N., and Croitoru, M. "Fuzzy Argumentation System for Decision Support." *SUM 2014*, LNCS 8720, 2014.
  — ASPIC 위에 퍼지 집합 이론 결합. 규칙과 전제에 [0,1] 중요도 부여, 인자에 강도 점수 부여.

- Da Costa Pereira, C., Tettamanzi, A., and Villata, S. "Fuzzy Labeling Semantics for Quantitative Argumentation." arXiv:2207.07339, 2022.
  — 퍼지 라벨링 의미론. 인자를 (수용도, 거부도, 미결정도) 삼중값으로 표현.

- Wu, Y., and Li, Z. "Gödel Fuzzy Argumentation Frameworks." *ECAI 2016*, FAIA Vol. 287, 2016.
  — Gödel t-norm 기반 퍼지 논증 프레임워크. Dung 의미론의 퍼지 확장.

### 베이지안/확률적 논증

- Hahn, U., and Oaksford, M. "A Normative Theory of Argument Strength." *Informal Logic*, 26(1):1–22, 2006.
  — 논증 강도를 P(Claim|Evidence)로 정량화. 우도비(LR)를 논증의 힘 지표로 사용.

- Prakken, H. "Probabilistic Strength of Arguments with Structure." *KR 2018*, pp. 158–167, 2018.
  — ASPIC+ 위에 확률 부여. 결합 확률 비일관성 문제 발견 — 내부 강도만으로 충돌 해소 불충분.

- Hunter, A., and Thimm, M. "Probabilistic Reasoning with Abstract Argumentation Frameworks." *Journal of Artificial Intelligence Research*, 59:565–611, 2017.
  — 확률적 논증의 두 접근: Constellations (인자 존재에 확률), Epistemic (수용에 확률).

### 시점 통합 / 어댑터

- Ahrendt, W., et al. "StaRVOOrS: A Tool for Combined Static and Runtime Verification of Java." *Runtime Verification*, LNCS 9333, 2015.
  — 하나의 명세로 정적 증명기와 런타임 모니터 양쪽에서 검증. Ground Adapter 개념의 원천.

- Ahrendt, W., et al. "A Specification Language for Static and Runtime Verification of Data and Control Properties." *FM 2015*, LNCS 9109, 2015.
  — ppDATE 명세 언어. 시점 통합 검증의 명세 수준 형식화.

- Cockburn, A. "Hexagonal Architecture (Ports and Adapters)." 2005.
  — 핵심 로직과 외부 시스템 사이에 포트와 어댑터를 두어 결합 제거. Ground Adapter의 아키텍처 패턴 원천.

### 인식론적 기반

- Gödel, K. "Über formal unentscheidbare Sätze der Principia Mathematica und verwandter Systeme I." *Monatshefte für Mathematik und Physik*, 38(1):173–198, 1931.
  — 불완전성 정리. 형식 체계는 자기 자신의 무모순성을 증명할 수 없다. Warrant를 fact로 취급하는 실용적 결단의 근거.

- Kelsen, H. *Reine Rechtslehre (Pure Theory of Law)*. Franz Deuticke, Wien, 1934.
  — 근본규범(Grundnorm). 법체계의 최종 효력 근거는 증명이 아니라 전제. Warrant의 공리적 지위와 대응.

### Toulmin의 소프트웨어 적용 선행 연구

- Gabriel, V.O., Panisson, A.R., Bordini, R.H., Adamatti, D.F., Billa, C.Z. "Reasoning in BDI agents using Toulmin's argumentation model." *Theoretical Computer Science*, 805:76–91, 2020.
  — Toulmin 모델을 BDI 에이전트에 구현한 유일한 선행 연구. 5요소 구현 (Backing 제외). qualify 함수로 SWW-SWR 계산. 본 엔진은 Backing 포함, Nute 3단계 적용, h-Categoriser로 연산 모델 교체.

- Panisson, A.R., McBurney, P., Bordini, R.H. "A Computational Model of Argumentation Schemes for Multi-agent Systems." *Argument & Computation*, 12(3):357–395, 2021.
  — Gabriel et al.의 후속. Toulmin에서 Walton의 argumentation schemes로 이동.

### 정책 엔진

- The Open Policy Agent Authors. "Rego Policy Language." https://www.openpolicyagent.org/docs/latest/policy-language/, 2024.
  — JSON 입력에 대한 논리적 규칙 평가. fact 모델의 대표적 구현. 본 엔진이 해결하려는 세 한계(경로 바인딩, defeasibility 부재, backing 미반영)의 원천.

### 안전성 사례

- Kelly, T.P. "Arguing Safety — A Systematic Approach to Managing Safety Cases." PhD Thesis, University of York, 1998. See also: ISO/IEC/IEEE 15026-2:2022.
  — Claim 기반 사고를 산업적으로 적용한 가장 성숙한 사례. GSN/CAE 프레임워크.

### 법률 규칙

- Palmirani, M., et al. "LegalRuleML: Design principles and foundations." *The 14th International Conference on AI and Law*, 2013.
  — 법률 규칙의 XML 표현. defeasible logic, 시간 조건, 위반-보상 관계 지원. fact 모델에서 defeasibility를 별도 메커니즘으로 추가해야 했던 사례.

- Kowalski, R., and Sergot, M. "The Use of Logical Models in Legal Problem Solving." *Ratio Juris*, 3(2):201–218, 1990.
  — 법률 조항의 논리 프로그램 표현. 법률 규정이 fact가 아니라 claim에 가까운 구조임을 보여주는 선행 사례.

## 결정 사항

### 프로젝트

- **라이브러리 이름**: `toulmin` — Toulmin 논증 모델을 규칙 엔진으로 구현한 Go 라이브러리
- **저장소**: `github.com/park-jun-woo/toulmin`
- **라이센스**: MIT
- **논문**: "Claim, Not Fact: Applying Toulmin's Argumentation Model to Software Rule Engines" (`toulmin-applying.md`)

### 구현

- **언어**: Go (OPA/Rego 엔진도 Go 구현)
- **Rule body**: Go 함수 `func(claim, ground) → bool`
- **메타데이터**: Go 어노테이션 (`//rule:` 접두사)
- **DSL**: 불필요. 엔진은 rule 내부를 볼 필요 없음 (h-Categoriser는 bool 결과와 그래프만 필요)
- **부분 평가**: PoC에서 제외 (미래의 fullend runtime 통합 시 재검토)
- **첫 적용 대상**: filefunc 22개 규칙

### Verdict

- **값 범위**: [-1.0, +1.0]. 중립점 0.0
- **변환**: `verdict = 2 * raw - 1` (h-Categoriser 원래 출력 [0,1] → [-1,1])
- **판정**: verdict > 0 → 위반, verdict == 0 → 판정불가, verdict < 0 → 반박됨
- **순환 공격**: maxDepth(100)에서 0.0(판정불가) 반환
- **Level(심각도)**: 독립 속성 아님. verdict에서 파생

### 논문의 위치

- **이론적 기여는 없음**: Toulmin(1958), Nute(1994), Amgoud(2017), StaRVOOrS(2015), Cockburn(2005) — 전부 기존 연구
- **기여는 융합**: 이 독립적 연구들이 연결된다는 발견 + 구현에 의한 실증
- **논문의 동기**: 학술적 공백 지적 (Toulmin이 60년간 규칙 엔진에 적용되지 않음)
- **구현의 동기**: fullend에 Rule DSL이 필요했으나 Rego의 구조적 한계 확인 → 독립 엔진 결정

### 소비자 관계

```
toulmin (라이브러리 — 도메인 무관)
  ├── filefunc (소비자 — 22개 코드 구조 규칙)
  └── fullend  (소비자 — 18개 crosscheck + 런타임 정책)
```

### Claim 타입과 Rule 매칭

- Claim 유형에 따라 적용 가능한 rule이 제한됨
- Go 타입 시스템이 claim-rule 매칭을 컴파일타임에 검증
- Claim 유형: 패턴 제약 (단일 subject), 존재 단언 (subject + 탐색 대상), 관계 단언 (subject 2개)

### Rego와의 관계

- Rego는 fullend crosscheck 18개를 표현 **가능** (불가능이 아님)
- 세 구조적 한계: (1) input 경로 바인딩 (2) defeasibility 비1급 (3) backing/qualifier 비시맨틱
- Rego에 Toulmin 요소를 끼워넣는 것은 fact 모델 위에 claim 시맨틱을 시뮬레이션하는 것 → 래퍼 복잡도 폭발
- 독립 엔진(2번 선택)이 근본 해결

## 미결정 사항

- [ ] Qualifier의 threshold 기본값 (0 기준 양수/음수로 단순화?)
- [ ] inline rebuttal과 독립 defeater 통합 여부
- [ ] Claim 모델 구체화 (RDF Triple 표현력 부족 — 관계 단언 미표현)
- [ ] phase 미지정 warrant의 양쪽 adapter 바인딩 제약
