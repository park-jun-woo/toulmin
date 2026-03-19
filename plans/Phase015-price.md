# Phase 013: 가격 판정 프레임워크 — pkg/price

## 목표

toulmin 기반 가격 판정 프레임워크를 `pkg/price`에 구현한다.
할인, 쿠폰, 멤버십, 프로모션 등 가격 규칙의 충돌과 예외를 defeats graph로 판정한다.
qualifier를 할인 강도로 활용하여 verdict가 최종 할인율을 나타낸다.

## 배경

### 현재 문제

1. **할인 규칙이 if-else로 얽힌다**: `if hasCoupon && isMember && !isAlreadyDiscounted || (isVIP && isBlackFriday)` 같은 분기가 결제 로직에 박힌다
2. **중복 할인 방지가 하드코딩이다**: 쿠폰 + 멤버십 중복 불가, VIP는 중복 가능 같은 비즈니스 규칙이 조건문으로 관리된다. 규칙 추가/변경 시 사이드이펙트가 발생한다
3. **할인 우선순위가 암묵적이다**: 어떤 할인이 먼저 적용되는지, 어떤 할인이 다른 할인을 무효화하는지 코드를 전부 읽어야 파악된다

### toulmin이 해결하는 것

- 할인 규칙 하나 = rule 함수 (1-2 depth)
- 규칙 충돌 = defeats graph (명시적 선언)
- 중복 할인 방지 = rebuttal + defeat edge
- **qualifier = 할인 강도** (1.0이 아닌 실제 할인율)
- 판정 근거 = EvaluateTrace (어떤 할인이 왜 적용/무효화됐는지)

### claim/ground 분리 원칙

toulmin의 `(claim any, ground any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 가격 프레임워크에서 claim은 구매 요청 (상품, 수량, 원가)
- **ground = 판정 재료**: ground는 가격 판정에 필요한 컨텍스트 (사용자, 쿠폰, 프로모션)

프레임워크는 Pricer 구조와 판정 흐름을 제공하고, **도메인 데이터는 ground로 사용자가 주입한다.**

| 역할 | 가격 프레임워크에서 |
|---|---|
| claim | PurchaseRequest (ProductID, Quantity, BasePrice) |
| ground | PriceContext (User, Coupons, Promotions, Metadata) |
| rule 함수 | claim/ground에서 조건 하나만 판단 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 할인 예외) |
| qualifier | 할인 강도 (0.3 = 30% 할인) |
| verdict | 최종 할인 강도 [-1, +1] |

이전 프레임워크들과의 차이: **qualifier가 처음으로 1.0이 아닌 실제 비즈니스 값으로 사용된다.** 정책/상태 전이/승인에서 qualifier는 전부 1.0(pass/fail)이었지만, 가격에서는 할인율 자체가 qualifier가 되어 verdict가 최종 할인 강도를 나타낸다.

## 핵심 설계

### PurchaseRequest

```go
// pkg/price/purchase_request.go
type PurchaseRequest struct {
    ProductID string
    Quantity  int
    BasePrice float64
    Metadata  map[string]any
}
```

### PriceContext

```go
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
    Membership string  // "none", "basic", "vip"
}

// pkg/price/coupon.go
type Coupon struct {
    Code     string
    Discount float64  // 0.0 ~ 1.0
    MinPrice float64
    Category string
}

// pkg/price/promotion.go
type Promotion struct {
    Name     string
    Discount float64  // 0.0 ~ 1.0
    Active   bool
}
```

### 범용 rule 함수

```go
// pkg/price/rule_has_coupon.go
func HasCoupon(claim any, ground any) (bool, any) {
    req := claim.(*PurchaseRequest)
    ctx := ground.(*PriceContext)
    for _, c := range ctx.Coupons {
        if req.BasePrice >= c.MinPrice {
            return true, c
        }
    }
    return false, nil
}

