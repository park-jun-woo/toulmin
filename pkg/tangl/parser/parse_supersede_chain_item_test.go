//ff:what TestParseSupersedeChainItem — chain expansion, validation, and duplicate edge normalization
package parser

import "testing"

func TestParseSupersedeChainItem(t *testing.T) {
	edges, ok, err := parseSupersedeChainItem(item{Text: "supersede in order: `A` > `B` > `C`", Line: 7}, "t.md")
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	want := [][2]string{{"A", "B"}, {"A", "C"}, {"B", "C"}}
	if len(edges) != len(want) {
		t.Fatalf("len=%d edges=%+v", len(edges), edges)
	}
	for i, pair := range want {
		if edges[i].Superseder != pair[0] || edges[i].Superseded != pair[1] || edges[i].Line != 7 {
			t.Fatalf("edge[%d]=%+v want=%v", i, edges[i], pair)
		}
	}
	for _, text := range []string{
		"supersede in order: `A` > `B`",
		"supersede in order: `A` > `B` > `A`",
		"supersede in order: `A` x `B` > `C`",
	} {
		if _, ok, err := parseSupersedeChainItem(item{Text: text, Line: 1}, "t.md"); !ok || err == nil {
			t.Errorf("%q: ok=%v err=%v", text, ok, err)
		}
	}
	six, _, err := parseSupersedeChainItem(item{Text: "supersede in order: `A` > `B` > `C` > `D` > `E` > `F`", Line: 8}, "t.md")
	if err != nil || len(six) != 15 {
		t.Fatalf("six-name closure len=%d err=%v", len(six), err)
	}
}

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
