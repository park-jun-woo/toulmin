## tangl:Subject
- this document is `접근제어`

## tangl:See
- see `정책` from `bank/policy`
- see `기록` from `tangl/log`

## tangl:Cases
- in case of `접근 가능`
  - `사용자` is required
  - `인증` is a general rule
  - `IP 차단` is a counter rule using `정책`.`IsIPBlocked` with `차단 목록`
  - `내부 IP 예외` is an except rule using `정책`.`IsInternalIP`
  - don't `인증` when `IP 차단`
  - don't `IP 차단` when `내부 IP 예외`
  - do `정책`.`Allow` when `인증`
  - do `기록`.`AccessResult` when `인증`
  - do `정책`.`Deny` when `IP 차단`

## tangl:Provides
- provides `접근`
  - `사용자` is required
  - run `접근 가능`
