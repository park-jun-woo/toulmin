# [의도적 위반] 섹션 순서 위반

일곱 섹션의 순서는 Subject → See → Definitions → Rules → Cases → Provides
→ Internal 로 고정이다(§일곱 섹션, `ParseSource`의 `sectionOrderIndex`
검사). 여기서는 `tangl:Provides`를 `tangl:Cases`보다 앞에 두어 파서가
정확히 거부하는지 검증한다.

기대: `go run ./cmd/tangl check test/bad-section-order.md` 는 0이 아닌
종료 코드로 실패하고, 메시지에 `tangl:Cases appears out of order` 가
포함되어야 한다(파서 단계 에러 — `Validate` 이전에 `parser.Parse`에서
멈춘다).

## tangl:Subject
- this document is `badorder`

## tangl:Provides
- provides `dummy`
  - run `can do thing`

## tangl:Cases
- in case of `can do thing`
  - `x` is required
