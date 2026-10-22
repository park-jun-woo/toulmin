## tangl:Subject
- this document is `대출심사`

## tangl:Rules
- `신용 충분` when `신용점수` is at least 650
- `우대 대상` when
  - either
    - `등급` equals "VIP"
    - and `연체 이력` is empty
  - or `직원` equals true

## tangl:Cases
- in case of `승인 가능`
  - `신용 충분` is a general rule
  - `우대 대상` is a general rule
