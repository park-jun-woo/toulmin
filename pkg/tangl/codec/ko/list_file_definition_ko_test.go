package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestListAndFileDefinitionKo(t *testing.T) {
	list, err := parseDefinitionKo(scan.Item{Text: "`도메인` 은 [\"a,b\", \"c\"] 을 뜻한다", No: 2}, "d.md")
	if err != nil || list.Kind != ast.StringListDef || !reflect.DeepEqual(list.Strings, []string{"a,b", "c"}) {
		t.Fatalf("list = %#v, err=%v", list, err)
	}
	file, err := parseDefinitionKo(scan.Item{Text: "`도메인` 은 `JSON 목록` 을 \"./domains.json\" 파일에서 읽는다", No: 3}, "d.md")
	if err != nil || file.Kind != ast.FileSpecDef || file.LoadType != "JSON 목록" || file.FilePath != "./domains.json" {
		t.Fatalf("file = %#v, err=%v", file, err)
	}
	if got := renderDefinitionKo(list); got != "`도메인` 은 [\"a,b\",\"c\"] 을 뜻한다" {
		t.Fatalf("render list = %q", got)
	}
	if got := renderDefinitionKo(file); got != "`도메인` 은 `JSON 목록` 을 \"./domains.json\" 파일에서 읽는다" {
		t.Fatalf("render file = %q", got)
	}
}
