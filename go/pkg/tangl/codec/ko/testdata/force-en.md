## tangl:Subject
- this document is `공격세기`

## tangl:Rules
- `신용 충분` when `신용점수` is at least 650
- `신용 부족` when `신용점수` is less than 500
- `연체 있음` when `연체 이력` is not empty
- `내부 직원` when `직원` equals true

## tangl:Cases
- in case of `승인 가능`
  - `신용 충분` is a general rule
  - `신용 부족` is a counter rule qualified 0.8
  - `연체 있음` is a counter rule, level fail
  - `내부 직원` is an except rule
  - don't `신용 충분` when `신용 부족`
  - `연체 있음` attacks `신용 충분` with force 0.5
  - `내부 직원` attacks `연체 있음` with force 1
  - `내부 직원` attacks `신용 부족` with force 0.25
