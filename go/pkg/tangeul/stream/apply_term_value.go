//ff:func feature=tangl type=util control=selection
//ff:what applyTermValue — put a means-edge value into a Definition by kind (JSON list, JSON file spec, or literal)
package stream

import (
	"encoding/json"
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// applyTermValue stores value, the TEXT a term's means edge points at, into def
// according to def.Kind: a string-list definition decodes it as a JSON string
// array, a file definition as {loadType, filePath}, and every other kind keeps
// it as the literal value. Malformed JSON is an error.
func applyTermValue(def *ast.Definition, value string) error {
	switch def.Kind {
	case ast.StringListDef:
		if err := json.Unmarshal([]byte(value), &def.Strings); err != nil {
			return fmt.Errorf("decode string-list definition %q: %w", def.Name, err)
		}
	case ast.FileSpecDef:
		var payload struct {
			LoadType string `json:"loadType"`
			FilePath string `json:"filePath"`
		}
		if err := json.Unmarshal([]byte(value), &payload); err != nil {
			return fmt.Errorf("decode file definition %q: %w", def.Name, err)
		}
		def.LoadType, def.FilePath = payload.LoadType, payload.FilePath
	default:
		def.Value = value
	}
	return nil
}
