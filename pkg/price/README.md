# pkg/price

**Stop nesting if-else for discount logic. Declare rules, exceptions, and caps.**

Price judgment framework built on toulmin defeats graph. Discount conditions are warrants, stacking conflicts are defeat edges, audit trail is built-in. Qualifier stays at 1.0 — discount info lives in spec where it belongs.

User is `any` — the framework does not impose a concrete User type. Membership extraction is done via ctx fields.

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/price"
```

## Quick Start

```go
g := toulmin.NewGraph("product:discount")
coupon := g.Rule(price.HasCoupon).With(&price.DiscountSpec{
    Name: "SAVE30", Rate: 0.3, Max: 50000,
})
basic := g.Rule(price.IsMemberLevel).With(&price.MemberSpec{
    Level: "basic",
    Discount: &price.DiscountSpec{Name: "basic", Rate: 0.1},
})
noStack := g.Counter(price.IsAlreadyDiscounted)
noStack.Attacks(coupon)
noStack.Attacks(basic)

pricer := price.NewPricer(g, nil)
result, _ := pricer.Evaluate(req, ctx)
// result.FinalPrice, result.TotalDiscount, result.AppliedDiscounts
```

## DiscountSpec

```go
type DiscountSpec struct {
    Name  string  // "SAVE30", "basic", "blackfriday"
    Rate  float64 // percentage (0.3 = 30%). 0 = not applied
    Fixed float64 // fixed amount (5000). 0 = not applied
    Min   float64 // minimum discount guarantee. 0 = no minimum
    Max   float64 // maximum discount cap. 0 = no cap
}
```

Discount = BasePrice * Rate + Fixed, then clamped by Min/Max.

## MemberSpec

```go
type MemberSpec struct {
    Level    string        // "basic", "gold", "vip"
    Discount *DiscountSpec // discount to apply if matched
}
```

## BulkOrderSpec

```go
type BulkOrderSpec struct {
    MinQuantity int // minimum quantity for bulk order
}
```

## Rules

| Rule | Spec | Description |
|---|---|---|
| `HasCoupon` | *DiscountSpec | Coupon applies (meets min price) |
| `IsMemberLevel` | *MemberSpec | User membership matches Level (via ctx) |
| `HasActivePromotion` | *DiscountSpec | Named promotion is active |
| `IsAlreadyDiscounted` | nil | Purchase already discounted (Rebuttal) |
| `IsBulkOrder` | *BulkOrderSpec | Order quantity >= spec.MinQuantity |

## Pricer

```go
pricer := price.NewPricer(g, totalCap)  // totalCap optional
result, err := pricer.Evaluate(req, ctx)
```

1. Evaluates graph → collects activated discounts from evidence
2. Computes each: BasePrice * Rate + Fixed, clamped by Min/Max
3. Sums all discounts, applies totalCap if set
4. Returns FinalPrice, TotalDiscount, AppliedDiscounts, Trace
