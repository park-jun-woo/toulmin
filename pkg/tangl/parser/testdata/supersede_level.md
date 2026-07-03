## tangl:Subject
- this document is `reins`

## tangl:See
- see `policy` from `bank/policy`

## tangl:Cases
- in case of `can pass`
  - `has license` is a general rule using `policy`.`HasLicense`
  - `overdrawn` is a counter rule using `policy`.`IsOverdrawn`, level fail
  - `pending` is a counter rule using `policy`.`IsPending`, level review
  - `overdrawn` supersedes `pending`
