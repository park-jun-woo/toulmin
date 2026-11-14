## tangl:Subject
- this document is `계좌이체`

## tangl:See
- see `은행` from `bank/core`
- see `기록` from `tangl/log`

## tangl:Cases
- in case of `출금 가능`
  - `보내는 계좌` is required
  - `금액` is required
  - `잔액 충분` is a general rule using `은행`.`hasBalance`
  - `계좌 동결` is a counter rule using `은행`.`isFrozen`
  - don't `잔액 충분` when `계좌 동결`
  - do `은행`.`withdraw` once when `잔액 충분`
  - undo `은행`.`refund` when `잔액 충분`
  - run `입금 가능` when `잔액 충분`

- in case of `입금 가능`
  - `받는 계좌` is required
  - `수취 계좌 유효` is a general rule using `은행`.`isValidAccount`
  - do `은행`.`deposit` once when `수취 계좌 유효`
  - do `기록`.`TransferComplete` when `수취 계좌 유효`

## tangl:Provides
- provides `이체`
  - `보내는 계좌` is required
  - `받는 계좌` is required
  - `금액` is required
  - run `출금 가능`
