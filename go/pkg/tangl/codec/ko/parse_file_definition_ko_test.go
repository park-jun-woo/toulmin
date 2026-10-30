//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseFileDefinitionKo — '`타입` 을 "경로" 파일에서 읽는다' parses; other bodies are unmatched; bad type/path fail
package ko

import (
	"reflect"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseFileDefinitionKo(t *testing.T) {
	def, matched, err := parseFileDefinitionKo("도메인", "`JSON 목록` 을 \"./d.json\" 파일에서 읽는다", "t.md", 5)
	if err != nil || !matched || !reflect.DeepEqual(def, ast.Definition{Name: "도메인", Kind: ast.FileSpecDef, LoadType: "JSON 목록", FilePath: "./d.json", Line: 5}) {
		t.Fatalf("valid: %#v matched=%v err=%v", def, matched, err)
	}
	if _, matched, err := parseFileDefinitionKo("x", "3 을 뜻한다", "t.md", 5); matched || err != nil {
		t.Fatalf("non-file body: matched=%v err=%v", matched, err)
	}
	for _, c := range []struct{ body, want string }{
		{"JSON 을 \"d.json\" 파일에서 읽는다", "expected backtick-quoted loader type"},
		{"`` 을 \"d.json\" 파일에서 읽는다", "expected backtick-quoted loader type"},
		{"`JSON 목록` 을 d.json 파일에서 읽는다", "expected one non-empty quoted file path"},
		{"`JSON 목록` 을 \"\" 파일에서 읽는다", "expected one non-empty quoted file path"},
	} {
		_, matched, err := parseFileDefinitionKo("x", c.body, "t.md", 5)
		if !matched || err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: matched=%v err=%v, want %q", c.body, matched, err, c.want)
		}
	}
}
