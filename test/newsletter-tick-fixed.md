# 뉴스레터 다이제스트 발송 — 정상판

`newsletter-tick-warn.md`의 결함을 두 가지 방식으로 고친 버전이다:

1. **`once` 가드** — `sendDigest`/`logDigestSent` 모두 `once`를 붙여 ctx
   생애주기(이 틱 러너 프로세스가 살아있는 동안, `every`가 공유하는 ctx)당
   최대 1회만 발화하게 한다.
2. **완료 상태를 반영하는 counter 가드** — `already sent today`를
   `counter rule`로 등록하고 `don't `subscription active` when `already
   sent today`` 로 무력화 간선을 걸어, 오늘 이미 보냈으면 트리거 노드 자체가
   더 이상 Active가 아니게 만든다(§틱 멱등성 요구의 (c)). `once`와 이 가드는
   서로 다른 방어선이며 함께 두어도 안전하다(둘 다 있어도 중복 실행은
   일어나지 않는다).

`tangl check`는 경고 없이 통과해야 한다.

## tangl:Subject
- this document is `newsletterfixed`

## tangl:See
- see `mailer` from `newsletter/mailer`

## tangl:Cases
- in case of `can send digest`
  - `subscriber id` is required
  - `subscription active` is a general rule using `mailer`.`isSubscribed`
  - `already sent today` is a counter rule using `mailer`.`hasSentToday`
  - don't `subscription active` when `already sent today`
  - do `mailer`.`sendDigest` once when `subscription active`
  - do `mailer`.`logDigestSent` once when `subscription active`

## tangl:Provides
- provides `send now`
  - `subscriber id` is required
  - run `can send digest`

## tangl:Internal
- every 15m
  - run `can send digest`
