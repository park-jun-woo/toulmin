# test/ — TANGL v0.3 검증용 예제 문서

`pkg/tangl/{ast,parser,validate,effects,gen}` + `cmd/tangl`(`internal/tanglcli`)
구현을 스펙(`files/TANGL문법v0.3.md`, 읽기 전용)에 없는 새 도메인으로 검증하는
예제 모음이다. `pkg/tangl/parser/testdata/{americano,access,transfer}.md`의
예시는 재사용하지 않았다.

모든 명령은 저장소 루트에서 실행한다: `go run ./cmd/tangl <sub> test/<file>.md ...`

## 정상 문서 (check 통과)

### 1. `library-renewal.md` — checking 합성 + Rules의 either/and/or/not

도서관 대출 갱신 심사. `tangl:Rules`에서 `either`/`and`/`or`/`not` 중첩
조건 트리로 인라인 규칙 2개(`no blocking overdue`, `has pending hold from
others`)를 정의하고, Cases에서 이름을 그대로 재사용하는 노드로 참조한다
(using 절 없는 노드 = 같은 이름의 Rules 함수 바인딩). `can process renewal
request`는 `checking`으로 `can renew loan`/`can waive late fee` 두 판정을
순수 합성한다.

검증한 것:
- `check test/library-renewal.md` → `ok` (경고 없음)
- `effects ... "renew"` → `catalog.stampRenewal once`, `log.RenewalResult` 2건
- `effects ... "renewal eligibility"` (check 엔드포인트) → **공집합**
- `effects ... "process renewal"` → `log.ProcessSummary` 1건만. `can renew
  loan`의 `stampRenewal`/`RenewalResult`는 **포함되지 않음** — `checking`
  경로가 Evaluate로 컴파일되어 순수하다는 것을 직접 확인
- `ast`로 Rules의 조건 트리가 `Logic{or}` → `Either` → `Logic{and}` /
  `Not{Compare}`로 정확히 중첩 파싱됨을 확인

```
$ go run ./cmd/tangl check test/library-renewal.md
ok

$ go run ./cmd/tangl effects test/library-renewal.md "renew"
do catalog.stampRenewal once (can renew loan / no blocking overdue)
do log.RenewalResult         (can renew loan / no blocking overdue)

$ go run ./cmd/tangl effects test/library-renewal.md "renewal eligibility"
(공집합 — check 엔드포인트는 항상 부수효과 0)

$ go run ./cmd/tangl effects test/library-renewal.md "process renewal"
do log.ProcessSummary  (can process renewal request / renewal allowed)
```

### 2/3. `newsletter-tick-warn.md` / `newsletter-tick-fixed.md` — once 없는 do의 Lint 경고 대비

같은 시나리오(구독 다이제스트 발송)를 두 파일로 대비시켰다.

- `newsletter-tick-warn.md`: `tangl:Internal`의 `every 15m`이 `can send
  digest`를 직접 `run`하는데, 그 안의 `do mailer.sendDigest`에 `once`도
  완료 가드도 없다 → **Lint 경고 발생**(§틱 멱등성 요구). 경고는 에러가
  아니므로 `check`는 여전히 exit 0.
- `newsletter-tick-fixed.md`: 같은 `do`에 `once`를 붙이고, 추가로 `already
  sent today`라는 counter 노드 + `don't` 무력화 간선으로 "완료 상태를
  반영하는 가드"(§틱 멱등성 요구의 (c))까지 이중으로 방어한다 → 경고 0건.

```
$ go run ./cmd/tangl check test/newsletter-tick-warn.md
test/newsletter-tick-warn.md:24: warning: case "can send digest" do mailer.sendDigest has no once guard and is reachable from tangl:Internal
ok

$ go run ./cmd/tangl check test/newsletter-tick-fixed.md
ok
```

**구현 관찰(스펙과의 미묘한 차이, 버그는 아님)**: 스펙 §틱 멱등성 요구는
무가드 `do`가 (a)멱등, (b)`once`, (c)완료 가드 중 하나로 보호되면 안전하다고
서술하지만, 실제 Lint 구현(`pkg/tangl/validate/unguardedDoWarnings`)은
**오직 `once` 플래그만** 검사한다 — (a)/(c)로 실질적으로 보호되는 `do`라도
`once`가 없으면 여전히 경고한다. `newsletter-tick-fixed.md`에서 `already
sent today` counter 가드를 걸어도 `logDigestSent`에 `once`를 빼면 경고가
그대로 뜨는 것으로 직접 확인했다(재현: 해당 줄에서 `once`만 제거하고
`check` 재실행). 문서화된 (a)/(c) 경로는 현재 Lint가 정적으로 인식하지
못하는 것으로 보인다 — 요구사항 문서 자체의 정확성 판단은 이 검증 작업
범위 밖이라 별도 보고로만 남긴다.

### 4. `order-saga.md` — 3단계 이상의 saga 보상 체인

