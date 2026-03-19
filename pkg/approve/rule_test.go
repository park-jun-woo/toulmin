package approve

import "testing"

type mockOrgTree struct {
	managers map[string]string // requesterID → managerID
}

func (m *mockOrgTree) IsDirectManager(approverID, requesterID string) bool {
	return m.managers[requesterID] == approverID
}
func (m *mockOrgTree) Level(userID string) int { return 0 }

func TestIsDirectManager(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}
	tests := []struct {
		name     string
		approver string
		want     bool
	}{
		{"is manager", "mgr-1", true},
		{"not manager", "mgr-2", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ApprovalRequest{RequesterID: "emp-1"}
			ctx := &ApprovalContext{Approver: &Approver{ID: tt.approver}, OrgTree: org}
			got, _ := IsDirectManager(req, ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestIsBudgetFrozen(t *testing.T) {
	tests := []struct {
		name   string
		frozen bool
		want   bool
	}{
		{"frozen", true, true},
		{"not frozen", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &ApprovalContext{Budget: &Budget{Frozen: tt.frozen}}
			got, _ := IsBudgetFrozen(nil, ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasApprovalRole(t *testing.T) {
	tests := []struct {
		name string
		role string
		have string
		want bool
	}{
		{"match", "finance", "finance", true},
		{"mismatch", "finance", "manager", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &ApprovalContext{Approver: &Approver{Role: tt.have}}
			got, _ := HasApprovalRole(nil, ctx, tt.role)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

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
			ctx := &ApprovalContext{Approver: &Approver{Level: tt.level}}
			got, _ := IsAboveLevel(nil, ctx, tt.min)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSmallAmount(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		threshold float64
		want      bool
	}{
		{"small", 5000, 10000, true},
		{"equal", 10000, 10000, true},
		{"large", 15000, 10000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ApprovalRequest{Amount: tt.amount}
			got, _ := IsSmallAmount(req, nil, tt.threshold)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUrgent(t *testing.T) {
	tests := []struct {
		name string
		meta map[string]any
		want bool
	}{
		{"urgent", map[string]any{"urgent": true}, true},
		{"not urgent", map[string]any{}, false},
		{"nil meta", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ApprovalRequest{Metadata: tt.meta}
			got, _ := IsUrgent(req, nil, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCEOOverride(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{"ceo", "ceo", true},
		{"not ceo", "manager", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &ApprovalContext{Approver: &Approver{Role: tt.role}}
			got, _ := IsCEOOverride(nil, ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
