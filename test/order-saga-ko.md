# 주문 처리 사가 (재고 → 결제 → 배송) — 한국어판

`test/order-saga.md`의 한국어 감사 표면판이다. 도메인 식별자(`inventory`,
`payment`, `shipping`, `log`와 case 이름)는 영어 코드베이스를 그대로 가리키고,
문장만 서술어형 정규형으로 옮겼다 — 같은 AST로 흡수되어 effects/gen 산출이
영어판과 동일하다. 재고를 차감하고, 결제를 받고, 배송을 예약하는 세 단계가
`진행한다` 캐스케이드로 연결된 하나의 최상위 Run 패스이므로 보상 스택도 하나를
공유한다: 배송 예약이 실패하면 이미 발화한 결제·재고 차감이 발화 역순(LIFO)으로
보상된다.

## tangl:주제
- 이 문서는 `order` 이다

## tangl:참조
- `inventory` 는 `shop/inventory` 에서 참조한다
- `payment` 는 `shop/payment` 에서 참조한다
- `shipping` 는 `shop/shipping` 에서 참조한다
- `log` 는 `tangl/log` 에서 참조한다

## tangl:사안
- `can reserve inventory` 의 경우에
  - `sku` 는 필수다
  - `quantity` 는 필수다
  - `stock available` 은 `inventory`.`hasStock` 쓰는 일반 규칙이다
  - `sku discontinued` 는 `inventory`.`isDiscontinued` 쓰는 예외 규칙이다
  - `sku discontinued` 이면 `stock available` 무효다
  - `stock available` 이면 `inventory`.`reserve` 한 번만 실행한다
  - `stock available` 이면 `inventory`.`release` 로 되돌린다
  - `stock available` 이면 `can charge payment` 진행한다

- `can charge payment` 의 경우에
  - `customer id` 는 필수다
  - `amount` 는 필수다
  - `payment method valid` 는 `payment`.`isValidMethod` 쓰는 일반 규칙이다
  - `card blocked` 는 `payment`.`isCardBlocked` 쓰는 반박 규칙이다
  - `card blocked` 이면 `payment method valid` 무효다
  - `payment method valid` 이면 `payment`.`charge` 한 번만 실행한다
  - `payment method valid` 이면 `payment`.`refund` 로 되돌린다
  - `payment method valid` 이면 `can book shipment` 진행한다

- `can book shipment` 의 경우에
  - `shipping address` 는 필수다
  - `address deliverable` 는 `shipping`.`isDeliverable` 쓰는 일반 규칙이다
  - `address deliverable` 이면 `shipping`.`bookSlot` 한 번만 실행한다
  - `address deliverable` 이면 `shipping`.`cancelSlot` 로 되돌린다
  - `address deliverable` 이면 `log`.`OrderComplete` 실행한다

## tangl:업무
- `place order` 를 한다면
  - `sku` 는 필수다
  - `quantity` 는 필수다
  - `customer id` 는 필수다
  - `amount` 는 필수다
  - `shipping address` 는 필수다
  - `can reserve inventory` 진행한다
