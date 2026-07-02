//ff:func feature=policy type=rule control=sequence
//ff:what TestIsRateLimited_Branches — covers rateLimiter and clientIP type-assertion failure branches of IsRateLimited
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsRateLimited_Branches(t *testing.T) {
	t.Run("limiter wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("clientIP", "1.2.3.4")
		ctx.Set("rateLimiter", "not-a-limiter")

		got, evidence := IsRateLimited(ctx, nil)
		if got {
			t.Errorf("expected false when rateLimiter is not a RateLimiter, got %v", got)
		}
		if evidence != nil {
			t.Errorf("expected nil evidence, got %v", evidence)
		}
	})

	t.Run("limiter unset", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("clientIP", "1.2.3.4")

		got, _ := IsRateLimited(ctx, nil)
		if got {
			t.Errorf("expected false when rateLimiter is unset, got %v", got)
		}
	})

	t.Run("clientIP wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		limiter := &mockLimiter{limited: map[string]bool{"1.2.3.4": true}}
		ctx.Set("rateLimiter", limiter)
		ctx.Set("clientIP", 12345)

		got, _ := IsRateLimited(ctx, nil)
		if got {
			t.Errorf("expected false when clientIP is not a string, got %v", got)
		}
	})
}
