## tangl:Subject
- this document is `pricing`

## tangl:See
- see `분류` from `mail/분류`

## tangl:Rules
- `신뢰 사용자` when `작성자`.`신뢰점수` is at least spec `신뢰 기준`.`최소 점수`
- `쿠폰 있음` when any of `쿠폰들` as `c` satisfies `c`.`최소가격` is at most `기본가`
- `IP 목록` when any of `차단목록` as `ip` satisfies `ip` equals `클라이언트IP`
- `활성 프로모션` when any of `프로모션들` as `p` satisfies
  - `p`.`활성` equals true
  - and `p`.`할인` is at least 10
- `게시 가능` when
  - `분류`.`isSpam` holds
  - and result of `분류`.`score` is greater than 0.7
  - or any of `쿠폰들` as `c` satisfies
    - `c`.`상태` equals "활성"
    - and `c`.`최소가격` is less than `기본가`

## tangl:Cases
- in case of `승인`
  - `신뢰 사용자` is a general rule
  - `쿠폰 있음` is a general rule
  - `IP 목록` is a general rule
  - `활성 프로모션` is a general rule
  - `게시 가능` is a general rule
