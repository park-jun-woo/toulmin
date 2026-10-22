//ff:func feature=cli type=command control=sequence
//ff:what TestSplitSubjectEndpoint — splitSubjectEndpoint splits "<subject>.<endpoint>" on its first dot
package tanglcli

import "testing"

// TestSplitSubjectEndpoint covers both branches of splitSubjectEndpoint: a
// missing '.' is an error, and a present '.' splits into subject and
// (backtick-trimmed) endpoint.
func TestSplitSubjectEndpoint(t *testing.T) {
	t.Run("no dot is an error", func(t *testing.T) {
		subject, endpoint, err := splitSubjectEndpoint("noDotHere")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if subject != "" || endpoint != "" {
			t.Errorf("subject=%q endpoint=%q, want both empty", subject, endpoint)
		}
	})

	t.Run("dot splits subject and bare endpoint", func(t *testing.T) {
		subject, endpoint, err := splitSubjectEndpoint("transfer.transfer")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if subject != "transfer" {
			t.Errorf("subject = %q, want %q", subject, "transfer")
		}
		if endpoint != "transfer" {
			t.Errorf("endpoint = %q, want %q", endpoint, "transfer")
		}
	})

	t.Run("dot splits subject and backticked endpoint", func(t *testing.T) {
		subject, endpoint, err := splitSubjectEndpoint("transfer.`transfer`")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if subject != "transfer" {
			t.Errorf("subject = %q, want %q", subject, "transfer")
		}
		if endpoint != "transfer" {
			t.Errorf("endpoint = %q, want %q (backticks trimmed)", endpoint, "transfer")
		}
	})
}
