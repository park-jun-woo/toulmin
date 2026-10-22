package parser

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseListAndFileDefinitions(t *testing.T) {
	list, err := parseDefinitionItem(item{Text: "`domains` means [\"a,b\", \"c\"]", Line: 2}, "d.md")
	if err != nil || list.Kind != ast.StringListDef || !reflect.DeepEqual(list.Strings, []string{"a,b", "c"}) {
		t.Fatalf("list = %#v, err=%v", list, err)
	}
	file, err := parseDefinitionItem(item{Text: "`domains` reads `JSON list` from \"./domains.json\"", Line: 3}, "d.md")
	if err != nil || file.Kind != ast.FileSpecDef || file.LoadType != "JSON list" || file.FilePath != "./domains.json" {
		t.Fatalf("file = %#v, err=%v", file, err)
	}
	for _, source := range []string{"[", "[1]", "[\"a\"] junk"} {
		if _, err := parseDefinitionItem(item{Text: "`x` means " + source, Line: 4}, "d.md"); err == nil {
			t.Errorf("accepted %q", source)
		}
	}
}
