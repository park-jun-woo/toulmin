//ff:func feature=tangl type=parser control=sequence
//ff:what TestSupersedeChainDeduplicatesIndividualEdge — a chain plus an explicit duplicate edge keeps all source edges
package parser

import "testing"

func TestSupersedeChainDeduplicatesIndividualEdge(t *testing.T) {
	src := "## tangl:Subject\n- this document is `x`\n\n## tangl:Cases\n- in case of `c`\n" +
		"  - `A` is a general rule\n  - `B` is a counter rule\n  - `C` is an except rule\n" +
		"  - supersede in order: `A` > `B` > `C`\n  - `A` supersedes `B`\n"
	doc, err := ParseSource(src, "t.md")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(doc.Cases[0].Supersedes); got != 4 {
		t.Fatalf("source supersedes=%d want=4 (3 expansion + 1 explicit): %+v", got, doc.Cases[0].Supersedes)
	}
}
