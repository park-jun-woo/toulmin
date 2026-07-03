//ff:func feature=tangl type=parser control=sequence
//ff:what TestVerifyCover — tests VerifyCover accepts a full tiling and rejects gaps and short cover
package scan

import "testing"

func TestVerifyCover(t *testing.T) {
	good := []Line{
		{Span: Span{Off: 0, Len: 3}, No: 1},
		{Span: Span{Off: 3, Len: 2}, No: 2},
	}
	if err := VerifyCover(good, 5); err != nil {
		t.Errorf("VerifyCover(good): unexpected error %v", err)
	}

	// A gap between spans is rejected.
	gap := []Line{
		{Span: Span{Off: 0, Len: 3}, No: 1},
		{Span: Span{Off: 4, Len: 2}, No: 2},
	}
	if err := VerifyCover(gap, 6); err == nil {
		t.Errorf("VerifyCover(gap): expected error, got nil")
	}

	// Spans that do not reach total are rejected.
	if err := VerifyCover(good, 6); err == nil {
		t.Errorf("VerifyCover(short): expected error, got nil")
	}
}
