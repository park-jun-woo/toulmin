//ff:func feature=tangl type=model control=sequence
//ff:what TestVersion — version stamp constants hold GEUL v0.3.0, profile 0x0001, codebook v0.2
package codebook

import "testing"

func TestVersion(t *testing.T) {
	if GEULVersionMajor != 0 || GEULVersionMinor != 3 || GEULVersionPatch != 0 {
		t.Errorf("GEUL version: got %d.%d.%d, want 0.3.0",
			GEULVersionMajor, GEULVersionMinor, GEULVersionPatch)
	}
	if TANGLProfileID != 0x0001 {
		t.Errorf("TANGLProfileID: got 0x%04X, want 0x0001", TANGLProfileID)
	}
	if CodebookVersionMajor != 0 || CodebookVersionMinor != 2 {
		t.Errorf("codebook version: got %d.%d, want 0.2",
			CodebookVersionMajor, CodebookVersionMinor)
	}
}
