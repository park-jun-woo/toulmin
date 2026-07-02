# 뉴스레터 다이제스트 발송 — 경고 유발판 (검증용 반례)

**주의**: 이 문서는 의도적으로 결함이 있는 반례다. 실제 배포용이 아니라
`tangl check`의 틱 멱등성 Lint 경고를 재현하기 위한 것이다. 정상적으로
고친 버전은 `newsletter-tick-fixed.md`를 보라.

`tangl:Internal`이 매 15분마다 `can send digest`를 통째로 재-Run한다.
`subscription active`가 활성인 동안(=구독이 살아있는 동안) 그 상태는 매 틱
계속 참이므로, `once`도 완료 가드도 없는 `do`는 매 틱 다시 발화한다 —
뉴스레터가 15분마다 중복 발송된다. `tangl check`는 이 무가드 `do`를
경고해야 한다(파싱은 성공, `Validate`도 통과 — Lint는 에러가 아니라
경고이므로 `check`의 종료 코드는 0이어야 한다).

## tangl:Subject
- this document is `newsletterwarn`

## tangl:See
- see `mailer` from `newsletter/mailer`

## tangl:Cases
- in case of `can send digest`
  - `subscriber id` is required
  - `subscription active` is a general rule using `mailer`.`isSubscribed`
  - do `mailer`.`sendDigest` when `subscription active`

## tangl:Provides
- provides `send now`
  - `subscriber id` is required
  - run `can send digest`

## tangl:Internal
- every 15m
  - run `can send digest`
