# Claim, Not Fact: 툴민 논증 모델의 소프트웨어 규칙 엔진 적용

© 2026 PARK JUN WOO(mail@parkjunwoo.com). 이 논문은 크리에이티브 커먼즈 저작자표시 4.0 국제 라이선스(CC BY 4.0, https://creativecommons.org/licenses/by/4.0/deed.ko)를 따른다. 이 저장소의 소프트웨어는 따로 MIT 라이선스를 따른다.

## 초록

소프트웨어 규칙 엔진(Rego/OPA, Drools, Semgrep, JSON Schema 등)은 대개 검증 대상 데이터를 "fact(사실)"로 다룬다. 본 논문은 검증 대상이 fact가 아니라 **claim(주장)**이며, 툴민의 논증 모델(1958)이 규칙 엔진 설계에서 거의 쓰이지 않았지만 그 설계 기초로 알맞다고 논증한다.

이 관점을 구현한 규칙 엔진 `toulmin`을 제시한다(이 판은 v0.9.0 기준). 툴민의 6요소(Claim, Ground, Warrant, Backing, Qualifier, Rebuttal)를 엔진 구성 요소에 대응시키고, 규칙 강도에는 Nute의 strict/defeasible/defeater 분류를, 판정값 계산에는 가중 h-Categoriser를 [-1, 1]로 옮겨 쓴다. 규칙은 판정 대상 맥락과 판정 기준(spec)을 받아 규칙값 v ∈ [-1, 1]과 증거를 돌려주는 함수이고, 규칙들은 공격 간선으로 이어진 그래프를 이룬다. v0.9.0의 판정 계산은 h-Categoriser를 네 방향으로 넓힌다: 규칙값으로 정해지는 적용도, 분모에 들어가는 반박과 적용도를 곱으로 깎는 예외의 구분, 공격 간선의 세기 F ∈ [0, 1], 그리고 판정값을 바꾸지 않고 하류 집계에서 뺄 노드를 표시하는 크리스프 배제(Supersedes)다. 같은 계산이 Go 참조 구현과 TypeScript·Python·C 포트에 있으며, 마크다운 정책 언어 TANGEUL이 코어 위의 선택적 선언 표면을 제공한다.

설계를 두 사례로 살핀다. filefunc의 코드 구조 규칙 22개를 툴민 warrant로 옮기고 그 규칙 집합을 적용한 세 프로젝트의 수치를 보고하며, fullend의 컴파일타임·런타임 정책을 한 규칙 함수로 묶는 시점 선택적 설계를 보인다.

---

## 1. 서론

### 1.1 문제: fact 전제

현대 소프트웨어 규칙 엔진은 흔히 입력을 "fact"로 다룬다. Drools는 Java 객체를 fact로 working memory에 넣는다. Rego는 `input` 문서를 주어진 것으로 놓고 정책을 평가한다. JSON Schema는 문서(instance)를 주어진 것으로 놓고 스키마에 맞는지 검사한다.

이 용어는 인식론적으로 어긋난다. **Fact(사실)**는 이미 참으로 확립된 것이어서 검증할 필요가 없다. 그런데 규칙 엔진이 존재하는 이유는 입력 데이터가 규칙을 충족하는지 **검증하는 것**이다. 검증 대상을 "이미 참인 것"이라 부르는 것은 모순이다.

검증 대상은 **claim(주장)**, 곧 참일 수도 거짓일 수도 있고 규칙에 비추어 타당성을 판정받아야 하는 단언이다. 본 논문은 이것이 용어 교정에 그치지 않는다고 본다. §3.1에서 "fact" 전제가 기존 규칙 엔진의 세 가지 구조적 한계와 맞물려 있다는 가설을 제시한다.

### 1.2 놓친 연결

툴민은 1958년에 논증의 구조를 6요소로 분석했다[1]: Claim, Ground, Warrant, Backing, Qualifier, Rebuttal. 그 뒤 여러 규칙 엔진이 나왔다(CLIPS 1985, Jess 1995, Drools 2001, OPA/Rego 2016). 이들은 fact 기반 설계를 공유하며, 우리가 조사한 범위에서는 툴민 모델을 설계 기초로 삼은 규칙 엔진을 찾지 못했다.

한 가지 이유로, 툴민의 저작이 철학과 수사학 분야에서 읽혀 왔고 규칙 엔진 개발의 참고 문헌이 되지 않았다는 점을 들 수 있다. 한편 각 규칙 엔진에는 툴민 모델의 일부 요소와 대응시켜 볼 수 있는 장치가 따로 들어갔다: LegalRuleML의 defeasibility 요소(Rebuttal에 대응), OPA의 `# METADATA` 주석(Backing의 일부에 대응), Drools의 salience(Qualifier의 일부에 대응)다.

### 1.3 기여

1. **개념**: 검증 대상을 claim으로 보는 관점을 제시하고, 툴민의 6요소를 규칙 엔진 구성 요소에 대응시킨다(§3).
2. **형식 정의**: 가중 h-Categoriser[3]를 규칙값·적용도, 반박과 예외의 구분, 공격 세기, 보고용 배제로 넓힌 판정 계산을 정의한다(§4).
3. **구현**: 이 계산을 구현한 `toulmin`의 Go 참조 구현과 TypeScript·Python·C 포트, 선언 표면 TANGEUL을 기술한다(§5).
4. **사례**: filefunc 규칙 22개의 변환(§6)과 fullend의 시점 통합 설계(§7)로 설계를 살핀다.

---

## 2. 배경 및 관련 연구

### 2.1 툴민의 논증 모델

툴민[1]은 논증을 6요소로 분석했다.

- **Claim(주장)**: 제기되는 단언.
- **Ground(근거, Data)**: 주장을 뒷받침하려고 내놓는 자료.
- **Warrant(보증)**: 근거에서 주장으로 넘어가도 된다고 허락하는 규칙이나 원리.
- **Backing(뒷받침)**: 보증이 왜 유효한지에 대한 정당화.
- **Qualifier(한정)**: 확신의 정도("확실히", "아마도", "어쩌면").
- **Rebuttal(반박)**: 주장이 성립하지 않는 예외 조건.

툴민 도식을 형식화해 논증 평가에 쓰려는 연구로 Verheij[19]가 있다. Gabriel et al.[4]은 이 중 5요소(Backing 제외)를 AgentSpeak/Jason 기반 BDI 에이전트 시스템에 구현했으며, qualify 함수로 SWW − SWR(warrant 가중치 합 − rebuttal 가중치 합)을 계산했다. 우리가 조사한 범위에서 툴민 모델을 범용 소프트웨어 규칙 엔진의 설계 기초로 쓴 선행 연구는 찾지 못했다.

### 2.2 무력화 가능 추론

Nute[2]는 규칙을 세 강도로 나눴다: **strict**(무력화되지 않음), **defeasible**(반대 증거에 의해 무력화될 수 있음), **defeater**(자기 결론 없이 다른 규칙의 결론만 막음). Nute 계열의 defeasible logic은 규칙 사이의 우선관계(superiority relation)도 둔다[18].

Pollock[12]은 무력화 요인(defeater)을 둘로 나눴다. **반박(rebutting) defeater**는 결론의 부정을 지지하고, **잠식(undercutting) defeater**는 근거와 결론 사이의 연결이 이 상황에서는 성립하지 않는다고 공격한다. 구조적 논증 틀 ASPIC+[5]는 strict 규칙과 defeasible 규칙을 구분하고 반박·잠식 공격을 함께 다루며, 그 의미는 Dung의 추상 논증 틀[6] 위에서 정해진다.

### 2.3 추상 논증과 점진 의미론

Dung의 추상 논증 틀[6]은 논증을 노드로, 공격을 간선으로 보고 받아들일 수 있는 논증 집합(extension)을 정한다. 라벨링 의미론[14]은 각 논증에 IN·OUT·UNDEC 라벨을 붙인다. 이런 의미론에서 판정은 크리스프하다.

점진 의미론(gradual semantics)은 논증마다 연속값의 수용도를 매긴다. h-Categoriser는 Besnard와 Hunter[15]가 도입했고, Amgoud et al.[3]은 노드마다 초기 가중치를 둔 가중 h-Categoriser를 정의했다.

```
Acc(a) = w(a) / (1 + Σ_{b ∈ Att(a)} Acc(b))
```

여기서 `w(a) ∈ [0, 1]`은 초기 가중치, `Att(a)`는 `a`를 공격하는 노드의 집합이다. 식에서 바로 보이듯 공격자의 수용도가 낮아지면(예컨대 공격자 자신이 공격받으면) `a`의 수용도는 올라간다. 순환이 있는 일반 그래프에서의 정의와 성질은 [3]을 보라. 순위 기반 의미론(ranking-based semantics)의 원리들은 Amgoud와 Ben-Naim[16]이 정리했다. 공격 간선에 가중치를 다는 착상은 Dunne et al.의 가중 논증 체계[17]에 있다. 그곳에서 가중치는 불일치 예산(inconsistency budget) 안에서 어떤 공격을 무시할 수 있는지를 정하는 데 쓰인다.

### 2.4 기존 규칙 시스템

| 시스템 | 입력 모델 | Defeasibility | 메타데이터 | 시점 |
|--------|----------|---------------|-----------|------|
| Rego/OPA | fact (`input`) | `default`/`else` 패턴 | `# METADATA` (판정 계산에 쓰이지 않음) | 런타임 정책 평가, Conftest로 설정 파일 검사 |
| Drools | fact (working memory) | salience/agenda (발화 순서) | 규칙 속성 | 런타임 |
| Semgrep | 소스 코드 패턴 | 없음 | 규칙 `metadata` (판정 계산에 쓰이지 않음) | 정적 분석 |
| JSON Schema | instance (JSON 문서) | 없음 | `description` 등 주석 키워드 | 검증기를 부르는 시점(정적·런타임) |
| LegalRuleML | 법 규범 표현(XML) | 핵심 모델 위에 덧붙인 요소 | 법적 메타데이터 | 표현 언어(자체 실행 시점 없음) |

우리가 살펴본 시스템 가운데 입력을 claim으로 다루면서 1급 defeasibility, 판정에 참여하는 spec·qualifier, 시점 선택적 평가를 함께 갖춘 것은 없었다.

### 2.5 시점 통합

StaRVOOrS[7]와 ppDATE[8]는 하나의 명세로 정적 검증(KeY)과 런타임 모니터링(LARVA)을 함께 할 수 있음을 보였다. Cockburn의 Hexagonal Architecture[9]는 도메인 로직과 외부 시스템을 분리하는 어댑터 패턴을 정리했다.

---

## 3. Fact에서 Claim으로

### 3.1 "fact" 명명과 설계 한계: 세 가설

규칙 엔진의 구조적 한계에는 여러 원인이 있다. 기술 스택의 제약, 성능 최적화, 유스케이스 범위가 모두 영향을 준다. 본 절은 여기에 더해 "fact"라는 이름이 특정 설계 방향을 자연스럽게 느끼게 하는 인지적 틀(cognitive frame)로 작용한다는 가설을 세운다. 아래 세 경로는 개념적 논증이며, 경험적으로 검증하지는 않았다(§8.2).

**경로 1: fact → 검증 인식 약화 → 정합성 불일치.**
Drools가 working memory 항목을 "fact"로 부르면, 개발자는 그것을 이미 검증된 것으로 여기기 쉽다. fact 사이의 교차 검증은 개발자가 규칙을 명시적으로 추가해야만 일어난다. Java의 객체 모델도 영향을 주지만, "fact"라는 이름은 검증 누락을 자연스럽게 느끼게 하는 틀을 제공한다.

**경로 2: fact → 데이터 구조 바인딩 → 규칙 중복.**
Rego는 `input`을 JSON 문서로 받고, 규칙은 구체적 경로(`input.user.role`)를 참조한다. Rego의 JSON 질의 문법도 한 요인이지만, fact 기반 설계는 규칙을 추상 속성이 아닌 구체적 데이터 모양에 묶는 방향을 부추긴다. 그 결과 같은 논리 규칙("관리자만 삭제할 수 있다")을 HTTP 요청, gRPC 메타데이터, CLI 인자마다 다시 써야 한다.

**경로 3: fact → 진리값 고정 → defeasibility의 사후 추가.**
fact는 참이고, 참인 것은 뒤집히지 않는다. 이 전제 위에 설계한 시스템은 defeasibility를 위해 바깥에 예외 장치를 붙여야 한다. LegalRuleML이 `<Override>`, `<Defeater>`, `<Superiority>` 같은 XML 요소를 핵심 모델 바깥에 더한 것이 그 예다. 입력을 claim으로 다시 정의하면 defeasibility를 처음부터 품은 설계가 쉬워지지만, 반박과 예외라는 구조적 장치(§4)는 여전히 필요하다.

**선례: JWT의 claims.**
JWT(RFC 7519[20])는 토큰 필드를 "facts"가 아니라 "claims"라 부른다. `sub`, `exp`, `iss`는 토큰 발급자의 주장이다. 서명 검증, 만료 확인, 발급자 대조는 이 데이터가 주장이라는 전제와 잘 맞는다. 이름과 검증 의무가 함께 가는 예다.

### 3.2 툴민 매핑

| 툴민 요소 | 엔진에서의 대응 | 수 |
|---|---|---|
| Claim | 판정 대상 명제. warrant 노드마다 하나이며 노드 이름으로 식별된다(`TraceEntry.Name`) | 노드당 1 |
| Ground | 판정 근거 데이터. 맥락 `ctx`로 전달되며 Ground Adapter가 채운다(`TraceEntry.Ground`) | 0..N |
| Warrant | 규칙 함수와 그 노드. 공격받을 수 있다 | 1..N |
| Backing | spec: 선언 시점에 고정되는 판정 기준(`TraceEntry.Specs`) | 규칙당 0..N |
| Qualifier | 규칙마다의 가중치 q ∈ [0, 1] **(§3.3)** | 규칙당 1 |
| Rebuttal | 공격하는 노드: 반박(Counter)과 예외(Except) **(§4.4)** | 0..N |

입력 데이터는 ground로 맥락에 실리지만, 엔진은 판정 대상 명제를 참으로 확정된 것으로 다루지 않는다. 명제가 얼마나 성립하는지는 규칙값, 공격, qualifier를 거친 판정값으로만 나온다.

### 3.3 Qualifier의 재배치: Claim에서 규칙으로

툴민의 원래 모델에서 Qualifier는 **Claim에 붙는다**. "**아마도(presumably)** 이 환자에게 페니실린을 투여해야 한다"에서 Qualifier는 주장의 확신도를 나타내는 양상 한정사(modal qualifier)이고, 화자가 Ground, Warrant, Rebuttal을 종합해 Claim에 부여하는 확신의 정도다.

본 모델은 Qualifier를 Claim에서 **각 규칙 노드(일반 규칙, 반박, 예외)**로 옮긴다. 논증 이론을 규칙 엔진에 적용하면서 생기는 공학적 조정이다.

규칙 엔진에서 "이 파일에는 함수가 3개 있다"는 관측은 ground이고, 판정할 claim은 "이 파일은 파일당 함수 하나 규칙을 어긴다" 같은 명제다. 관측 자체에 확신도를 매길 일은 드물다. 판정의 질을 좌우하는 것은 관측에서 claim으로 넘어가는 **규칙의 확신도**다.

- "파일당 함수 하나": qualifier 1.0 (확실한 규칙)
- "권장 100줄 이하": qualifier 0.7 (유연한 규칙)
- "알레르기 가능성": qualifier 0.95 (강한 반박)

이 재배치는 가중 h-Categoriser[3]와 자연스럽게 맞물린다. 노드의 초기 가중치 `w(a)` 자리에 규칙의 qualifier q가 들어가고, 최종 산출물인 판정값(verdict)이 툴민 모델에서 Qualifier가 맡던 역할, 곧 판정의 확신도를 대신한다. v0.9.0에서는 규칙 함수도 정도를 낸다(규칙값 v, §4.2). q는 규칙이 얼마나 믿을 만한지를, v는 이 사안에 규칙이 얼마나 해당하는지를 나타내며, 둘은 곱으로 결합된다(§4.3).

| | 툴민 원래 모델 | 본 모델 |
|--|--------------|--------|
| Qualifier 위치 | Claim (출력 쪽) | 각 규칙 노드 (입력 쪽) |
| 역할 | 최종 판정의 확신도 | 규칙의 가중치 q |
| 확신도 산출 | 화자가 종합 판단 | 엔진이 q, 규칙값 v, 공격과 그 세기로 계산 |
| 최종 확신도 | Qualifier 자체 | **판정값** (verdict) |

### 3.4 Warrant의 인식론적 지위

엄밀히 말하면 warrant 자체도 claim이다. "관리자만 삭제할 수 있다"는 조직의 합의이지 물리법칙이 아니고, 물리법칙도 재현성이 매우 높은 claim이다. 그러나 이 인식을 끝없이 적용하면 정당화의 무한 후퇴에 빠진다. 형식 체계에서도 공리는 체계 안에서 증명되지 않고 전제된다. 괴델의 제2불완전성 정리[10]는 산술을 담는 무모순 형식 체계가 자기 무모순성조차 체계 안에서 증명할 수 없음을 보인다. 켈젠의 근본규범(Grundnorm)[11]은 법체계의 효력 사슬이 전제된 규범에서 멈춘다는 비슷한 구조를 보여 준다.

따라서 본 모델은 실용적으로 결단한다: **warrant는 시스템 안에서 fact로 다룬다.** 이는 공리 체계가 공통으로 지는 전제이며, 모델의 적용 범위를 "warrant가 주어졌을 때 claim의 타당성을 판정하는 것"으로 한정한다.

---

## 4. 판정 모델

본 절은 `toulmin` v0.9.0의 판정 계산을 정의한다. 정의는 Go 참조 구현(`go/pkg/toulmin`)과 그 테스트를 따르며, 네 포트가 같은 계산을 구현한다(§5.5).

### 4.1 판정 그래프

판정 그래프는 다음으로 이루어진다.

- 노드 집합 N. 각 노드 x는 규칙 함수 `f_x`, spec 목록 `σ_x`, qualifier `q_x ∈ [0, 1]`(기본 1), 역할 `role(x) ∈ {rule, counter, except}`, 강도 `str(x) ∈ {defeasible, strict, defeater}`를 갖는다. `Rule`과 `Counter`로 등록한 노드는 defeasible, `Except`로 등록한 노드는 defeater다.
- 공격 간선 집합 A ⊆ N × N. 간선 (y, x)는 "y가 x를 공격한다"이며 세기 `F_yx ∈ [0, 1]`(기본 1)을 갖는다.
- 배제 간선 집합 P ⊆ N × N. 간선 (y, x)는 "y가 x를 배제한다(supersedes)"이다.

A와 P는 각각 순환이 없어야 한다. 엔진은 평가 시작 때 두 간선 집합을 따로 깊이 우선 탐색해, 순환이 있으면 규칙 함수를 하나도 실행하지 않고 오류를 돌려준다.

노드 x의 공격자 집합을 `Att(x) = {y | (y, x) ∈ A}`라 하고, 이를 둘로 나눈다.

- 예외 `E(x) = {y ∈ Att(x) | str(y) = defeater 또는 role(y) = except}`
- 반박 `R(x) = Att(x) \ E(x)`

곧 공격이 반박으로 계산될지 예외로 계산될지는 간선이 아니라 **공격자의 역할**이 정한다. 공격자가 아니고 defeater도 아닌 노드를 **warrant**라 하며, 평가 결과는 warrant마다 하나씩 나온다.

### 4.2 규칙값과 적용도

맥락 c(`ctx`)에서 규칙 함수는 규칙값 `v_x = f_x(c, σ_x) ∈ [-1, 1]`과 증거를 돌려준다. +1은 "완전히 해당", 0은 "반쯤 해당", -1은 "해당 없음"이다. 규칙값이 NaN이거나 [-1, 1] 밖이면, 또는 규칙 함수가 패닉(예외)을 일으키면 평가 전체가 오류로 끝난다.

노드 x의 적용도는 `a(x) = (v_x + 1) / 2 ∈ [0, 1]`이다. `v_x > -1`이면 x가 **적용된다**고 한다(`TraceEntry.Activated`).

옛 bool 규칙 함수도 그대로 받는다. 참은 `v = +1`, 거짓은 `v = -1`로 바뀐다. "해당 없음"은 0이 아니라 -1이라는 점이 중요하다. 0을 돌려주는 규칙은 적용도 0.5로 계산에 남으며, 0은 합법적인 값이므로 엔진은 이 실수를 검출할 수 없다.

### 4.3 판정값 계산

판정값 `verdict(x)`는 다음과 같이 재귀적으로 정한다.

- `v_x = -1`이면 `verdict(x) = -1`이다. 이때 x의 공격자는 평가하지 않는다(지연 평가).
- 그 밖에는

```
S(x)       = Σ_{y ∈ R(x)} F_yx · (verdict(y) + 1) / 2         반박: 분모 항
G(x)       = Π_{y ∈ E(x)} (1 − F_yx · (verdict(y) + 1) / 2)    예외: 적용도 게이트
a_eff(x)   = a(x) · G(x)                                       유효 적용도
raw(x)     = q_x · a_eff(x) / (1 + S(x))                       ∈ [0, 1]
verdict(x) = 2 · raw(x) − 1                                    ∈ [-1, 1]
```

- `str(x) = strict`이면 공격을 모두 무시한다: `S(x) = 0`, `G(x) = 1`.

적용되지 않는 공격자 y는 `verdict(y) = -1`이므로 기여 `(verdict(y) + 1) / 2`가 0이다. `q_x`, `a_eff(x)`가 [0, 1]에 있고 `S(x) ≥ 0`이므로 `raw(x) ∈ [0, 1]`, `verdict(x) ∈ [-1, 1]`이다. 각 warrant의 유효 적용도 `a_eff`는 `EvalResult.Applicability`로 보고된다.

**h-Categoriser와의 관계.** 모든 규칙이 bool(v ∈ {−1, +1})이고, 예외 공격이 없고, 모든 세기가 1이면 적용된 노드 x에 대해

```
raw(x) = q_x / (1 + Σ_{y ∈ Att(x), v_y = +1} raw(y))
```

이다. 곧 적용된 노드만으로 이루어진 부분 그래프 위의 가중 h-Categoriser[3]이며, 가중치는 qualifier다. 저장소의 테스트(`graph_evaluate_bool_compat_test.go`)는 이 경우 bool 규칙 그래프와 ±1을 돌려주는 등급 규칙 그래프가 모두 v0.9.0 이전의 식 `raw = q / (1 + Σ raw(attackers))`와 비트 단위로 같은 값을 냄을 확인한다.

### 4.4 반박과 예외, 공격 세기

반박과 예외는 같은 공격 간선으로 선언하지만 계산이 다르다.

- **반박(Counter)**은 결론에 맞서는 주장이다. 분모 `1 + S(x)`에 더해져 판정값을 낮춘다. 완전한 반박 하나는 공격받지 않은 완전한 주장(q = 1, v = +1)을 0(균형)으로 만들 뿐 뒤집지 않으며, 반박만으로는 `q_x · a_eff(x) > 0`인 한 판정값이 -1에 이르지 않는다.
- **예외(Except)**는 "이 경우에는 규칙이 적용되지 않는다"는 말이다. 분모가 아니라 대상의 적용도에 곱해진다. 완전한 예외 하나는 대상의 유효 적용도를 0으로, 판정값을 -1로 만든다.

이 구분은 Pollock의 반박·잠식 구분[12]과 대응한다고 볼 수 있다. 반박은 결론을 다투고, 예외는 근거에서 결론으로 넘어가는 연결이 이 사안에 적용되지 않음을 말한다. 예외의 예외도 같은 간선 하나로 표현된다. 예외 노드 자신도 공격받을 수 있고, 그 판정값이 낮아지면 게이트가 덜 닫힌다.

세기 `F`는 공격자의 기여에 곱해지는 간선의 속성이며, 노드의 강도나 qualifier와는 다른 축이다. 아래 표는 공격받지 않은 완전한 주장(q = 1, v = +1) 하나를 완전히 적용되는 공격자 하나(v = +1)가 세기 F로 공격할 때의 판정값이다. 식으로는 반박이 `(1 − F)/(1 + F)`, 예외가 `1 − 2F`이다.

| F | 0 | 0.25 | 0.5 | 0.75 | 1 |
|---|---|---|---|---|---|
| 반박 | 1 | 0.6 | 0.333 | 0.143 | 0 |
| 예외 | 1 | 0.5 | 0 | −0.5 | −1 |

표의 값은 v0.9.0의 Go 코어로 이 그래프를 만들어 `Evaluate`해 얻었고(소수 셋째 자리 반올림), F = 0, 0.5, 1의 반박과 F = 0, 0.25, 0.5의 예외는 저장소 테스트(`graph_evaluate_attack_force_test.go`)가 확인하는 값과 같다. 예외 여럿은 곱으로 겹친다. 반쯤 해당하는 예외(v = 0.5) 셋은 적용도를 `(1 − 0.75)³ ≈ 0.016`만 남긴다(`graph_evaluate_except_gate_test.go`). 반박 셋은 같은 주장을 -0.5까지만 내린다.

세기가 NaN이거나 [0, 1] 밖이면 선언 시점에 거부된다. qualifier도 같은 규칙을 따른다.

공격 간선에 가중치를 다는 착상 자체는 가중 논증 체계[17]에도 있지만, 여기서 세기는 공격을 지울지 정하는 예산이 아니라 h-Categoriser의 합 항과 게이트에 곱하는 계수다.

### 4.5 그래프 제약으로서의 강도

Nute의 분류[2]는 엔진에서 다음과 같이 작동한다.

| 강도 | 엔진에서의 의미 |
|----------|------|
| Strict | 선언된 공격을 계산에서 모두 무시한다(`S = 0`, `G = 1`) |
| Defeasible | 공격을 계산에 반영한다 |
| Defeater | 자기 결론이 없어 warrant가 되지 않는다. 공격자로서 예외 게이트로 계산되며, 자신도 공격받을 수 있다 |

v0.9.0의 공개 빌더 API에는 노드를 strict로 지정하는 메서드가 없다(§8.2). 공격 간선을 받지 않는 defeasible 규칙은 strict 규칙과 같은 판정값을 얻으므로, 사례(§6)의 strict 분류는 엔진에서 "공격 간선이 없는 규칙"으로 드러난다.

### 4.6 배제(Supersedes): 보고용 크리스프 우선권

공격 간선 외에 엔진은 **배제(Supersedes)** 간선을 제공한다. 배제 간선 집합 P에 대해

```
Excluded(x)  ⇔  ∃ y : (y, x) ∈ P  그리고  a_eff(y) > 0
```

이다. 배제자 y의 유효 적용도는 예외 게이트를 반영하고 반박은 반영하지 않으며, strict 배제자는 `a_eff = a`다. 적용되지 않는 배제자는 `a_eff = 0`이다. 유효 적용도가 0.3처럼 일부만 남아도 대상 전체를 배제한다.

배제는 판정값을 바꾸지 않는다. 배제된 노드는 trace에 `Excluded`로 표시될 뿐이며, 판정값은 그대로이고 그 노드가 다른 노드에 가하는 공격과 배제도 그대로 유효하다. 배제 간선은 h-Categoriser 순회에 들어가지 않는다. 표시를 따라 노드를 집계에서 빼는 일은 하류 집계(예: 할인 하나만 고르기)가 맡는다.

저장소 테스트의 예(`run_supersedes_except_gate_test.go`)를 보자. VIP 할인이 일반 할인을 배제한다. "정지된 VIP" 예외가 VIP 할인을 완전히 막으면(`a_eff = 0`) 배제가 풀려 일반 할인이 집계에 남는다. 세기 0.5의 예외(`a_eff = 0.5`)나 VIP 할인에 대한 반박(`a_eff = 1`)은 배제를 유지한다. 배제 간선의 순환은 서로 배제하는 작성 오류로 보고 평가 전에 거부한다.

이론적으로 배제는 공격과 다른 층위의 장치다. Raz의 배제적 이유(exclusionary reason)[13]는 저울에서 상대 이유를 눌러 이기는 1차 이유가 아니라 특정 이유를 저울에 올리지 말라는 2차 이유이며, 배제 간선은 이 집계 배제의 구조를 따른다. defeasible logic의 우선관계[18]도 규칙 사이의 크리스프한 우선을 표현한다. 다만 라벨링 의미론[14]에서 OUT 라벨을 받은 논증이 공격력을 잃는 것과 달리, 배제된 노드는 공격력과 배제권을 잃지 않는다. 배제는 점진 계산을 그대로 두고 그 위에 얹은 보고용 크리스프 층이다.

### 4.7 판정값의 해석

| 판정값 | 의미 |
|---------|------|
| +1 | 완전히 적용되고(`q · a_eff = 1`) 유효한 반박이 없다 |
| 0 | 균형. 예: 완전한 주장 하나와 완전한 반박 하나, 또는 공격 없이 `q · a_eff = 1/2` |
| −1 | 유효 적용도가 0(완전한 예외)이거나 q = 0 |

판정값은 **주장이 성립하는 정도**다. 양수 쪽이 긍정(주장이 지지됨), 음수 쪽이 부정(주장이 반박됨)이다. 적용되지 않는 warrant(`v = -1`)는 판정값 -1을 내는 대신 결과에서 빠진다. 곧 "해당 없음"과 "반박됨"이 구별된다. 엔진은 값만 계산하며, 문턱 같은 도메인별 해석은 사용자가 정한다(예: 접근 제어는 `verdict > 0`이면 허용). 직관에 맞추려면 규칙의 주장을 긍정형으로 설계하기를 권한다("접근 허용"에서 +1은 허용). 부정형("위반 감지") 규칙도 유효하지만 +1이 위반을 뜻하게 되어 해석이 반직관적이다.

---

## 5. 구현

### 5.1 증거와 판정값을 돌려주는 함수로서의 규칙

모든 규칙(일반 규칙, 반박, 예외)은 같은 시그니처를 갖는다. Go에서는

```go
func(ctx toulmin.Context, specs toulmin.Specs) (float64, any) // 규칙값 v ∈ [-1, 1], 증거
func(ctx toulmin.Context, specs toulmin.Specs) (bool, any)    // 옛 시그니처: 참 → +1, 거짓 → -1
```

이다. 첫 매개변수 `ctx`는 `Get(key string) (any, bool)`과 `Set(key string, value any)`를 가진 `Context` 인터페이스이고, 둘째 매개변수 `specs`는 규칙의 판정 기준을 런타임 값으로 전달한다. 둘째 반환값은 도메인 특화 증거(예: 위반 위치와 기대값)이며, 엔진은 이를 해석하지 않고 결과와 trace로 넘긴다. 위치·기대값·실제값을 담는 선택적 관례 타입 `Fact{Where, Expected, Actual}`도 있다.

| | Context | Spec |
|--|---------|------|
| 툴민 요소 | Ground | Backing |
| 정의 | 판정 대상의 사실(Get/Set 접근) | 규칙의 판정 기준 |
| 시점 | 요청마다 다름 | 선언 시점에 고정 |
| 예시 | 파일의 함수 개수, 사용자의 역할 | 허용 함수 개수의 임계값, 역할 이름 |

```go
// 이진 규칙: bool 시그니처로 쓴다.
func CheckOneFileOneFunc(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
    gf, _ := ctx.Get("file")
    f := gf.(*FileGround)
    if len(f.Funcs) > 1 {
        return true, &Evidence{Got: len(f.Funcs), Expected: 1}
    }
    return false, nil
}

// spec을 판정 기준으로 쓰는 규칙
func CheckLineCount(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
    gf, _ := ctx.Get("file")
    f := gf.(*FileGround)
    limit := specs[0].(*LineLimit)
    if f.Lines > limit.Max {
        return true, &Evidence{Got: f.Lines, Limit: limit.Max}
    }
    return false, nil
}

// 등급 규칙: 분류기 점수 p ∈ [0, 1]을 규칙값 v = 2p − 1로 옮긴다.
func LooksLikeSpam(ctx toulmin.Context, specs toulmin.Specs) (float64, any) {
    p, _ := ctx.Get("spam_score")
    score := p.(float64)
    return 2*score - 1, score
}
```

warrant, 반박, 예외의 구분은 함수 시그니처가 아니라 그래프에서의 위치와 등록 역할로 정해진다.

| 역할 | 그래프에서의 위치 |
|------|----------------|
| Warrant (`Rule`) | 공격받을 수 있고, 다른 노드를 공격하지 않는 노드. 결과를 낸다 |
| 반박 (`Counter`) | 다른 노드를 공격하는 노드(자기 결론이 있다). 분모 항으로 계산된다 |
| 예외 (`Except`) | 자기 결론 없이 공격하는 노드. 적용도 게이트로 계산된다 |

### 5.2 Graph API: 로직과 구조의 분리

Graph API는 규칙 로직(함수)과 그래프 구조(역할, qualifier, 공격, 배제)를 분리한다.

```go
g := toulmin.NewGraph("api:access")
auth    := g.Rule(isAuthenticated)
blocked := g.Counter(isIPBlocked)
exempt  := g.Except(isInternalIP)
blocked.Attacks(auth)        // 반박, 세기 1
exempt.Attacks(blocked)      // 반박에 대한 예외: 차단 규칙의 적용도를 깎는다
hint := g.Counter(isSuspicious)
hint.Attacks(auth, 0.3)      // 세기 0.3의 약한 반박
```

`g.Rule(fn)`, `g.Counter(fn)`, `g.Except(fn)`은 `*Rule` 참조를 돌려주고, 빌더 메서드 `.With(spec)`(spec 추가)과 `.Qualifier(q)`(기본 1.0)를 이어 쓸 수 있다. spec은 `Spec` 인터페이스(`SpecName() string`, `Validate() error`)를 구현해야 하며, 함수 필드를 가진 spec은 거부된다. `Attacks(target, force...)`는 공격 간선을, `Supersedes(target)`은 배제 간선을 선언한다.

같은 함수를 서로 다른 spec으로 등록하면 별개의 노드가 된다. 노드 식별자는 함수 식별자에 spec들의 이름과 값 직렬화를 붙여 만든다.

```go
g := toulmin.NewGraph("access")
admin  := g.Rule(isInRole).With(&RoleSpec{Role: "admin"})
editor := g.Rule(isInRole).With(&RoleSpec{Role: "editor"}).Qualifier(0.8)
```

함수가 식별자이므로 문자열 이름이 필요 없다. 같은 함수를 다른 그래프에서 다른 역할과 공격 관계로 다시 쓸 수 있다.

```go
// 같은 IsAdult 함수, 다른 그래프, 다른 공격 관계
votingGraph := toulmin.NewGraph("voting")
vw := votingGraph.Rule(IsAdult)
vr := votingGraph.Counter(HasCriminalRecord)
vr.Attacks(vw)

contractGraph := toulmin.NewGraph("contract")
cw := contractGraph.Rule(IsAdult)
cr := contractGraph.Counter(IsBankrupt)
cr.Attacks(cw)
```

배제는 다음과 같이 선언한다.

```go
general := g.Rule(generalDiscount)
vip     := g.Rule(vipDiscount)
vip.Supersedes(general)               // vip의 a_eff > 0이면 general에 Excluded 표시
g.Except(suspendedVIP).Attacks(vip)   // vip를 완전히 막으면 배제도 풀린다
```

구조적 메타데이터(역할, qualifier, 공격, 세기, 배제)는 모두 함수가 아니라 그래프에 속하고, 함수는 판정 로직만 담는다.

### 5.3 평가 절차

`Evaluate`는 다음 순서로 판정한다.

```
0. 순환 감지: 공격 간선과 배제 간선을 각각 DFS. 순환이면 오류(규칙 함수 실행 전)
1. 각 warrant에서 시작한다
2. 규칙 함수 f(ctx, specs)를 실행해 규칙값 v를 얻는다
   (패닉, NaN, 범위 밖이면 평가 오류)
3. v = -1이면 판정값 -1, 공격자는 평가하지 않는다. warrant라면 결과에서 빠진다
4. v > -1이고 strict가 아니면 공격자를 재귀적으로 평가해,
   예외는 게이트 G에, 반박은 합 S에 넣는다
5. a_eff = a·G, raw = q·a_eff/(1+S), verdict = 2·raw − 1. 노드별로 캐시한다
6. trace를 끄면 한 번의 Evaluate 안에서 각 노드의 함수는 많아야 한 번 실행되고
   결과는 warrant 사이에 공유된다. trace를 켜면 warrant마다 상태를 비우고
   다시 평가하므로 함수가 여러 번 실행될 수 있다
7. trace를 켜면 warrant마다 배제 패스를 돌린다: 배제자를 평가해 a_eff > 0이면
   대상의 TraceEntry에 Excluded를 표시한다(판정값은 그대로)
```

지연 평가 때문에 trace를 끈 `Evaluate`는 warrant에서 공격 사슬로 닿는 규칙만 실행한다.

### 5.4 결과, trace, 실행

```go
results, err := g.Evaluate(ctx)                                   // 판정값 + 증거
results, err = g.Evaluate(ctx, toulmin.EvalOption{Trace: true})   // + 노드별 trace
results, err = g.Evaluate(ctx, toulmin.EvalOption{Duration: true}) // + 실행 시간(trace 자동 켜짐)
```

```go
type EvalResult struct {
    Name          string       `json:"name"`
    Verdict       float64      `json:"verdict"`
    Evidence      any          `json:"evidence,omitempty"`
    Applicability float64      `json:"applicability"` // a_eff
    Trace         []TraceEntry `json:"trace"`
}

type TraceEntry struct {
    Name        string        `json:"name"`         // = Claim
    Role        string        `json:"role"`
    Activated   bool          `json:"activated"`    // v > -1
    RuleVerdict float64       `json:"rule_verdict"` // 규칙값 v
    Qualifier   float64       `json:"qualifier"`
    Verdict     float64       `json:"verdict"`
    Evidence    any           `json:"evidence,omitempty"`
    Ground      any           `json:"ground,omitempty"` // = Ground (ctx)
    Specs       Specs         `json:"specs,omitempty"`  // = Backing
    Duration    time.Duration `json:"duration,omitempty"`
    Excluded    bool          `json:"excluded,omitempty"`
    Level       Level         `json:"level"`
}
```

`Evaluate`의 trace는 warrant마다 평가된 노드의 역할, 규칙값, 적용 여부, qualifier, spec, 증거, 배제 표시를 남긴다. spec이 trace에 들어가므로 같은 함수를 다른 spec으로 등록한 노드들의 판정 기준을 구분할 수 있다. 판정값의 도출 과정을 이렇게 감사할 수 있다는 것이 심볼릭 추론의 장점이다. 다만 `Evaluate` 경로의 trace 항목에는 노드별 최종 판정값(`Verdict`)이 채워지지 않는다. 노드별 판정값과 배제 표시를 포함한 그래프 전체의 trace는 `Run`이 돌려준다. trace를 끈 `Evaluate`의 결과에는 배제 표시가 없다. 판정 방법 `Method`의 기본값 `Matrix`는 지연 재귀 h-Categoriser이며, `Recursive`는 아직 구현되지 않아 지정하면 오류를 돌려준다.

`Evaluate`는 부수효과가 없다. `Run`은 그래프 전체를 평가한 뒤, 적용되고 판정값이 양수인 노드에 등록된 처리기(`RunOn`)를 등록 순서로 실행하고, 하위 그래프가 선언된 노드는 같은 맥락으로 그 그래프를 `Run`한다(실행 합성, 순환 거부, 깊이 64 제한). Go 포트는 판정값과 별개인 보고 등급 `Level`(`Unset`/`Fail`/`Review`)도 노드에 달 수 있다. 규칙의 존재 이유를 기록할 때는 일반 주석을 쓴다. 별도의 주석 문법은 없다.

### 5.5 네 포트

같은 판정 계산이 저장소의 네 디렉터리에 구현되어 있다: `go/`(참조 구현, 모듈 `github.com/park-jun-woo/toulmin/go`), `ts/`, `py/`, `c/`. Go 모듈은 GEUL 라이브러리 `github.com/geul-org/gogeul`에 의존하며, 이 의존은 TANGEUL 툴체인이 쓴다.

| | Go | TypeScript | Python | C |
|---|---|---|---|---|
| 등급 규칙 | `(float64, any)` | `[number, unknown]` | `(float, Any)` | `rc_verdict_func` + `rc_graph_*_verdict` |
| bool 규칙 구분 | 등록 시 함수 타입으로 | 호출마다 반환값 타입으로 | 호출마다 반환값 타입으로 | 등록 함수를 따로 둔다(`rc_graph_rule` 등) |
| 공격 세기 | `Attacks(t, f)` | `attacks(t, f)` | `attacks(t, force=f)` | `rc_rule_attacks_with_force(r, t, f)` |
| 배제 | `Supersedes(t)` | `supersedes(t)` | `supersedes(t)` | `rc_rule_supersedes(r, t)` |
| 유효 적용도 | `EvalResult.Applicability` | `applicability` | `applicability` | `rc_eval_result.applicability` |
| 잘못된 세기·qualifier | 선언 시 panic | 예외 | `ValueError` | 지연 빌더 오류(평가 시 `RC_ERR_BUILDER`), 세기는 호출에서도 `RC_ERR_INVALID_ARG` |
| 잘못된 규칙값 | `Evaluate` 오류 | 예외 | `RuntimeError` | `RC_ERR_RULE` |

판정값 계산(§4.3)과 배제 판정(§4.6)은 네 포트가 같다. 다른 점은 다음과 같다. TypeScript와 Python에서는 수도 bool도 아닌 반환값이 "거짓"이 아니라 오류다. C의 규칙 함수는 상태 코드를 돌려주며 실패하면 `RC_ERR_RULE`이 된다. `Evaluate`에서 trace를 켰을 때 항목의 역할 표시는 Go가 등록 역할을, 나머지 세 포트가 그래프 위치로 추론한 역할을 쓴다(일반 규칙으로 등록한 노드가 다른 노드를 공격하면 표시가 다르다. 판정값에는 영향이 없다). 보고 등급 `Level`과 TANGEUL 툴체인, 도메인 프레임워크 패키지는 Go 포트에만 있다.

### 5.6 선언 표면: TANGEUL

코어 엔진은 별도의 DSL을 요구하지 않는다. 규칙은 각 언어의 함수다. 그 위에 선택적으로, 비개발자가 읽고 감사할 수 있는 마크다운 정책 언어 TANGEUL(Toulmin Amgoud Nute Graph - General Encoding Unified Language)이 있다(Go 포트 `go/pkg/tangl`, `go/pkg/tangeul`, CLI `tangl`). TANGEUL 문서는 별도 런타임 없이 `toulmin.Graph`로 컴파일된다.

- **규칙과 사안의 분리.** `tangl:Rules` 절이 판정의 구현(조건식 규칙, 또는 호스트 함수를 쓰는 `uses` 규칙)과 판정 입력(`with` spec, `grounded on` 근거)을 가진다. `tangl:Cases` 절은 같은 이름을 인용해 역할(일반·반박·예외), qualifier, 보고 등급만 붙인다.
- **정의.** 문자열 목록 정의(`` `domains` means ["a.com", "b.com"] ``)는 정확 일치 소속 검사를 하는 내장 spec이 되고, 파일 정의(`` `domains` reads `JSON list` from "./domains.json" ``)는 평가나 코드 생성 전에 한 번만 읽어 고정한다.
- **공격과 세기.** `` don't `X` when `Y` ``는 세기 1의 공격, `` `Y` attacks `X` with force 0.5 ``는 `Y.Attacks(X, 0.5)`로 컴파일된다. 세기는 [0, 1], 소수 넷째 자리까지다. 같은 공격 문장이라도 공격자가 예외 규칙이면 게이트로, 그 밖이면 분모 항으로 계산된다(§4.4).
- **배제.** `` `X` supersedes `Y` ``는 배제 간선이고, `` supersede in order: `A` > `B` > `C` ``는 앞→뒤의 모든 쌍(A→B, A→C, B→C)으로 펼쳐진다.
- **원본은 `.geul`.** `tangl encode`는 저작한 `.md` 표면을 로케일에 독립인 `.geul` 파일, 곧 TOULMIN 코덱을 따르는 GEUL 문서 하나[21]로 흡수한다. 이 파일은 AST와 저작한 원문 바이트를 함께 담으며, 공격 세기는 ATTACKS 레코드의 `F_FORCE` 슬롯에 실린다. 검사·효과 분석·코드 생성은 이 스트림 위에서 돌고, `tangl decode`는 원문을 바이트 그대로 되살리거나(모드 1) AST에서 영어·한국어 표면을 다시 만든다(모드 2). 영어와 한국어 표면은 같은 AST로 파싱된다.

TANGEUL의 조건식 규칙과 사안 합성(`checking`)은 참·거짓만 낸다(+1/−1). 등급 판정은 `(float64, any)` 호스트 함수를 `uses`로 쓸 때만 나오며, 참조 평가기 `pkg/tangl/eval`은 아직 bool 전용이다. 다른 문서의 사안을 참조하는 `checking`은 파서와 코드 생성기는 받지만 `.geul`에는 아직 실을 수 없다.

---

## 6. 사례 1: filefunc 규칙의 툴민 변환

### 6.1 배경

filefunc은 LLM 네이티브 Go 개발을 위한 코드 구조 컨벤션이자 CLI 도구다. 변환 당시 filefunc은 22개 규칙(F·Q·A·C 계열)을 Go 함수로 구현해 AST를 대상으로 컴파일타임에 평가했고, 모든 규칙을 툴민 warrant로 옮겼다. 아래 분류와 수치는 변환 당시의 규칙 집합 기준이다.

### 6.2 강도 분류

| 강도 | 수 | 비율 | 예 |
|----------|---|------|------|
| Strict | 15 | 68% | F1–F4, A1–A3 등 |
| Defeasible | 4 | 18% | Q1, Q2, Q3, C4 |
| Defeater | 3 | 14% | F5, F6, 테스트 파일 예외 |

대부분의 규칙이 strict다. 코드 구조 컨벤션은 본디 예외를 줄이려 하기 때문이다. defeasible 규칙은 맥락에 따라 임계값이 달라지는 것들(깊이 제한 Q1–Q2, 크기 제한 Q3)이다. §4.5에서 말했듯 엔진에서 strict 규칙은 공격 간선을 받지 않는 규칙으로, defeater는 `Except` 노드로 표현된다.

### 6.3 규칙 집합을 적용한 결과

| 프로젝트 | 파일 수 (전→후) | 평균 LOC/파일 (전→후) | SRP 위반 해소 | depth 위반 해소 |
|---------|---------------|---------------------|-------------|--------------|
| filefunc | — (처음부터 준수) | 25.1 | 0 | 0 |
| fullend | 87→1,260 | 244→25.4 | 66→0 | 148→0 |
| whyso | 12→99 | 147.8→24.4 | 12→0 | 23→0 |

이 수치는 저자가 각 프로젝트를 filefunc 규칙에 맞게 고친 전후를 잰 것이다. 툴민 그래프로 표현한 규칙 집합이 실제 코드베이스에 적용되어 위반을 0으로 만들 수 있었음을 보여 주지만, 엔진 설계(툴민 변환) 자체의 효과를 if-else 구현과 견주어 분리하지는 않는다(§8.2).

### 6.4 툴민 변환의 이점

1. **Spec의 1급 지위**: 판정 기준(임계값, 정책 등)이 런타임 spec 값으로 함수에 전달된다. 같은 함수를 다른 spec으로 등록하면 별개의 규칙이 된다.
2. **반박과 예외의 구조화**: 예외 조건(F5 테스트 파일, F6 그룹 상수)이 명시적 공격 간선으로 선언되고, 코드 속 if 분기에 묻히지 않는다.
3. **강도 분류**: "이 규칙에는 예외가 없다"(strict)와 "이 규칙은 무력화될 수 있다"(defeasible)가 그래프 구조로 드러난다.

---

## 7. 사례 2: fullend의 시점 선택적 설계

### 7.1 배경

fullend는 9개 SSOT(OpenAPI, SSaC, Rego 등)의 정합성을 검증하고 코드를 생성하는 시스템이다. 18개 교차 검증 규칙이 Go 함수로 구현되어 있고, 런타임 정책은 별도의 Rego 파일이다. filefunc(§6)이 한 시점(컴파일타임)에서만 작동하는 반면, fullend는 시점 선택적 설계가 필요한 사례다.

### 7.2 문제: 규칙 중복

같은 논리 규칙이 두 번 작성된다(개략).

**컴파일타임** (Go 교차 검증):
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

### 7.3 툴민 설계

이진 판정 규칙이므로 bool 시그니처로 쓴다(§5.1).

```go
// OWASP API Security Top 10 (2023) API2: Broken Authentication
func AuthEndpointRequiresClaims(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
    ep, _ := ctx.Get("endpoint")
    g := ep.(*EndpointGround)
    if g.Security.Contains("bearerAuth") && !g.PolicyRule.References("claims") {
        return true, &Evidence{Endpoint: g.Path, Missing: "claims reference"}
    }
    return false, nil
}

// 공개 API endpoint는 설계상 인증이 필요 없다
func PublicEndpointException(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
    ep, _ := ctx.Get("endpoint")
    g := ep.(*EndpointGround)
    return g.Annotation.Contains("x-public"), nil
}

// 그래프 정의
g := toulmin.NewGraph("endpoint-auth")
w := g.Rule(AuthEndpointRequiresClaims)
d := g.Except(PublicEndpointException)
d.Attacks(w)   // 공개 endpoint면 유효 적용도 0 → 판정값 -1
```

**컴파일타임 Ground Adapter**: `EndpointGround.Security` ← OpenAPI 명세 파싱, `EndpointGround.PolicyRule` ← Rego AST 파싱.

**런타임 Ground Adapter**: `EndpointGround.Security` ← 요청 라우트 미들웨어, `EndpointGround.PolicyRule` ← OPA 평가.

하나의 warrant 함수가 컴파일타임 교차 검증과 런타임 정책 적용을 함께 맡는다. 그래프 구조는 같고 Ground Adapter만 바뀐다. 함수 자체도 그래프 사이에서 다시 쓸 수 있어, 같은 `AuthEndpointRequiresClaims`가 다른 맥락에서 다른 공격 관계에 참여할 수 있다. 이 사례는 설계 수준의 사례이며, 두 Ground Adapter가 같은 입력에 같은 판정을 내는지에 대한 정량 평가는 담지 않는다.

---

## 8. 논의

### 8.1 본 논문이 주장하지 않는 것

본 논문은 구성 요소의 이론적 독창성을 주장하지 않는다. 구조적 기초는 툴민(1958)[1], 강도 분류는 Nute(1994)[2], 판정 계산의 뼈대는 h-Categoriser(Besnard·Hunter 2001[15], 가중판 Amgoud et al. 2017[3]), 반박과 잠식의 구분은 Pollock(1987)[12], 어댑터 패턴은 Cockburn(2005)[9]과 StaRVOOrS(2015)[7]에서 왔다. 규칙값과 적용도, 예외 게이트, 공격 세기, 보고용 배제는 이 엔진이 더한 공학적 확장이며, 그 형식 성질은 이 논문에서 분석하지 않는다. 우리의 기여는 **따로 발전한 이 조각들이 이어진다는 것**, 곧 툴민의 논증 모델이 소프트웨어 규칙 엔진의 빠진 설계 기초라는 관점을 제시하고, **구현으로 이를 보인 것**이다.

### 8.2 한계와 타당성 위협

**모델과 구현**

- 규칙 본문(호스트 함수)은 엔진에 불투명하므로 부분 평가와 정적 분석을 지원하지 않는다. 그래프 구조와 TANGEUL 조건식 규칙은 `.geul`로 직렬화되지만 호스트 함수의 본문은 그렇지 않다. 네 포트가 같은 계산을 구현해도 규칙 함수는 언어마다 다시 써야 한다.
- 예외 게이트, 공격 세기, 배제가 점진 의미론 문헌의 원리들을 만족하는지 분석하지 않았다. [-1, 1] 변환은 선형이며, 복잡한 규칙 상호작용에는 더 정교한 집계가 필요할 수 있다.
- 공격 간선과 배제 간선의 순환을 오류로 거부하므로, 서로 공격하는 논증처럼 순환 그래프가 필요한 경우를 다루지 못한다.
- 규칙값 0은 "반쯤 해당"이다. bool 습관대로 "아니오"에 0을 돌려주는 실수는 조용히 판정을 바꾸며, 엔진은 이를 검출할 수 없다.
- 배제는 보고용 표시다. 판정값을 바꾸지 않으며, 하류 집계가 표시를 따라야 효과가 있다. trace를 끈 `Evaluate`의 결과에는 표시가 없다.
- v0.9.0의 공개 빌더 API와 TANGEUL에는 노드를 strict로 지정하는 수단이 없다. 엔진 내부의 strict 처리는 구현되어 있지만 사용자는 공격 간선을 두지 않는 방식으로만 strict 규칙을 표현할 수 있다.
- Claim 모델링에 추가 작업이 필요하다. 두 주어 사이의 관계 단언은 단순한 (주어, 속성, 값) 삼중쌍의 표현력을 넘는다.

**평가**

- 사례가 저자 자신의 프로젝트에 한정되어 있다. 외부 오픈소스 프로젝트에 대한 검증은 향후 과제다.
- §6.3의 수치는 규칙 집합을 적용한 효과이며, 툴민 변환 자체의 효과를 분리하지 않는다. 측정 명령과 환경은 이 판에 싣지 않았다.
- §7은 설계 사례이며 정량 평가가 없다.
- §3.1의 "fact" 명명에 관한 논증은 개념적이며 경험적으로 검증하지 않았다.

### 8.3 기존 시스템과의 관계

본 엔진은 Rego, Drools, Semgrep을 대체하지 않고 **상위 추상화**를 제공한다. Rego는 런타임 Ground Adapter의 백엔드로 쓸 수 있고, Semgrep 규칙은 컴파일타임 Ground Adapter에서 다시 쓸 수 있다. 기존 시스템은 특정 시점의 Ground Adapter 구현체가 된다.

---

## 9. 결론

생산 규칙 시스템이 나온 이래 소프트웨어 규칙 엔진은 검증 대상을 "fact"로 다루고 이 전제 위에 설계를 쌓아 왔다. 툴민의 논증 모델(1958)은 이미 알맞은 구조, 곧 Claim(fact가 아닌), Ground, Warrant, Backing, Qualifier, Rebuttal을 제공했지만 규칙 엔진 설계에는 거의 쓰이지 않았다.

우리는 이 간극을 메우는 규칙 엔진 `toulmin`을 구현했다. 규칙은 맥락과 spec을 받아 규칙값과 증거를 돌려주는 함수다. spec은 선언 시점에 고정되는 판정 기준(Backing)이고, 맥락은 요청마다 달라지는 판정 대상의 사실(Ground)이다. 엔진은 규칙을 공격 그래프로 조직하고, 가중 h-Categoriser를 넓힌 계산으로 [-1, 1]의 판정값을 낸다. 반박은 분모로, 예외는 적용도 게이트로 계산되고, 공격 간선에는 세기를 달 수 있으며, 배제는 판정값을 건드리지 않는 크리스프 우선권을 보고한다. 같은 계산이 Go, TypeScript, Python, C 네 포트에 있고, TANGEUL이 코어 위에 감사 가능한 선언 표면과 `.geul` 원본 형식을 더한다. trace는 어떤 규칙이 어떤 역할로, 어떤 spec과 증거로 판정에 참여했는지 남긴다.

코어에는 별도의 DSL이 필요 없다. v0.9.0의 Go 코어(`go/pkg/toulmin`)는 테스트를 뺀 Go 파일 95개, 빈 줄과 `//` 주석 줄을 뺀 1,078줄이다(커밋 `4a6e568`, Linux에서 `find go/pkg/toulmin -name '*.go' ! -name '*_test.go' | xargs cat | grep -v '^\s*//' | grep -v '^\s*$' | wc -l`로 잰 값). 같은 규칙 함수를 다른 그래프에서 다른 공격 관계로, 또는 다른 spec으로 다시 쓸 수 있어 규칙 로직과 그래프 구조가 분리된다. filefunc의 22개 규칙 변환은 툴민 모델이 규칙 강도 분류, 구조화된 defeasibility, 명시적 spec을 자연스럽게 수용함을 보여 주고, fullend 사례는 Ground Adapter로 같은 규칙을 여러 시점에서 다시 쓰는 설계를 보여 준다.

`toulmin` 소프트웨어는 MIT 라이선스로 공개된다.

---

## 참고문헌

[1] Toulmin, S. *The Uses of Argument*. Cambridge University Press, 1958.

[2] Nute, D. "Defeasible Logic." In D.M. Gabbay, C.J. Hogger, J.A. Robinson (eds.), *Handbook of Logic in Artificial Intelligence and Logic Programming*, Vol. 3, Oxford University Press, 1994.

[3] Amgoud, L., Ben-Naim, J., Doder, D., and Vesic, S. "Acceptability Semantics for Weighted Argumentation Frameworks." *IJCAI 2017*, pp. 56–62, 2017.

[4] Gabriel, V.O., Panisson, A.R., Bordini, R.H., Adamatti, D.F., Billa, C.Z. "Reasoning in BDI agents using Toulmin's argumentation model." *Theoretical Computer Science*, 805:76–91, 2020.

[5] Modgil, S., and Prakken, H. "The ASPIC+ Framework for Structured Argumentation: A Tutorial." *Argument and Computation*, 5(1):31–62, 2014.

[6] Dung, P.M. "On the Acceptability of Arguments and its Fundamental Role in Nonmonotonic Reasoning, Logic Programming and n-person Games." *Artificial Intelligence*, 77(2):321–357, 1995.

[7] Ahrendt, W., et al. "StaRVOOrS: A Tool for Combined Static and Runtime Verification of Java." *Runtime Verification*, LNCS 9333, 2015.

[8] Ahrendt, W., et al. "A Specification Language for Static and Runtime Verification of Data and Control Properties." *FM 2015*, LNCS 9109, 2015.

[9] Cockburn, A. "Hexagonal Architecture (Ports and Adapters)." 2005.

[10] Gödel, K. "Über formal unentscheidbare Sätze der Principia Mathematica und verwandter Systeme I." *Monatshefte für Mathematik und Physik*, 38(1):173–198, 1931.

[11] Kelsen, H. *Reine Rechtslehre (Pure Theory of Law)*. Franz Deuticke, Wien, 1934.

[12] Pollock, J.L. "Defeasible Reasoning." *Cognitive Science*, 11(4):481–518, 1987.

[13] Raz, J. *Practical Reason and Norms*. Hutchinson, London, 1975. (2nd ed. Oxford University Press, 1999.)

[14] Caminada, M. "On the Issue of Reinstatement in Argumentation." *JELIA 2006*, LNCS 4160, Springer, 2006.

[15] Besnard, P., and Hunter, A. "A logic-based theory of deductive arguments." *Artificial Intelligence*, 128(1–2):203–235, 2001.

[16] Amgoud, L., and Ben-Naim, J. "Ranking-based Semantics for Argumentation Frameworks." *SUM 2013*, LNCS 8078, 2013.

[17] Dunne, P.E., Hunter, A., McBurney, P., Parsons, S., and Wooldridge, M. "Weighted argument systems: Basic definitions, algorithms, and complexity results." *Artificial Intelligence*, 175(2):457–486, 2011.

[18] Antoniou, G., Billington, D., Governatori, G., and Maher, M.J. "Representation results for defeasible logic." *ACM Transactions on Computational Logic*, 2(2):255–287, 2001.

[19] Verheij, B. "Evaluating arguments based on Toulmin's scheme." *Argumentation*, 19(3):347–371, 2005.

[20] Jones, M., Bradley, J., and Sakimura, N. "JSON Web Token (JWT)." RFC 7519, IETF, 2015.

[21] GEUL 문법: TOULMIN 코덱. https://geul.org/grammar/toulmin/
