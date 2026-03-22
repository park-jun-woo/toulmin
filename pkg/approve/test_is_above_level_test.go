//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsAboveLevel — tests IsAboveLevel rule
package approve

import "testing"

func TestIsAboveLevel(t *testing.T) {
	tests := []struct {
		name  string
		level int
		min   int
		want  bool
	}{
		{"above", 5, 3, true},
		{"equal", 3, 3, true},
		{"below", 2, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := &ApproverBacking{Level: tt.min, LevelFunc: testAB.LevelFunc}
			ctx := &ApprovalContext{Approver: &testApprover{Level: tt.level}}
			got, _ := IsAboveLevel(nil, ctx, ab)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
