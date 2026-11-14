//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseStringListDefinitionKo — a JSON string array parses; non-arrays are unmatched; bad JSON or trailing values fail
package ko

import (
	"reflect"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseStringListDefinitionKo(t *testing.T) {
	def, matched, err := parseStringListDefinitionKo("도메인", ` ["a,b", "c"]`, "t.md", 2)
	if err != nil || !matched || def.Kind != ast.StringListDef || def.Name != "도메인" || def.Line != 2 || !reflect.DeepEqual(def.Strings, []string{"a,b", "c"}) {
		t.Fatalf("valid: %#v matched=%v err=%v", def, matched, err)
	}
	if _, matched, err := parseStringListDefinitionKo("x", "3", "t.md", 2); matched || err != nil {
		t.Fatalf("scalar: matched=%v err=%v", matched, err)
	}
	for _, c := range []struct{ src, want string }{
		{"[1]", "invalid JSON string array"},
		{"[\"a\"", "invalid JSON string array"},
		{"[\"a\"] 1", "unexpected value after JSON string array"},
	} {
		_, matched, err := parseStringListDefinitionKo("x", c.src, "t.md", 2)
		if !matched || err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: matched=%v err=%v, want %q", c.src, matched, err, c.want)
		}
	}
}
