## tangl:Subject
- this document is `퀘스트 판정`

## tangl:See
- see `게이트` from `reins/gate`

## tangl:Cases
- in case of `주문 접수`
  - `유효 주문` is a general rule using `게이트`.`isValidOrder`
  - `컵 없음` is a counter rule using `게이트`.`noCup`, level fail
  - `재고 부족` is an except rule using `게이트`.`lowStock`, level review
  - don't `유효 주문` when `컵 없음`
  - `컵 없음` supersedes `주문 접수`
