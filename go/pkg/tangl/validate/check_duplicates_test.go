//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckDuplicates — tests checkDuplicates for first-occurrence, second-occurrence-report, and already-reported-skip branches via subtests
package validate

import (
	"strings"
	"testing"
)

func TestCheckDuplicates(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		errs := checkDuplicates("test.md", "thing", nil)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("AllUnique", func(t *testing.T) {
		locs := []nameLoc{
			{Name: "a", Line: 1},
			{Name: "b", Line: 2},
		}
		errs := checkDuplicates("test.md", "thing", locs)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("SecondOccurrenceReports", func(t *testing.T) {
		locs := []nameLoc{
			{Name: "a", Line: 1},
			{Name: "a", Line: 3},
		}
		errs := checkDuplicates("test.md", "thing", locs)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "first declared at line 1") {
			t.Errorf("expected mention of first declaration line, got %v", errs[0])
		}
	})

	t.Run("ThirdOccurrenceSkipsAlreadyReported", func(t *testing.T) {
		locs := []nameLoc{
			{Name: "a", Line: 1},
			{Name: "a", Line: 3},
			{Name: "a", Line: 5},
		}
		errs := checkDuplicates("test.md", "thing", locs)
		if len(errs) != 1 {
			t.Fatalf("expected exactly 1 error (not one per repeat), got %v", errs)
		}
	})
}