재고 차감 → 결제 → 배송 예약, 3개 case가 `run` 캐스케이드로 연결된 하나의
최상위 Run 패스. 각 단계 `do`마다 `undo`가 걸려 있어 보상 스택 하나를
캐스케이드 전체가 공유한다(§undo 간선 구현 노트).

```
$ go run ./cmd/tangl check test/order-saga.md
ok

$ go run ./cmd/tangl effects test/order-saga.md "place order"
do   inventory.reserve   once (can reserve inventory / stock available)
undo inventory.release        (can reserve inventory / stock available)
do   payment.charge      once (can charge payment / payment method valid)
undo payment.refund           (can charge payment / payment method valid)
do   shipping.bookSlot   once (can book shipment / address deliverable)
undo shipping.cancelSlot      (can book shipment / address deliverable)
do   log.OrderComplete        (can book shipment / address deliverable)
```

**이 문서로 실제 버그를 하나 발견했다 — 아래 "발견된 버그" 절 참고.**
`can book shipment` case의 `address deliverable` 노드는 `do bookSlot once`
→ `undo cancelSlot` → `do log.OrderComplete` 순서로 선언되어 있다(스펙
문법상 완전히 유효). 코드 생성기는 같은 노드의 모든 `do`를 먼저 실행하고
그 다음에야 `undo`를 무장하므로, `bookSlot`이 성공한 뒤 `OrderComplete`가
실패하면 `cancelSlot` 보상이 전혀 무장되지 않는다.

### 5. `factory-line-cascade.md` — run 캐스케이드 4레벨 + checking 케이스 분리

전원 투입 → 가압 → 승온 → 생산 시작(캐스케이드 4레벨). 진동 이상 감지는
별도 case(`can flag anomaly`, 자기 `do log.AlertOps`를 가짐)로 분리해
생산 시작 case에서 `checking`으로만 참조한다.

```
$ go run ./cmd/tangl check test/factory-line-cascade.md
ok

$ go run ./cmd/tangl effects test/factory-line-cascade.md "start line"
do plc.energize          once (can power on / breaker closed)
do plc.openValve         once (can pressurize / pressure in range)
do plc.enableHeater      once (can heat / temperature in range)
do plc.startLine         once (can start production / ready to produce)
do log.ProductionStarted      (can start production / ready to produce)

$ go run ./cmd/tangl effects test/factory-line-cascade.md "diagnostics"
(공집합)
```

`start line`의 effect summary에 `log.AlertOps`(checking 전용 참조)가
**포함되지 않음**을 확인했고, `can flag anomaly`를 직접 `check`로 노출하는
`diagnostics` 엔드포인트의 effect summary가 (그 case 자신이 `do`를 가짐에도
불구하고) 공집합임을 확인했다 — 스펙의 두 MUST 요구사항을 모두 직접
검증한다.

### gen (모든 정상 문서 공통)

다섯 정상 문서 전부 `go run ./cmd/tangl gen test/<file>.md`가 exit 0으로
Go 소스를 stdout에 출력했고, `Generate` 내부에서 이미 `go/format.Source`를
거치므로 별도로 `gofmt -l`을 걸어도 변경 사항 없음(포맷 완료 상태)을
재확인했다. `order-saga.md`는 실제로 격리된 임시 모듈에 손으로 쓴 stub
패키지(`inventory`/`payment`/`shipping`/`log`)를 붙여 `go build`까지 통과함을
확인했다(아래 "발견된 버그" 재현 사례와 같은 방식, `pkg/tangl/gen`의
`TestCompileSyntheticDoc`이 쓰는 것과 동일한 replace-모듈 기법).

## 위반 문서 (check 실패, 정확한 에러로 거부되어야 함)

| 파일 | 위반 종류 | 검증 단계 | 실제 에러 메시지 |
|---|---|---|---|
| `bad-undo-without-do.md` | `undo`에 선행하는 `do`가 없음 | `validate.checkUndoRequiresDo` | `case "can refund": undo when "balance sufficient" has no preceding do on that node` |
| `bad-run-cycle.md` | `run` 캐스케이드 순환 (A→B→A) | `validate.checkRunCycle` | `run cycle detected at "can start a"` |
| `bad-section-order.md` | 섹션 순서 위반 (Provides가 Cases보다 앞) | `parser.ParseSource` (Validate 이전) | `tangl:Cases appears out of order` |

```
$ go run ./cmd/tangl check test/bad-undo-without-do.md
Error: test/bad-undo-without-do.md:22: case "can refund": undo when "balance sufficient" has no preceding do on that node
(exit 1)

$ go run ./cmd/tangl check test/bad-run-cycle.md
Error: test/bad-run-cycle.md:18: run cycle detected at "can start a"
(exit 1)

$ go run ./cmd/tangl check test/bad-section-order.md
Error: test/bad-section-order.md:20: tangl:Cases appears out of order
(exit 1)
```

