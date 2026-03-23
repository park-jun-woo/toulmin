//ff:func feature=feature type=rule control=sequence
//ff:what TestIsUserInPercentage — tests IsUserInPercentage rule
package feature

import "testing"

func TestIsUserInPercentage(t *testing.T) {
	ctx := &UserContext{ID: "user-123"}
	h := hashPercentage("user-123")

	got1, _ := IsUserInPercentage(nil, ctx, &PercentageBacking{Percentage: 1.0})
	if !got1 {
		t.Error("expected true for 100% rollout")
	}

	got2, _ := IsUserInPercentage(nil, ctx, &PercentageBacking{Percentage: 0.0})
	if got2 {
		t.Error("expected false for 0% rollout")
	}

	// deterministic
	got3, _ := IsUserInPercentage(nil, ctx, &PercentageBacking{Percentage: h + 0.01})
	got4, _ := IsUserInPercentage(nil, ctx, &PercentageBacking{Percentage: h + 0.01})
	if got3 != got4 {
		t.Error("expected deterministic result")
	}
}
