## tangl:Subject
- this document is `아메리카노`

## tangl:See
- see `팔` from `robot/arm`
- see `센서` from `robot/sensor`

## tangl:Definitions
- `서빙 온도` means 65°C as `센서`.`Temperature`

## tangl:Cases
- in case of `컵 놓기 가능`
  - `주문 접수` is a general rule using `센서`.`isOrdered`
  - `컵 없음` is a counter rule using `센서`.`noCupAvailable`
  - don't `주문 접수` when `컵 없음`
  - do `팔`.`placeCup` once when `주문 접수`
  - run `물 붓기 가능` when `주문 접수`

- in case of `물 붓기 가능`
  - `컵 놓임` is a general rule using `센서`.`cupIsPlaced`
  - do `팔`.`pourWater` once when `컵 놓임`
  - run `추출 가능` when `컵 놓임`

- in case of `추출 가능`
  - `물 부어짐` is a general rule using `센서`.`waterIsPoured`
  - do `팔`.`extractShot` once when `물 부어짐`
  - run `서빙 가능` when `물 부어짐`

- in case of `서빙 가능`
  - `추출 완료` is a general rule using `센서`.`isComplete`
  - `저온` is an except rule using `센서`.`temperatureLow` with `서빙 온도`
  - don't `추출 완료` when `저온`
  - do `팔`.`serve` once when `추출 완료`
  - do `기록`.`CoffeeStatus` when `추출 완료`

## tangl:Provides
- provides `커피 만들기`
  - run `컵 놓기 가능`

## tangl:Internal
- every 1s until `서빙 가능`
  - run `컵 놓기 가능`