세 파일 모두 정확히 기대한 에러로, 정확한 줄 번호와 함께 실패했다.

## 발견된 버그 — 같은 노드의 undo가 뒤따르는 do의 실패로 인해 무장되지 않음

**증상**: 한 노드에 `do A once` → `undo A'` → `do B` 순서로 선언되면(문법상
완전히 유효 — `order-saga.md`의 `can book shipment`가 정확히 이 모양이다),
`A`가 성공한 뒤 `B`가 실패할 경우 `A'`(A의 보상)가 **전혀 무장되지 않는다**.
스펙 §undo 간선은 "노드의 `do`가 성공 완료된 **직후** 보상이 무장된다"고
명시하고, 원칙 12는 "실패는 침묵하지 않는다 — 보상하거나 사람에게 올린다"고
못박는다. 그러나 실제로는 `A`의 실제 부수효과(예: 배송 슬롯 예약)가 이미
세상에 반영된 채로, 보상도 REVIEW 승급도 없이 조용히 원래 에러만 반환된다.

**원인 (코드 레벨)**: `pkg/tangl/gen/render_node_exec_block.go`가 한
노드의 실행 블록을 만들 때 `writeDoStatements`(그 노드의 모든 `do`를
문서 순서로, 각각 실패 시 즉시 `return err`)를 먼저 전부 쓰고, 그 다음에야
`writeUndoPushes`(그 노드의 모든 `undo`를 무조건 push)를 쓴다
(`write_undo_pushes.go`의 주석: "armed unconditionally once control
reaches this point in the handler — i.e. after every do on the node has
already succeeded"). 즉 실제 무장 시점은 "**그 do 직후**"가 아니라 "**그
노드의 모든 do가 끝난 후**"다. 같은 노드에 `do`가 하나뿐이면 두 표현이
같아 문제가 드러나지 않지만(스펙의 계좌이체 예시, `access`/`americano`
예시가 전부 이 패턴), 같은 노드에 `undo`로 보호받는 `do` 뒤에 **보호받지
않는 또 다른 `do`가 이어지면** 그 사이에서 조용한 데이터 유실이 생긴다.

**최소 재현 사례** (저장소 밖 스크래치 디렉터리에서 실제로 빌드·실행해
확인함 — `pkg/tangl`은 수정하지 않았다):

```
## tangl:Subject
- this document is `repro`

## tangl:See
- see `ext` from `repro/ext`

## tangl:Cases
- in case of `can run`
  - `x` is required
  - `trigger` is a general rule using `ext`.`Always`
  - do `ext`.`StepA` once when `trigger`
  - undo `ext`.`UndoA` when `trigger`
  - do `ext`.`StepB` when `trigger`

## tangl:Provides
- provides `go`
  - `x` is required
  - run `can run`
```

`ext.Always`는 항상 true, `ext.StepA`는 성공(실제 부수효과 있음),
`ext.StepB`는 항상 에러 반환, `ext.UndoA`는 보상(성공 시 로그 출력)이라고
하면, 생성된 `Go(ctx)` 실행 결과는:

```
StepA: ran (side effect applied)
Go() returned error: runOn "Always": StepB: boom
```

`UndoA: COMPENSATION FIRED`가 **한 번도 출력되지 않는다** — `StepA`의
효과는 세상에 남은 채, 보상도 REVIEW 승급도 없이 조용히 에러만 반환된다.

이 재현은 `test/order-saga.md`를 검증 문서로 그대로 두고(문법상 유효한
정상 문서이므로 수정하지 않음) 별도로 최소화한 것이다. `pkg/tangl` 프로덕션
코드는 이번 작업에서 수정하지 않았다(요청 범위 밖 — 버그 수정은 별도
작업으로 남긴다). 수정한다면 `write_undo_pushes`를 각 `do`문 **직후**로
옮기거나(도큐먼트 선언 순서와 무관하게, do–undo 페어 단위로 인터리브),
동일 노드에 `undo`로 보호되는 `do`가 있고 그 뒤에 무가드 `do`가 이어지면
파서/린트 경고를 추가하는 방향이 될 것이다.

## 파일 목록

```
test/
├── README.md                      (이 파일)
├── library-renewal.md             (a) checking + Rules either/and/or/not
├── newsletter-tick-warn.md        (b) once 없는 do → Lint 경고 재현
├── newsletter-tick-fixed.md       (b) once + 완료 가드로 고친 버전
├── order-saga.md                  (c) 3단계 saga 보상 체인 (+ 위 버그의 재현 근거)
├── factory-line-cascade.md        (d) run 캐스케이드 4레벨 + checking 분리 case
├── bad-undo-without-do.md         (e) 위반: undo without do
├── bad-run-cycle.md               (e) 위반: 순환 run 캐스케이드
└── bad-section-order.md           (e) 위반: 섹션 순서 위반
```
