//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseStringListDefinition — a JSON string array parses; non-arrays are not matched; bad JSON or trailing values fail
package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseStringListDefinition(t *testing.T) {
	def, matched, err := parseStringListDefinition("d", ` ["a,b", "c"]`, "t.md", 2)
	if err != nil || !matched || def.Kind != ast.StringListDef || def.Name != "d" || def.Line != 2 || !reflect.DeepEqual(def.Strings, []string{"a,b", "c"}) {
		t.Fatalf("valid: %#v matched=%v err=%v", def, matched, err)
	}
	if _, matched, err := parseStringListDefinition("d", "3", "t.md", 2); matched || err != nil {
		t.Fatalf("scalar: matched=%v err=%v", matched, err)
	}
	for _, c := range []struct{ src, want string }{
		{"[1, 2]", "invalid JSON string array"},
		{"[\"a\"", "invalid JSON string array"},
		{"[\"a\"] [\"b\"]", "unexpected value after JSON string array"},
	} {
		_, matched, err := parseStringListDefinition("d", c.src, "t.md", 2)
		if !matched || err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: matched=%v err=%v, want %q", c.src, matched, err, c.want)
		}
	}
}
