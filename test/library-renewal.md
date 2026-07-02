# 도서관 대출 갱신 심사

회원의 연체 이력과 대기 예약 상태를 종합해 대출 갱신 가능 여부를 판정한다.
`tangl:Rules`의 `either`/`and`/`or`/`not`으로 판정 로직을 인라인 조합하고,
`checking`으로 서로 다른 사안(연체 판정, 연체료 면제 판정)의 verdict를 합성한다.

검증 포인트:
- `checking`이 참조하는 사안의 실행(`do`)은 **결코 발화하지 않는다** — `run
  `can process renewal request`` 가 실제로 실행돼도 `can renew loan`의
  `stampRenewal`/`RenewalResult`는 발화하지 않고, `can process renewal
  request` 자신의 `ProcessSummary`만 발화해야 한다.
- `check` 엔드포인트(`renewal eligibility`)는 부수효과가 0이어야 한다 —
  같은 `can renew loan` 사안을 `run`(엔드포인트 `renew`)으로도, `check`
  (엔드포인트 `renewal eligibility`)로도 노출해 Evaluate/Run 이중성을 직접
  대비시킨다.

## tangl:Subject
- this document is `library`

## tangl:See
- see `catalog` from `library/catalog`
- see `log` from `tangl/log`

## tangl:Definitions
- `long term membership years` means 3 as `catalog`.`Threshold`

## tangl:Rules
1. `no blocking overdue` when
    - either
      - `overdue count` is at most 0
      - and `membership tier` equals "gold"
    - or `overdue count` is less than 2
- `has pending hold from others` when
    - not `hold conflict` equals "none"

## tangl:Cases
- in case of `can renew loan`
  - `member id` is required
  - `book id` is required
  - `catalog record exists` is a general rule using `catalog`.`recordExists`
  - `no blocking overdue` is a general rule
  - `has pending hold from others` is a counter rule
  - don't `no blocking overdue` when `has pending hold from others`
  - don't `catalog record exists` when `has pending hold from others`
  - do `catalog`.`stampRenewal` once when `no blocking overdue` if at least 60% certain
  - do `log`.`RenewalResult` when `no blocking overdue`

- in case of `can waive late fee`
  - `member id` is required
  - `long term member` is a general rule using `catalog`.`isLongTermMember` with `long term membership years`
  - do `catalog`.`waiveLateFee` once when `long term member`

- in case of `can process renewal request`
  - `member id` is required
  - `book id` is required
  - `renewal allowed` is a general rule checking `can renew loan`
  - `fee waived` is a general rule checking `can waive late fee`
  - do `log`.`ProcessSummary` when `renewal allowed`

## tangl:Provides
- provides `renew`
  - `member id` is required
  - `book id` is required
  - run `can renew loan`
- provides `renewal eligibility`
  - `member id` is required
  - `book id` is required
  - check `can renew loan`
- provides `process renewal`
  - `member id` is required
  - `book id` is required
  - run `can process renewal request`
