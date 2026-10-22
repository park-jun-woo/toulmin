## tangl:Subject
- this document is `pricing`

## tangl:See
- see `nlp` from `mail/nlp`
- see `check` from `mail/check`

## tangl:Rules
- `credit ok` when `credit score` is at least 650
- `trusted author` when `author`.`trustScore` is at least 0.7
- `is internal` when `attributes`.`internal` equals true
- `within budget` when spec `discount`.`min price` is at most `basePrice`
- `high nlp` when result of `nlp`.`score` is at least 0.7
- `spam` when `check`.`isSpam` holds
- `hasCoupon` when any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`
- `activePromotion` when any of `promotions` as `p` satisfies
  - `p`.`active` equals true
  - and `p`.`discount` is at least 10
- `ipInList` when any of `blocklist` as `ip` satisfies `ip` equals `clientIp`
- `combo` when
  - any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`
  - or `credit score` is at least 700
- `flagged` when
  - `check`.`isSpam` holds
  - or `manualFlag` equals true
- `no manual flag` when
  - not `manualFlag` equals true
- `vip member` when
  - either
    - `grade` equals VIP
    - and `overdue` is empty
  - or `staff` equals true
- `grade listed` when `grade` is in VIP, GOLD
- `meta ok` when `credit ok` holds

## tangl:Cases
- in case of `approve`
  - `credit ok` is a general rule
  - `trusted author` is a general rule
  - `hasCoupon` is a general rule
  - `activePromotion` is a general rule
  - `ipInList` is a general rule
