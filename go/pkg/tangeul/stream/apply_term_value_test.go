//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestApplyTermValue — list, file and literal values land in the Definition; malformed JSON fails
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestApplyTermValue(t *testing.T) {
	list := ast.Definition{Kind: ast.StringListDef}
	if err := applyTermValue(&list, `["a","b"]`); err != nil || !reflect.DeepEqual(list.Strings, []string{"a", "b"}) {
		t.Fatalf("list: %+v %v", list, err)
	}
	file := ast.Definition{Kind: ast.FileSpecDef}
	if err := applyTermValue(&file, `{"loadType":"lines","filePath":"x.txt"}`); err != nil || file.LoadType != "lines" || file.FilePath != "x.txt" {
		t.Fatalf("file: %+v %v", file, err)
	}
	lit := ast.Definition{Kind: ast.ConstDef}
	if err := applyTermValue(&lit, "5"); err != nil || lit.Value != "5" {
		t.Fatalf("literal: %+v %v", lit, err)
	}
	for _, k := range []ast.DefKind{ast.StringListDef, ast.FileSpecDef} {
		if err := applyTermValue(&ast.Definition{Kind: k}, "{bad"); err == nil {
			t.Fatalf("kind %d bad JSON: want error", k)
		}
	}
}
