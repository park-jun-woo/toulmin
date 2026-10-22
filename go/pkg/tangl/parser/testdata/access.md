## tangl:Subject
- this document is `api`

## tangl:See
- see `policy` from `bank/policy`
- see `log` from `tangl/log`

## tangl:Cases
- in case of `can access`
  - `user` is required
  - `authenticate` is a general rule
  - `block ip` is a counter rule using `policy`.`IsIPBlocked` with `blocklist`
  - `exempt internal ip` is an except rule using `policy`.`IsInternalIP`
  - don't `authenticate` when `block ip`
  - don't `block ip` when `exempt internal ip`
  - do `policy`.`Allow` when `authenticate`
  - do `log`.`AccessResult` when `authenticate`
  - do `policy`.`Deny` when `block ip`

## tangl:Provides
- provides `access`
  - `user` is required
  - run `can access`
