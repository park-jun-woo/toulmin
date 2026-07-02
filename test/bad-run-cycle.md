# [의도적 위반] 순환 run 캐스케이드

Cases는 반드시 비순환 DAG여야 한다(§비순환 제약, `checkRunCycle`).
`can start a`가 `can start b`를 `run`하고, `can start b`가 다시 `can
start a`를 `run`하는 순환을 만들어 `tangl check`가 정확히 거부하는지
검증한다.

기대: `go run ./cmd/tangl check test/bad-run-cycle.md` 는 0이 아닌 종료
코드로 실패하고, 메시지에 `run cycle detected` 가 포함되어야 한다.

## tangl:Subject
- this document is `badcycle`

## tangl:See
- see `x` from `cycle/x`

## tangl:Cases
- in case of `can start a`
  - `flag` is a general rule using `x`.`isFlagged`
  - run `can start b` when `flag`

- in case of `can start b`
  - `flag` is a general rule using `x`.`isFlagged`
  - run `can start a` when `flag`

## tangl:Provides
- provides `start`
  - `flag` is required
  - run `can start a`