// pkg/price/rule_is_member.go
func IsMember(membership string) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*PriceContext)
        return ctx.User.Membership == membership, nil
    }
}

// pkg/price/rule_is_vip.go
func IsVIP(claim any, ground any) (bool, any) {
    ctx := ground.(*PriceContext)
    return ctx.User.Membership == "vip", nil
}

// pkg/price/rule_has_active_promotion.go
func HasActivePromotion(name string) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*PriceContext)
        for _, p := range ctx.Promotions {
            if p.Name == name && p.Active {
                return true, p
            }
        }
        return false, nil
    }
}

// pkg/price/rule_is_already_discounted.go
func IsAlreadyDiscounted(claim any, ground any) (bool, any) {
    req := claim.(*PurchaseRequest)
    discounted, _ := req.Metadata["discounted"].(bool)
    return discounted, nil
}

// pkg/price/rule_is_bulk_order.go
func IsBulkOrder(minQty int) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        req := claim.(*PurchaseRequest)
        return req.Quantity >= minQty, nil
    }
}
```

### Pricer — 가격 판정 실행

```go
// pkg/price/pricer.go
type Pricer struct {
    graph *toulmin.Graph
}

func NewPricer(g *toulmin.Graph) *Pricer

// Evaluate — 최종 할인 판정
func (p *Pricer) Evaluate(req *PurchaseRequest, ctx *PriceContext) (*PriceResult, error)

// FinalPrice — 최종 가격 계산 (BasePrice * (1 - discount))
func (p *Pricer) FinalPrice(req *PurchaseRequest, ctx *PriceContext) (float64, error)
```

### PriceResult

```go
// pkg/price/price_result.go
type PriceResult struct {
    BasePrice  float64
    Discount   float64            // verdict로부터 계산된 할인율
    FinalPrice float64
    Trace      []toulmin.TraceEntry
}
```

### qualifier를 할인 강도로 사용

**주의**: 클로저 rule은 변수에 저장 후 재사용해야 한다. Rebuttal만으로는 공격이 일어나지 않으며 반드시 Defeat edge를 선언해야 한다. 예외를 처리하는 rule은 Defeater로 등록해야 한다.

```go
// 클로저는 변수에 저장 후 재사용
isMemberBasic := price.IsMember("basic")
blackfriday := price.HasActivePromotion("blackfriday")
bulkOrder := price.IsBulkOrder(100)

// 쿠폰 30%, 멤버십 10%, 중복 할인 방지, VIP는 중복 허용
g := toulmin.NewGraph("product:discount").
    Warrant(price.HasCoupon, 0.3).              // 30% 할인
    Warrant(isMemberBasic, 0.1).                // 10% 할인
    Warrant(blackfriday, 0.2).                  // 20% 할인
    Rebuttal(price.IsAlreadyDiscounted, 1.0).   // 중복 할인 방지
    Defeater(price.IsVIP, 1.0).                 // 예외 rule은 Defeater로 등록
    Defeater(bulkOrder, 1.0).
    Defeat(price.IsAlreadyDiscounted, price.HasCoupon).     // Rebuttal → Warrant 공격 edge 필수
    Defeat(price.IsVIP, price.IsAlreadyDiscounted).         // VIP는 중복 할인 허용
    Defeat(bulkOrder, price.IsAlreadyDiscounted)            // 대량 주문도 중복 허용

pricer := price.NewPricer(g)

result, err := pricer.Evaluate(req, ctx)
// result.Discount: h-Categoriser가 계산한 최종 할인 강도
// result.FinalPrice: BasePrice * (1 - Discount)
// result.Trace: 어떤 할인이 적용/무효화됐는지
```

### 사용 예시

```go
req := &price.PurchaseRequest{
    ProductID: "prod-001",
    Quantity:  1,
    BasePrice: 100000,
}

