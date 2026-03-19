# Phase 015: 가격 판정 프레임워크 — pkg/price (구현 완료)

## 목표

toulmin 기반 가격 판정 프레임워크를 `pkg/price`에 구현한다.
할인 규칙의 적용 여부를 defeats graph로 판정한다.
할인 정보(비율/정액/최소/최대)는 backing에 담는다.
qualifier는 Toulmin 원래 목적(확신도)으로만 사용한다.

## 배경

### 설계 원칙: qualifier ≠ 할인율

qualifier는 Toulmin의 확신도(confidence weight)이다. 할인율을 qualifier에 넣으면 Toulmin 모델을 오용하는 것이다.

- **qualifier** = 이 규칙이 얼마나 확실한가 (0.0~1.0)
- **backing** = 이 규칙의 판정 기준 + 할인 정보 (Rate, Fixed, Min, Max)

```go
// 잘못됨 — qualifier를 할인율로 오용
g.Warrant(HasCoupon, nil, 0.3)

// 올바름 — qualifier는 확신도, 할인 정보는 backing
g.Warrant(HasCoupon, &DiscountBacking{Name: "SAVE30", Rate: 0.3, Max: 50000}, 1.0)
```

### 현재 문제

1. **할인 규칙이 if-else로 얽힌다**: 쿠폰 + 멤버십 + 프로모션 + 중복 방지 + VIP 예외가 분기로 관리된다
2. **중복 할인 방지가 하드코딩이다**: 어떤 할인끼리 겹치면 안 되는지가 if문 안에 있다
3. **총 할인 상한이 별도 로직이다**: 최대 30% 제한 같은 규칙이 할인 계산 후 별도로 적용된다

### toulmin이 해결하는 것

- 할인 조건 = warrant (이 할인이 적용되는가?)
- 할인 정보 = backing (Rate, Fixed, Min, Max)
- 중복 방지 = rebuttal + defeat edge
- 판정 근거 = EvaluateTrace (어떤 할인이 왜 적용/차단됐는지)

### claim/ground/backing 분리

| 역할 | 가격 프레임워크에서 |
|---|---|
| claim | PurchaseRequest (ProductID, Quantity, BasePrice) |
| ground | PriceContext (User, Coupons, Promotions) |
| backing | *DiscountBacking (Name, Rate, Fixed, Min, Max) |
| rule 함수 | 할인 조건 판정 (1-2 depth), evidence로 DiscountBacking 반환 |
| graph | 할인 간 충돌/예외 관계 |
| verdict | 할인 적용 여부. Pricer가 activated rule의 evidence에서 할인 계산 |

## 핵심 설계

### DiscountBacking

```go
// pkg/price/discount_backing.go
type DiscountBacking struct {
    Name  string  // 할인 이름 ("SAVE30", "basic", "blackfriday")
    Rate  float64 // 비율 할인 (0.3 = 30%). 0이면 미적용
    Fixed float64 // 정액 할인 (5000원). 0이면 미적용
    Min   float64 // 최소 할인 보장. 0이면 제한 없음
    Max   float64 // 최대 할인 상한. 0이면 제한 없음
}
```

Rate + Fixed 조합으로 모든 할인 유형을 표현한다:

```go
// 30% 할인, 최대 50000원
&DiscountBacking{Name: "SAVE30", Rate: 0.3, Max: 50000}

// 정액 5000원 할인
&DiscountBacking{Name: "welcome", Fixed: 5000}

// 20% + 정액 2000원, 최소 3000원 보장
&DiscountBacking{Name: "combo", Rate: 0.2, Fixed: 2000, Min: 3000}

// VIP — 할인 자체는 없고, 중복 허용 defeater용
&DiscountBacking{Name: "vip"}
```

### PurchaseRequest / PriceContext

