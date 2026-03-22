//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsUnderBudget — tests IsUnderBudget rule
package approve

import "testing"

func TestIsUnderBudget(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		remaining float64
		want      bool
	}{
		{"under", 5000, 10000, true},
		{"equal", 10000, 10000, true},
		{"over", 15000, 10000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ApprovalRequest{Amount: tt.amount}
			ctx := &ApprovalContext{Budget: &Budget{Remaining: tt.remaining}}
			got, _ := IsUnderBudget(req, ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
