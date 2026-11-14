# 계좌 이체

출금하고 입금한다. 입금이 실패하면 출금을 환불한다.
어느 시점에 실패해도 계좌 총합은 보존되거나, 사람의 검토로 올라간다.

## tangl:Subject
- this document is `transfer`

## tangl:See
- see `bank` from `bank/core`
- see `log` from `tangl/log`

## tangl:Cases
- in case of `can withdraw`
  - `from account` is required
  - `amount` is required
  - `balance sufficient` is a general rule using `bank`.`hasBalance`
  - `account frozen` is a counter rule using `bank`.`isFrozen`
  - don't `balance sufficient` when `account frozen`
  - do `bank`.`withdraw` once when `balance sufficient`
  - undo `bank`.`refund` when `balance sufficient`
  - run `can deposit` when `balance sufficient`

- in case of `can deposit`
  - `to account` is required
  - `recipient valid` is a general rule using `bank`.`isValidAccount`
  - do `bank`.`deposit` once when `recipient valid`
  - do `log`.`TransferComplete` when `recipient valid`

## tangl:Provides
- provides `transfer`
  - `from account` is required
  - `to account` is required
  - `amount` is required
  - run `can withdraw`
