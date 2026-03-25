//ff:func feature=feature type=rule control=sequence
//ff:what TestIsUserInPercentage — tests IsUserInPercentage rule
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsUserInPercentage(t *testing.T) {
	h := hashPercentage("user-123")

	ctx1 := toulmin.NewContext()
	ctx1.Set("id", "user-123")
	got1, _ := IsUserInPercentage(ctx1, &PercentageBacking{Percentage: 1.0})
	if !got1 {
		t.Error("expected true for 100% rollout")
	}

	ctx2 := toulmin.NewContext()
	ctx2.Set("id", "user-123")
	got2, _ := IsUserInPercentage(ctx2, &PercentageBacking{Percentage: 0.0})
	if got2 {
		t.Error("expected false for 0% rollout")
	}

	// deterministic
	ctx3 := toulmin.NewContext()
	ctx3.Set("id", "user-123")
	got3, _ := IsUserInPercentage(ctx3, &PercentageBacking{Percentage: h + 0.01})
	ctx4 := toulmin.NewContext()
	ctx4.Set("id", "user-123")
	got4, _ := IsUserInPercentage(ctx4, &PercentageBacking{Percentage: h + 0.01})
	if got3 != got4 {
		t.Error("expected deterministic result")
	}
}