```go
// pkg/price/purchase_request.go
type PurchaseRequest struct {
    ProductID string
    Quantity  int
    BasePrice float64
    Metadata  map[string]any
}

// pkg/price/price_context.go
type PriceContext struct {
    User       *User
    Coupons    []Coupon
    Promotions []Promotion
    Metadata   map[string]any
}

// pkg/price/user.go
type User struct {
    ID         string
    Membership string // "none", "basic", "gold", "vip"
}

// pkg/price/coupon.go
type Coupon struct {
    Code     string
    MinPrice float64
}

// pkg/price/promotion.go
type Promotion struct {
    Name   string
    Active bool
}
```

### 범용 rule 함수

```go
// pkg/price/rule_has_coupon.go
// backing: *DiscountBacking
func HasCoupon(claim any, ground any, backing any) (bool, any) {
    req := claim.(*PurchaseRequest)
    ctx := ground.(*PriceContext)
    db := backing.(*DiscountBacking)
    for _, c := range ctx.Coupons {
        if req.BasePrice >= c.MinPrice {
            return true, db
        }
    }
    return false, nil
}

// pkg/price/rule_is_member_level.go
// backing: *DiscountBacking (Name = 멤버십 등급)
func IsMemberLevel(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*PriceContext)
    db := backing.(*DiscountBacking)
    return ctx.User.Membership == db.Name, db
}

// pkg/price/rule_has_active_promotion.go
// backing: *DiscountBacking (Name = 프로모션 이름)
func HasActivePromotion(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*PriceContext)
    db := backing.(*DiscountBacking)
    for _, p := range ctx.Promotions {
        if p.Name == db.Name && p.Active {
            return true, db
        }
    }
    return false, nil
}

// pkg/price/rule_is_already_discounted.go
// backing: nil
func IsAlreadyDiscounted(claim any, ground any, backing any) (bool, any) {
    req := claim.(*PurchaseRequest)
    discounted, _ := req.Metadata["discounted"].(bool)
    return discounted, nil
}

// pkg/price/rule_is_bulk_order.go
// backing: int (최소 수량)
func IsBulkOrder(claim any, ground any, backing any) (bool, any) {
    req := claim.(*PurchaseRequest)
    minQty := backing.(int)
    return req.Quantity >= minQty, nil
}
```

### Pricer — 가격 판정 + 할인 계산

```go
// pkg/price/pricer.go
type Pricer struct {
    graph    *toulmin.Graph
    totalCap *DiscountBacking // optional: 총 할인 상한 (Max 필드만 사용)
}

func NewPricer(g *toulmin.Graph, totalCap *DiscountBacking) *Pricer

// Evaluate — 할인 판정 + 최종 가격 계산
func (p *Pricer) Evaluate(req *PurchaseRequest, ctx *PriceContext) (*PriceResult, error)
```

Evaluate 흐름:
1. `graph.EvaluateTrace(req, ctx)` → activated warrant 목록
2. activated warrant의 evidence에서 `*DiscountBacking` 수집
3. 각 DiscountBacking의 할인액 계산:
   - `discount = BasePrice * Rate + Fixed`
   - Min 적용: `discount = max(discount, Min)` (Min > 0일 때)
   - Max 적용: `discount = min(discount, Max)` (Max > 0일 때)
4. 전체 할인액 합산
5. totalCap 적용 (있으면): `totalDiscount = min(totalDiscount, totalCap.Max)`
6. `FinalPrice = BasePrice - totalDiscount`

### PriceResult

```go
// pkg/price/price_result.go
type PriceResult struct {
    BasePrice        float64
    TotalDiscount    float64
    FinalPrice       float64
    AppliedDiscounts []*DiscountBacking
    Trace            []toulmin.TraceEntry
}
```

### 사용 예시

