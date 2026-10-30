//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseFileDefinition — "reads `Type` from \"path\"" parses; missing type, from, or a quoted path is an error
package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseFileDefinition(t *testing.T) {
	def, matched, err := parseFileDefinition("d", ` reads `+"`JSON list`"+` from "./d.json"`, "t.md", 4)
	if err != nil || !matched {
		t.Fatalf("valid: matched=%v err=%v", matched, err)
	}
	if !reflect.DeepEqual(def, ast.Definition{Name: "d", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "./d.json", Line: 4}) {
		t.Fatalf("def = %#v", def)
	}
	if _, matched, err := parseFileDefinition("d", " means 3", "t.md", 4); matched || err != nil {
		t.Fatalf("non-file body: matched=%v err=%v", matched, err)
	}
	for _, c := range []struct{ body, want string }{
		{"reads JSON from \"d.json\"", "expected backtick-quoted loader type"},
		{"reads `` from \"d.json\"", "expected backtick-quoted loader type"},
		{"reads `JSON list` at \"d.json\"", "expected 'from' after loader type"},
		{"reads `JSON list` from d.json", "expected one non-empty quoted file path"},
		{"reads `JSON list` from \"\"", "expected one non-empty quoted file path"},
	} {
		_, matched, err := parseFileDefinition("d", c.body, "t.md", 4)
		if !matched || err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: matched=%v err=%v, want %q", c.body, matched, err, c.want)
		}
	}
}
