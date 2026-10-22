//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCertainty — "if at least N% certain" gate, not used by the spec's three examples
package parser

import "testing"

// TestParseCertainty exercises a do edge's certainty clause.
func TestParseCertainty(t *testing.T) {
	src := "" +
		"## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `x`\n" +
		"  - `a` is a general rule\n" +
		"  - do `arm`.`placeCup` once when `a` if at least 75% certain\n"
	doc, err := ParseSource(src, "certainty.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	execs := doc.Cases[0].Execs
	if len(execs) != 1 {
		t.Fatalf("Execs = %+v, want 1", execs)
	}
	cert := execs[0].Certainty
	if cert == nil || cert.Op != "at least" || cert.Percent != 75 {
		t.Errorf("Certainty = %+v, want {at least 75}", cert)
	}
}
