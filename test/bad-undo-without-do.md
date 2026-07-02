# [의도적 위반] undo without do

`undo`는 같은 case 안에서 같은 노드에 선행하는 `do`가 있어야 한다
(§undo 간선, 파서 검증 `checkUndoRequiresDo`). 여기서는 `balance
sufficient` 노드에 `do` 없이 `undo`만 걸어, `tangl check`가 정확히
거부하는지 검증한다.

기대: `go run ./cmd/tangl check test/bad-undo-without-do.md` 는 0이 아닌
종료 코드로 실패하고, 메시지는
`case "can refund" undo when "balance sufficient" has no preceding do on that node` 형태여야 한다.

## tangl:Subject
- this document is `badundo`

## tangl:See
- see `bank` from `bank/core`

## tangl:Cases
- in case of `can refund`
  - `account id` is required
  - `balance sufficient` is a general rule using `bank`.`hasBalance`
  - undo `bank`.`refund` when `balance sufficient`

## tangl:Provides
- provides `refund`
  - `account id` is required
  - run `can refund`