```go
g := toulmin.NewGraph("product:discount")

// 할인 조건 — qualifier는 전부 1.0
coupon := g.Warrant(price.HasCoupon, &price.DiscountBacking{
    Name: "SAVE30", Rate: 0.3, Max: 50000,
}, 1.0)
basic := g.Warrant(price.IsMemberLevel, &price.DiscountBacking{
    Name: "basic", Rate: 0.1,
}, 1.0)
promo := g.Warrant(price.HasActivePromotion, &price.DiscountBacking{
    Name: "blackfriday", Fixed: 5000,
}, 1.0)

// 중복 방지
noStack := g.Rebuttal(price.IsAlreadyDiscounted, nil, 1.0)
g.Defeat(noStack, coupon)
g.Defeat(noStack, basic)

// VIP는 중복 허용
vip := g.Defeater(price.IsMemberLevel, &price.DiscountBacking{
    Name: "vip",
}, 1.0)
g.Defeat(vip, noStack)

// 대량 주문도 중복 허용
bulk := g.Defeater(price.IsBulkOrder, 100, 1.0)
g.Defeat(bulk, noStack)

// 총 할인 상한: 최대 50000원
totalCap := &price.DiscountBacking{Max: 50000}
pricer := price.NewPricer(g, totalCap)

result, _ := pricer.Evaluate(req, ctx)
// result.FinalPrice: 할인 적용 후 최종 가격
// result.AppliedDiscounts: 적용된 할인 목록
// result.TotalDiscount: 총 할인액 (cap 적용 후)
// result.Trace: 판정 근거
```

## 범위

### 포함

1. **DiscountBacking 구조체**: Rate/Fixed/Min/Max 할인 정보
2. **PurchaseRequest, PriceContext 구조체**: 요청/컨텍스트
3. **User, Coupon, Promotion 구조체**: 도메인 모델
4. **범용 rule 함수**: HasCoupon, IsMemberLevel, HasActivePromotion, IsAlreadyDiscounted, IsBulkOrder
5. **Pricer**: 할인 판정 + 개별 할인 계산(Rate+Fixed+Min+Max) + 총합 cap + 최종 가격
6. **PriceResult**: 결과 (할인액, 최종가, 적용 목록, trace)
7. **테스트**: rule 단위 + Pricer 통합

### 제외

- 통화 변환, 세금 계산, 결제 연동
- 할인 이력 퍼시스턴스

## 산출물

```
pkg/
  price/
    discount_backing.go            — DiscountBacking (Name, Rate, Fixed, Min, Max)
    purchase_request.go            — PurchaseRequest 구조체
    price_context.go               — PriceContext 구조체
    user.go                        — User 구조체
    coupon.go                      — Coupon 구조체
    promotion.go                   — Promotion 구조체
    price_result.go                — PriceResult 구조체
    rule_has_coupon.go             — HasCoupon (backing: *DiscountBacking)
    rule_is_member_level.go        — IsMemberLevel (backing: *DiscountBacking)
    rule_has_active_promotion.go   — HasActivePromotion (backing: *DiscountBacking)
    rule_is_already_discounted.go  — IsAlreadyDiscounted (backing: nil)
    rule_is_bulk_order.go          — IsBulkOrder (backing: int)
    pricer.go                      — Pricer (NewPricer, Evaluate)
    rule_test.go                   — rule 함수 단위 테스트
    pricer_test.go                 — Pricer 통합 테스트
```

## 검증 기준

1. qualifier는 전부 1.0 — 할인율로 사용하지 않는다
2. DiscountBacking의 Rate/Fixed/Min/Max 조합이 정확히 계산된다
3. activated warrant의 evidence에서 DiscountBacking을 수집하여 할인 합산한다
4. totalCap으로 총 할인 상한을 적용한다
5. 중복 방지 rebuttal + VIP/대량 주문 defeat가 정확히 동작한다
6. IsMemberLevel이 backing.Name으로 모든 멤버십 등급을 구분한다
7. PriceResult에 적용된 할인 목록과 trace가 포함된다
8. 전체 테스트 PASS

## 의존성

- Phase 010: backing 일급 시민
- Phase 011: Rule 참조 반환, Defeat(*Rule, *Rule)
