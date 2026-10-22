# 아메리카노 제조

고객이 주문하면 컵을 놓고, 물을 붓고, 추출하여 서빙한다. 센서 갱신이 늦으면
다음 틱이 이어받는다 — 팔 동작은 once로 틱 재실행에서 보호된다.

## tangl:Subject
- this document is `americano`

## tangl:See
- see `arm` from `robot/arm`
- see `sensor` from `robot/sensor`

## tangl:Definitions
- `serve temperature` means 65°C as `sensor`.`Temperature`

## tangl:Cases
- in case of `can place cup`
  - `order received` is a general rule using `sensor`.`isOrdered`
  - `no cup` is a counter rule using `sensor`.`noCupAvailable`
  - don't `order received` when `no cup`
  - do `arm`.`placeCup` once when `order received`
  - run `can pour water` when `order received`

- in case of `can pour water`
  - `cup placed` is a general rule using `sensor`.`cupIsPlaced`
  - do `arm`.`pourWater` once when `cup placed`
  - run `can brew espresso` when `cup placed`

- in case of `can brew espresso`
  - `water poured` is a general rule using `sensor`.`waterIsPoured`
  - do `arm`.`extractShot` once when `water poured`
  - run `can serve coffee` when `water poured`

- in case of `can serve coffee`
  - `brewing complete` is a general rule using `sensor`.`isComplete`
  - `low temperature` is an except rule using `sensor`.`temperatureLow` with `serve temperature`
  - don't `brewing complete` when `low temperature`
  - do `arm`.`serve` once when `brewing complete`
  - do `log`.`CoffeeStatus` when `brewing complete`

## tangl:Provides
- provides `make coffee`
  - run `can place cup`

## tangl:Internal
- every 1s until `can serve coffee`
  - run `can place cup`
