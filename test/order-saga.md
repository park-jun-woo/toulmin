# 주문 처리 사가 (재고 → 결제 → 배송)

재고를 차감하고, 결제를 받고, 배송을 예약한다. 세 단계가 `run` 캐스케이드로
연결된 하나의 최상위 Run 패스이므로 보상 스택도 하나를 공유한다 — 배송
예약 단계에서 실패하면 이미 발화한 결제·재고 차감이 **발화 역순(LIFO)**
으로 보상된다: 배송 실패 → (결제 없음, 배송 단계는 undo 없음이므로 스킵)
→ 결제 환불 → 재고 복원. 결제 단계에서 실패하면 재고만 복원된다.

검증 포인트: 3단계 이상의 사가 체인, `undo`가 같은 case 안에서 선행하는
`do`를 요구하는 파서 검증(§undo), `tangl effects`가 캐스케이드 전체의
`do`/`undo`를 문서 선언 순서로 열거하는지.

## tangl:Subject
- this document is `order`

## tangl:See
- see `inventory` from `shop/inventory`
- see `payment` from `shop/payment`
- see `shipping` from `shop/shipping`
- see `log` from `tangl/log`

## tangl:Cases
- in case of `can reserve inventory`
  - `sku` is required
  - `quantity` is required
  - `stock available` is a general rule using `inventory`.`hasStock`
  - `sku discontinued` is an except rule using `inventory`.`isDiscontinued`
  - don't `stock available` when `sku discontinued`
  - do `inventory`.`reserve` once when `stock available`
  - undo `inventory`.`release` when `stock available`
  - run `can charge payment` when `stock available`

- in case of `can charge payment`
  - `customer id` is required
  - `amount` is required
  - `payment method valid` is a general rule using `payment`.`isValidMethod`
  - `card blocked` is a counter rule using `payment`.`isCardBlocked`
  - don't `payment method valid` when `card blocked`
  - do `payment`.`charge` once when `payment method valid`
  - undo `payment`.`refund` when `payment method valid`
  - run `can book shipment` when `payment method valid`

- in case of `can book shipment`
  - `shipping address` is required
  - `address deliverable` is a general rule using `shipping`.`isDeliverable`
  - do `shipping`.`bookSlot` once when `address deliverable`
  - undo `shipping`.`cancelSlot` when `address deliverable`
  - do `log`.`OrderComplete` when `address deliverable`

## tangl:Provides
- provides `place order`
  - `sku` is required
  - `quantity` is required
  - `customer id` is required
  - `amount` is required
  - `shipping address` is required
  - run `can reserve inventory`
