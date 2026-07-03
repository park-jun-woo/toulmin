//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseDefinitionFieldsKo — parse a StructDef's child field items
package ko

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseDefinitionFieldsKo parses the nested "`f` 를 가진다, 타입 T" child items of a
// struct 뜻한다 definition into ast.Field values.
func parseDefinitionFieldsKo(items []scan.Item, path string) ([]ast.Field, error) {
	fields := make([]ast.Field, 0, len(items))
	for _, it := range items {
		f, err := parseFieldKo(it, path)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, nil
}
