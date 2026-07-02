# 공장 라인 기동 캐스케이드

전원 투입 → 가압 → 승온 → 생산 시작까지 4단계(`run` 캐스케이드 3단 이상)로
표현한다. 진동 이상 감지는 별도 사안(`can flag anomaly`)으로 분리해 생산
시작 사안에서 `checking`으로만 참조한다.

검증 포인트: `tangl effects factoryline.`start line`` 이 캐스케이드
전체(전원 → 가압 → 승온 → 생산)의 `do`만 문서 순서로 열거하고, `checking`
으로만 참조되는 `can flag anomaly`의 `AlertOps`는 **포함하지 않아야** 한다
— `checking`은 Evaluate로 컴파일되어 순수하기 때문이다. 대비되는 두 번째
엔드포인트 `diagnostics`는 `can flag anomaly`를 `check`로 직접 노출하는데,
`check` 엔드포인트의 effect summary는 스펙상 항상 공집합이어야 한다(설령
그 case 자신이 `do`를 갖고 있어도).

## tangl:Subject
- this document is `factoryline`

## tangl:See
- see `plc` from `factory/plc`
- see `sensor` from `factory/sensor`
- see `log` from `tangl/log`

## tangl:Cases
- in case of `can flag anomaly`
  - `line id` is required
  - `vibration abnormal` is a general rule using `sensor`.`isVibrationAbnormal`
  - do `log`.`AlertOps` when `vibration abnormal`

- in case of `can power on`
  - `line id` is required
  - `breaker closed` is a general rule using `plc`.`isBreakerClosed`
  - do `plc`.`energize` once when `breaker closed`
  - run `can pressurize` when `breaker closed`

- in case of `can pressurize`
  - `pressure in range` is a general rule using `sensor`.`isPressureInRange`
  - do `plc`.`openValve` once when `pressure in range`
  - run `can heat` when `pressure in range`

- in case of `can heat`
  - `temperature in range` is a general rule using `sensor`.`isTemperatureInRange`
  - do `plc`.`enableHeater` once when `temperature in range`
  - run `can start production` when `temperature in range`

- in case of `can start production`
  - `anomaly detected` is a general rule checking `can flag anomaly`
  - `ready to produce` is a general rule using `plc`.`isReadyToProduce`
  - don't `ready to produce` when `anomaly detected`
  - do `plc`.`startLine` once when `ready to produce`
  - do `log`.`ProductionStarted` when `ready to produce`

## tangl:Provides
- provides `start line`
  - `line id` is required
  - run `can power on`
- provides `diagnostics`
  - `line id` is required
  - check `can flag anomaly`
