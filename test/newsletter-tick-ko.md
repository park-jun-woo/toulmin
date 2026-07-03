# 뉴스레터 다이제스트 발송 — 경고 유발판 (한국어판, 검증용 반례)

`test/newsletter-tick-warn.md`의 한국어 감사 표면판이다. 의도적으로 결함이 있는
반례로, `tangl check`의 틱 멱등성 Lint 경고를 재현하기 위한 것이다. `tangl:내부`가
매 15분마다 `can send digest`를 통째로 재-Run하는데, 그 안의 `do`
(`mailer`.`sendDigest`)에 `한 번만`도 완료 가드도 없다 → 매 틱 중복 발송된다.
`check`는 이 무가드 `do`를 경고해야 한다(파싱·Validate는 통과, 경고는 에러가
아니므로 종료 코드는 0).

## tangl:주제
- 이 문서는 `newsletterwarn` 이다

## tangl:참조
- `mailer` 는 `newsletter/mailer` 에서 참조한다

## tangl:사안
- `can send digest` 의 경우에
  - `subscriber id` 는 필수다
  - `subscription active` 은 `mailer`.`isSubscribed` 쓰는 일반 규칙이다
  - `subscription active` 이면 `mailer`.`sendDigest` 실행한다

## tangl:업무
- `send now` 를 한다면
  - `subscriber id` 는 필수다
  - `can send digest` 진행한다

## tangl:내부
- 매 15분마다
  - `can send digest` 진행한다