ctx := &price.PriceContext{
    User:    &price.User{ID: "user-1", Membership: "vip"},
    Coupons: []price.Coupon{{Code: "SAVE30", Discount: 0.3, MinPrice: 50000}},
    Promotions: []price.Promotion{{Name: "blackfriday", Discount: 0.2, Active: true}},
}

finalPrice, err := pricer.FinalPrice(req, ctx)
// VIP이므로 중복 할인 허용 → 쿠폰 + 멤버십 + 프로모션 적용
// h-Categoriser가 qualifier 기반으로 최종 할인 강도 계산
```

## 범위

### 포함

1. **PurchaseRequest, PriceContext 구조체**: 가격 판정에 필요한 요청/컨텍스트
2. **User, Coupon, Promotion 구조체**: 도메인 모델
3. **범용 rule 함수**: HasCoupon, IsMember, IsVIP, HasActivePromotion, IsAlreadyDiscounted, IsBulkOrder
4. **Pricer**: 가격 판정 실행 (Evaluate, FinalPrice)
5. **PriceResult**: 판정 결과 (할인율, 최종 가격, trace)
6. **테스트**: rule 함수 단위 테스트, Pricer 통합 테스트

### 제외

- 통화 변환
- 세금 계산
- 결제 연동
- 할인 이력 퍼시스턴스 (DB 저장은 사용자 책임)

## 산출물

```
pkg/
  price/
    purchase_request.go            — PurchaseRequest 구조체
    price_context.go               — PriceContext 구조체
    user.go                        — User 구조체
    coupon.go                      — Coupon 구조체
    promotion.go                   — Promotion 구조체
    rule_has_coupon.go             — HasCoupon
    rule_is_member.go              — IsMember (클로저)
    rule_is_vip.go                 — IsVIP
    rule_has_active_promotion.go   — HasActivePromotion (클로저)
    rule_is_already_discounted.go  — IsAlreadyDiscounted
    rule_is_bulk_order.go          — IsBulkOrder (클로저)
    pricer.go                      — Pricer (NewPricer, Evaluate, FinalPrice)
    price_result.go                — PriceResult 구조체
    rule_test.go                   — rule 함수 단위 테스트
    pricer_test.go                 — Pricer 통합 테스트
```

## 단계

### Step 1: 구조체 정의

- PurchaseRequest, PriceContext, User, Coupon, Promotion
- PriceResult

### Step 2: rule 함수 구현

- 각 rule 함수를 파일 하나에 하나씩 구현 (filefunc 규칙 준수)
- 각 함수는 1-2 depth 유지
- 클로저 rule: IsMember, HasActivePromotion, IsBulkOrder

### Step 3: Pricer 구현

- NewPricer: graph를 받아 Pricer 생성
- Evaluate: graph.EvaluateTrace → verdict를 할인율로 변환 → PriceResult
- FinalPrice: Evaluate → BasePrice * (1 - Discount)

### Step 4: 테스트

- rule 함수 단위 테스트: 각 조건별 true/false
- Pricer 통합 테스트:
  - 쿠폰만 적용 → 30% 할인
  - 쿠폰 + 멤버십 + 중복 방지 rebuttal → 하나만 적용
  - VIP defeat → 중복 할인 허용
  - 프로모션 비활성 → 할인 미적용
  - 대량 주문 defeat → 중복 할인 허용
  - trace에 각 rule의 activated/defeated 상태 포함

### Step 5: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. HasCoupon, IsMember 등 rule 함수가 claim/ground에서 올바르게 판정한다
2. qualifier가 할인율로 작동한다 (0.3 = 30%)
3. Pricer.Evaluate가 h-Categoriser verdict를 할인율로 변환한다
4. Pricer.FinalPrice가 BasePrice * (1 - Discount)를 정확히 계산한다
5. IsAlreadyDiscounted rebuttal이 중복 할인을 방지한다
6. IsVIP defeat가 중복 할인 방지를 무효화한다
7. PriceResult.Trace에 각 rule의 판정 근거가 포함된다
8. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
